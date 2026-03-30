// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package partitioner // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter/internal/partitioner"

import "github.com/twmb/franz-go/pkg/kgo"

// Config is the interface implemented by all partitioner-specific
// configuration structs.
type Config interface {
	Validate() error
}

// Factory creates a configured kgo.Partitioner and exposes its default
// configuration.
type Factory interface {
	// CreatePartitioner returns a new kgo.Partitioner for the given Config.
	// cfg must be the concrete type expected by the factory; passing an
	// incompatible type returns an error.
	CreatePartitioner(cfg Config) (kgo.Partitioner, error)
	// CreateDefaultConfig returns a new Config instance pre-populated with
	// sensible defaults.
	CreateDefaultConfig() Config
}

// NewFactory constructs a Factory from a partitioner-creation function and a
// default-config function.
func NewFactory(createFn func(Config) (kgo.Partitioner, error), defaultFn func() Config) Factory {
	return &factoryImpl{createFn: createFn, defaultFn: defaultFn}
}

type factoryImpl struct {
	createFn  func(Config) (kgo.Partitioner, error)
	defaultFn func() Config
}

func (f *factoryImpl) CreatePartitioner(cfg Config) (kgo.Partitioner, error) {
	return f.createFn(cfg)
}

func (f *factoryImpl) CreateDefaultConfig() Config {
	return f.defaultFn()
}

// topicPartitioner is a reusable kgo.TopicPartitioner that delegates
// partition selection to a stateful closure.  The closure is created fresh
// for every topic (see ForTopic), so each topic maintains independent state.
type topicPartitioner struct {
	partitionFn func(r *kgo.Record, numPartitions int) int
}

func (t *topicPartitioner) RequiresConsistency(*kgo.Record) bool { return false }

func (t *topicPartitioner) Partition(r *kgo.Record, n int) int {
	return t.partitionFn(r, n)
}
