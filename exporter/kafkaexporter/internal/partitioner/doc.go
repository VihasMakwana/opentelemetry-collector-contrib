// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package partitioner provides franz-go kgo.Partitioner implementations for
// the Kafka exporter.
// Each partitioner is created through a Factory, which also exposes a default Config.
// Use NewRoundRobinFactory or NewRandomFactory to obtain the corresponding factory.
package partitioner // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter/internal/partitioner"
