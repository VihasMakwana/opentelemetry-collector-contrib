// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package statefulconverter // github.com/open-telemetry/opentelemetry-collector-contrib/confmap/converter/statefulconverter

import (
	"context"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/featuregate"
)

var Stateful = featuregate.GlobalRegistry().MustRegister(
	"stateful",
	featuregate.StageAlpha,
)

type converter struct{}

func NewFactory() confmap.ConverterFactory {
	return confmap.NewConverterFactory(newConverter)
}

func newConverter(set confmap.ConverterSettings) confmap.Converter {
	return converter{}
}

func (converter) Convert(_ context.Context, conf *confmap.Conf) error {
	if !Stateful.IsEnabled() {
		return nil
	}
	featuregate.GlobalRegistry().Set("confmap.enableMergeAppendOption", true)
	defer featuregate.GlobalRegistry().Set("confmap.enableMergeAppendOption", false)

	fs := map[string]any{
		"extensions": map[string]any{
			"file_storage": map[string]any{
				"directory": "/Users/vihasmakwana/Desktop/Vihas/OTeL/stateful-collector",
			},
		},
		"service": map[string]any{
			"extensions": []any{"file_storage"},
		},
	}
	return confmap.NewFromStringMap(fs).Merge(conf)
}
