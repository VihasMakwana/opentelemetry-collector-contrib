// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package partitioner

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/twmb/franz-go/pkg/kgo"
)

// RandomConfig configures the random partitioner.
type RandomConfig struct {
	// EventsPerPartition is the number of records sent to the same partition
	// before randomly selecting the next one. Must be at least 1.
	EventsPerPartition int `mapstructure:"events_per_partition"`
}

// Validate returns an error if the configuration is invalid.
func (c *RandomConfig) Validate() error {
	if c.EventsPerPartition <= 0 {
		return errors.New("partitioner.events_per_partition must be a positive integer")
	}
	return nil
}

// NewRandomFactory returns a Factory that creates random kgo.Partitioner
// instances.
func NewRandomFactory() Factory {
	return NewFactory(newRandomPartitioner, newRandomDefaultConfig)
}

func newRandomDefaultConfig() Config {
	return &RandomConfig{EventsPerPartition: 1}
}

func newRandomPartitioner(cfg Config) (kgo.Partitioner, error) {
	rCfg, ok := cfg.(*RandomConfig)
	if !ok {
		return nil, fmt.Errorf("expected *RandomConfig, got %T", cfg)
	}
	return &randomPartitioner{n: rCfg.EventsPerPartition}, nil
}

type randomPartitioner struct {
	n int
}

// ForTopic returns a new per-topic TopicPartitioner with its own independent
// state (count and current partition).  The initial partition is chosen at
// random, and a new random partition is selected after every n records.
func (p *randomPartitioner) ForTopic(_ string) kgo.TopicPartitioner {
	n := p.n
	count, partition := 0, -1
	return &topicPartitioner{
		partitionFn: func(_ *kgo.Record, numPartitions int) int {
			if count == n || partition < 0 {
				count = 0
				partition = rand.IntN(numPartitions)
			}
			count++
			return partition
		},
	}
}
