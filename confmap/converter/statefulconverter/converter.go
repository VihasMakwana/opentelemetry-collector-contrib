// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package statefulconverter // github.com/open-telemetry/opentelemetry-collector-contrib/confmap/converter/statefulconverter

import (
	"context"
	"fmt"
	"strings"

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

	if !conf.IsSet("receivers") {
		return nil
	}

	receiverCfg, err := conf.Sub("receivers")
	if err != nil {
		return err
	}

	if found := addFilelogStorage(receiverCfg); !found {
		// No stateful receivers found, skipping conversion"
		return nil

	}

	fs := confmap.NewFromStringMap(map[string]any{
		"extensions": map[string]any{
			"file_storage/_stateful": map[string]any{
				"directory": "DEFAULT_PATH",
			},
		},
		"service": map[string]any{
			"extensions": []any{"file_storage/_stateful"},
		},
		// the sub-config of receiver is a copied version from the original. We need to merge it as well.
		"receivers": receiverCfg.ToStringMap(),
	})

	// Merge stateful storage config with the current config
	if err := fs.Merge(conf); err != nil {
		return err
	}

	// Update the input configuration
	return conf.Merge(fs)
}

// addFilelogStorage adds storage config to any stateful receivers if not already set.
func addFilelogStorage(receiverCfg *confmap.Conf) (found bool) {
	for id := range receiverCfg.ToStringMap() {
		typeStr, _, _ := strings.Cut(id, "/")
		typeStr = strings.TrimSpace(typeStr)
		fmt.Println(typeStr, id)
		if typeStr == "" {
			continue
		}
		switch typeStr {
		case "filelog":
			found = true
			if receiverCfg.IsSet(id + "::storage") {
				return
			}
			receiverCfg.Merge(confmap.NewFromStringMap(map[string]any{
				id: map[string]any{
					"storage": "file_storage/_stateful",
				},
			}))
		}
	}
	return
}
