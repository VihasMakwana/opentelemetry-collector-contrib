// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestMetricsBuilderConfig(t *testing.T) {
	tests := []struct {
		name string
		want MetricsBuilderConfig
	}{
		{
			name: "default",
			want: DefaultMetricsBuilderConfig(),
		},
		{
			name: "all_set",
			want: MetricsBuilderConfig{
				Metrics: MetricsConfig{
					VcenterClusterCPUEffective:      MetricConfig{Enabled: true},
					VcenterClusterCPULimit:          MetricConfig{Enabled: true},
					VcenterClusterHostCount:         MetricConfig{Enabled: true},
					VcenterClusterMemoryEffective:   MetricConfig{Enabled: true},
					VcenterClusterMemoryLimit:       MetricConfig{Enabled: true},
					VcenterClusterMemoryUsed:        MetricConfig{Enabled: true},
					VcenterClusterVMCount:           MetricConfig{Enabled: true},
					VcenterDatastoreDiskUsage:       MetricConfig{Enabled: true},
					VcenterDatastoreDiskUtilization: MetricConfig{Enabled: true},
					VcenterHostCPUUsage:             MetricConfig{Enabled: true},
					VcenterHostCPUUtilization:       MetricConfig{Enabled: true},
					VcenterHostDiskLatencyAvg:       MetricConfig{Enabled: true},
					VcenterHostDiskLatencyMax:       MetricConfig{Enabled: true},
					VcenterHostDiskThroughput:       MetricConfig{Enabled: true},
					VcenterHostMemoryUsage:          MetricConfig{Enabled: true},
					VcenterHostMemoryUtilization:    MetricConfig{Enabled: true},
					VcenterHostNetworkPacketCount:   MetricConfig{Enabled: true},
					VcenterHostNetworkPacketErrors:  MetricConfig{Enabled: true},
					VcenterHostNetworkThroughput:    MetricConfig{Enabled: true},
					VcenterHostNetworkUsage:         MetricConfig{Enabled: true},
					VcenterResourcePoolCPUShares:    MetricConfig{Enabled: true},
					VcenterResourcePoolCPUUsage:     MetricConfig{Enabled: true},
					VcenterResourcePoolMemoryShares: MetricConfig{Enabled: true},
					VcenterResourcePoolMemoryUsage:  MetricConfig{Enabled: true},
					VcenterVMCPUUsage:               MetricConfig{Enabled: true},
					VcenterVMCPUUtilization:         MetricConfig{Enabled: true},
					VcenterVMDiskLatencyAvg:         MetricConfig{Enabled: true},
					VcenterVMDiskLatencyMax:         MetricConfig{Enabled: true},
					VcenterVMDiskThroughput:         MetricConfig{Enabled: true},
					VcenterVMDiskUsage:              MetricConfig{Enabled: true},
					VcenterVMDiskUtilization:        MetricConfig{Enabled: true},
					VcenterVMMemoryBallooned:        MetricConfig{Enabled: true},
					VcenterVMMemorySwapped:          MetricConfig{Enabled: true},
					VcenterVMMemorySwappedSsd:       MetricConfig{Enabled: true},
					VcenterVMMemoryUsage:            MetricConfig{Enabled: true},
					VcenterVMMemoryUtilization:      MetricConfig{Enabled: true},
					VcenterVMNetworkPacketCount:     MetricConfig{Enabled: true},
					VcenterVMNetworkThroughput:      MetricConfig{Enabled: true},
					VcenterVMNetworkUsage:           MetricConfig{Enabled: true},
				},
				ResourceAttributes: ResourceAttributesConfig{
					VcenterClusterName:      ResourceAttributeConfig{Enabled: true},
					VcenterDatastoreName:    ResourceAttributeConfig{Enabled: true},
					VcenterHostName:         ResourceAttributeConfig{Enabled: true},
					VcenterResourcePoolName: ResourceAttributeConfig{Enabled: true},
					VcenterVMID:             ResourceAttributeConfig{Enabled: true},
					VcenterVMName:           ResourceAttributeConfig{Enabled: true},
				},
			},
		},
		{
			name: "none_set",
			want: MetricsBuilderConfig{
				Metrics: MetricsConfig{
					VcenterClusterCPUEffective:      MetricConfig{Enabled: false},
					VcenterClusterCPULimit:          MetricConfig{Enabled: false},
					VcenterClusterHostCount:         MetricConfig{Enabled: false},
					VcenterClusterMemoryEffective:   MetricConfig{Enabled: false},
					VcenterClusterMemoryLimit:       MetricConfig{Enabled: false},
					VcenterClusterMemoryUsed:        MetricConfig{Enabled: false},
					VcenterClusterVMCount:           MetricConfig{Enabled: false},
					VcenterDatastoreDiskUsage:       MetricConfig{Enabled: false},
					VcenterDatastoreDiskUtilization: MetricConfig{Enabled: false},
					VcenterHostCPUUsage:             MetricConfig{Enabled: false},
					VcenterHostCPUUtilization:       MetricConfig{Enabled: false},
					VcenterHostDiskLatencyAvg:       MetricConfig{Enabled: false},
					VcenterHostDiskLatencyMax:       MetricConfig{Enabled: false},
					VcenterHostDiskThroughput:       MetricConfig{Enabled: false},
					VcenterHostMemoryUsage:          MetricConfig{Enabled: false},
					VcenterHostMemoryUtilization:    MetricConfig{Enabled: false},
					VcenterHostNetworkPacketCount:   MetricConfig{Enabled: false},
					VcenterHostNetworkPacketErrors:  MetricConfig{Enabled: false},
					VcenterHostNetworkThroughput:    MetricConfig{Enabled: false},
					VcenterHostNetworkUsage:         MetricConfig{Enabled: false},
					VcenterResourcePoolCPUShares:    MetricConfig{Enabled: false},
					VcenterResourcePoolCPUUsage:     MetricConfig{Enabled: false},
					VcenterResourcePoolMemoryShares: MetricConfig{Enabled: false},
					VcenterResourcePoolMemoryUsage:  MetricConfig{Enabled: false},
					VcenterVMCPUUsage:               MetricConfig{Enabled: false},
					VcenterVMCPUUtilization:         MetricConfig{Enabled: false},
					VcenterVMDiskLatencyAvg:         MetricConfig{Enabled: false},
					VcenterVMDiskLatencyMax:         MetricConfig{Enabled: false},
					VcenterVMDiskThroughput:         MetricConfig{Enabled: false},
					VcenterVMDiskUsage:              MetricConfig{Enabled: false},
					VcenterVMDiskUtilization:        MetricConfig{Enabled: false},
					VcenterVMMemoryBallooned:        MetricConfig{Enabled: false},
					VcenterVMMemorySwapped:          MetricConfig{Enabled: false},
					VcenterVMMemorySwappedSsd:       MetricConfig{Enabled: false},
					VcenterVMMemoryUsage:            MetricConfig{Enabled: false},
					VcenterVMMemoryUtilization:      MetricConfig{Enabled: false},
					VcenterVMNetworkPacketCount:     MetricConfig{Enabled: false},
					VcenterVMNetworkThroughput:      MetricConfig{Enabled: false},
					VcenterVMNetworkUsage:           MetricConfig{Enabled: false},
				},
				ResourceAttributes: ResourceAttributesConfig{
					VcenterClusterName:      ResourceAttributeConfig{Enabled: false},
					VcenterDatastoreName:    ResourceAttributeConfig{Enabled: false},
					VcenterHostName:         ResourceAttributeConfig{Enabled: false},
					VcenterResourcePoolName: ResourceAttributeConfig{Enabled: false},
					VcenterVMID:             ResourceAttributeConfig{Enabled: false},
					VcenterVMName:           ResourceAttributeConfig{Enabled: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := loadMetricsBuilderConfig(t, tt.name)
			if diff := cmp.Diff(tt.want, cfg, cmpopts.IgnoreUnexported(MetricConfig{}, ResourceAttributeConfig{})); diff != "" {
				t.Errorf("Config mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}

func loadMetricsBuilderConfig(t *testing.T, name string) MetricsBuilderConfig {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	sub, err := cm.Sub(name)
	require.NoError(t, err)
	cfg := DefaultMetricsBuilderConfig()
	require.NoError(t, component.UnmarshalConfig(sub, &cfg))
	return cfg
}
