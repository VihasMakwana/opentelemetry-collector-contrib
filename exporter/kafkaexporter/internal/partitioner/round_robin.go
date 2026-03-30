// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package partitioner

import (
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// RoundRobinConfig configures the round-robin partitioner.
type RoundRobinConfig struct {
	// EventsPerPartition is the number of records sent to the same partition
	// before the partitioner advances to the next one. Must be at least 1.
	EventsPerPartition int `mapstructure:"events_per_partition"`
}

// Validate returns an error if the configuration is invalid.
func (c *RoundRobinConfig) Validate() error {
	if c.EventsPerPartition <= 0 {
		return errors.New("partitioner.events_per_partition must be a positive integer")
	}
	return nil
}

// NewRoundRobinFactory returns a Factory that creates round-robin
// kgo.Partitioner instances.
func NewRoundRobinFactory() Factory {
	return NewFactory(newRoundRobinPartitioner, newRoundRobinDefaultConfig)
}

func newRoundRobinDefaultConfig() Config {
	return &RoundRobinConfig{EventsPerPartition: 1}
}

func newRoundRobinPartitioner(cfg Config) (kgo.Partitioner, error) {
	rCfg, ok := cfg.(*RoundRobinConfig)
	if !ok {
		return nil, fmt.Errorf("expected *RoundRobinConfig, got %T", cfg)
	}
	return &roundRobinPartitioner{n: rCfg.EventsPerPartition}, nil
}

type roundRobinPartitioner struct {
	n int
}

// ForTopic returns a new per-topic TopicPartitioner with its own independent
// state (count and current partition), so different topics cycle through their
// partitions independently.
func (p *roundRobinPartitioner) ForTopic(_ string) kgo.TopicPartitioner {
	n := p.n
	count, partition := 0, 0
	return &topicPartitioner{
		partitionFn: func(_ *kgo.Record, numPartitions int) int {
			if count == n {
				count = 0
				partition = (partition + 1) % numPartitions
			}
			count++
			return partition
		},
	}
}
