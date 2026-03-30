// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package partitioner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundRobinConfig_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, (&RoundRobinConfig{EventsPerPartition: 1}).Validate())
		assert.NoError(t, (&RoundRobinConfig{EventsPerPartition: 100}).Validate())
	})
	t.Run("zero_is_invalid", func(t *testing.T) {
		assert.Error(t, (&RoundRobinConfig{EventsPerPartition: 0}).Validate())
	})
	t.Run("negative_is_invalid", func(t *testing.T) {
		assert.Error(t, (&RoundRobinConfig{EventsPerPartition: -1}).Validate())
	})
}

func TestRandomConfig_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, (&RandomConfig{EventsPerPartition: 1}).Validate())
		assert.NoError(t, (&RandomConfig{EventsPerPartition: 100}).Validate())
	})
	t.Run("zero_is_invalid", func(t *testing.T) {
		assert.Error(t, (&RandomConfig{EventsPerPartition: 0}).Validate())
	})
	t.Run("negative_is_invalid", func(t *testing.T) {
		assert.Error(t, (&RandomConfig{EventsPerPartition: -1}).Validate())
	})
}

func TestFactory_WrongConfigType(t *testing.T) {
	t.Run("round_robin_rejects_random_config", func(t *testing.T) {
		_, err := NewRoundRobinFactory().CreatePartitioner(&RandomConfig{EventsPerPartition: 1})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RandomConfig")
	})
	t.Run("random_rejects_round_robin_config", func(t *testing.T) {
		_, err := NewRandomFactory().CreatePartitioner(&RoundRobinConfig{EventsPerPartition: 1})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RoundRobinConfig")
	})
}

func TestFactory_DefaultConfig(t *testing.T) {
	t.Run("round_robin_default", func(t *testing.T) {
		cfg := NewRoundRobinFactory().CreateDefaultConfig().(*RoundRobinConfig)
		assert.Equal(t, 1, cfg.EventsPerPartition)
	})
	t.Run("random_default", func(t *testing.T) {
		cfg := NewRandomFactory().CreateDefaultConfig().(*RandomConfig)
		assert.Equal(t, 1, cfg.EventsPerPartition)
	})
}

func TestRoundRobinPartitioner_OneEventPerPartition(t *testing.T) {
	p := mustCreateRoundRobin(t, 1)
	tp := p.ForTopic("test")
	n := 4

	// With eventsPerPartition=1 the sequence must be 0, 1, 2, 3, 0, 1, ...
	want := []int{0, 1, 2, 3, 0, 1, 2, 3}
	for i, w := range want {
		assert.Equal(t, w, tp.Partition(nil, n), "call %d", i+1)
	}
}

func TestRoundRobinPartitioner_MultipleEventsPerPartition(t *testing.T) {
	const n, eventsPerPartition = 3, 3
	p := mustCreateRoundRobin(t, eventsPerPartition)
	tp := p.ForTopic("test")

	// 3 events on 0, then 3 on 1, then 3 on 2, then wraps back to 0.
	want := []int{0, 0, 0, 1, 1, 1, 2, 2, 2, 0, 0, 0}
	for i, w := range want {
		assert.Equal(t, w, tp.Partition(nil, n), "call %d", i+1)
	}
}

func TestRoundRobinPartitioner_WrapsAroundPartitions(t *testing.T) {
	p := mustCreateRoundRobin(t, 1)
	tp := p.ForTopic("test")
	n := 3

	got := make([]int, 9)
	for i := range got {
		got[i] = tp.Partition(nil, n)
	}
	assert.Equal(t, []int{0, 1, 2, 0, 1, 2, 0, 1, 2}, got)
}

func TestRoundRobinPartitioner_IndependentTopics(t *testing.T) {
	p := mustCreateRoundRobin(t, 1)
	tpA := p.ForTopic("topic-a")
	tpB := p.ForTopic("topic-b")
	n := 3

	// Advance topic-a by two steps.
	tpA.Partition(nil, n) // → 0
	tpA.Partition(nil, n) // → 1

	// topic-b must still be at its own independent starting position.
	assert.Equal(t, 0, tpB.Partition(nil, n))
}

func TestRoundRobinPartitioner_RequiresConsistencyFalse(t *testing.T) {
	p := mustCreateRoundRobin(t, 1)
	assert.False(t, p.ForTopic("t").RequiresConsistency(nil))
}

func TestRandomPartitioner_UsesMultiplePartitions(t *testing.T) {
	p := mustCreateRandom(t, 1)
	tp := p.ForTopic("test")
	n := 100

	seen := make(map[int]bool)
	for range 300 {
		seen[tp.Partition(nil, n)] = true
	}
	// With 300 draws over 100 partitions the probability of visiting < 2
	// distinct partitions is negligible.
	assert.Greater(t, len(seen), 1, "random partitioner should spread across partitions")
}

func TestRandomPartitioner_EventsPerPartitionWindowIsRespected(t *testing.T) {
	const eventsPerPartition = 5
	p := mustCreateRandom(t, eventsPerPartition)
	tp := p.ForTopic("test")
	n := 10

	// Collect enough calls to cover several complete windows.
	total := eventsPerPartition * 6
	partitions := make([]int, total)
	for i := range partitions {
		partitions[i] = tp.Partition(nil, n)
	}

	// Every record within the same eventsPerPartition window must share the
	// same partition.
	for windowStart := 0; windowStart < total; windowStart += eventsPerPartition {
		windowEnd := min(windowStart+eventsPerPartition, total)
		for i := windowStart + 1; i < windowEnd; i++ {
			assert.Equal(t, partitions[windowStart], partitions[i],
				"all records in window [%d, %d) must share the same partition", windowStart, windowEnd)
		}
	}
}

func TestRandomPartitioner_IndependentTopics(t *testing.T) {
	p := mustCreateRandom(t, 1)
	n := 100

	tpA := p.ForTopic("topic-a")
	tpB := p.ForTopic("topic-b")

	// Both topic partitioners must be distinct objects with independent state.
	assert.NotSame(t, tpA.(*topicPartitioner), tpB.(*topicPartitioner))
	for range 20 {
		assert.Less(t, tpA.Partition(nil, n), n)
		assert.Less(t, tpB.Partition(nil, n), n)
	}
}

func TestRandomPartitioner_RequiresConsistencyFalse(t *testing.T) {
	p := mustCreateRandom(t, 1)
	assert.False(t, p.ForTopic("t").RequiresConsistency(nil))
}

func mustCreateRoundRobin(t *testing.T, eventsPerPartition int) *roundRobinPartitioner {
	t.Helper()
	p, err := NewRoundRobinFactory().CreatePartitioner(&RoundRobinConfig{EventsPerPartition: eventsPerPartition})
	require.NoError(t, err)
	return p.(*roundRobinPartitioner)
}

func mustCreateRandom(t *testing.T, eventsPerPartition int) *randomPartitioner {
	t.Helper()
	p, err := NewRandomFactory().CreatePartitioner(&RandomConfig{EventsPerPartition: eventsPerPartition})
	require.NoError(t, err)
	return p.(*randomPartitioner)
}
