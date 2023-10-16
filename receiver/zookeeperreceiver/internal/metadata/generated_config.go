// Code generated by mdatagen. DO NOT EDIT.

package metadata

import "go.opentelemetry.io/collector/confmap"

// MetricConfig provides common config for a particular metric.
type MetricConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (ms *MetricConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(ms, confmap.WithErrorUnused())
	if err != nil {
		return err
	}
	ms.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// MetricsConfig provides config for zookeeper metrics.
type MetricsConfig struct {
	ZookeeperConnectionActive            MetricConfig `mapstructure:"zookeeper.connection.active"`
	ZookeeperDataTreeEphemeralNodeCount  MetricConfig `mapstructure:"zookeeper.data_tree.ephemeral_node.count"`
	ZookeeperDataTreeSize                MetricConfig `mapstructure:"zookeeper.data_tree.size"`
	ZookeeperFileDescriptorLimit         MetricConfig `mapstructure:"zookeeper.file_descriptor.limit"`
	ZookeeperFileDescriptorOpen          MetricConfig `mapstructure:"zookeeper.file_descriptor.open"`
	ZookeeperFollowerCount               MetricConfig `mapstructure:"zookeeper.follower.count"`
	ZookeeperFsyncExceededThresholdCount MetricConfig `mapstructure:"zookeeper.fsync.exceeded_threshold.count"`
	ZookeeperLatencyAvg                  MetricConfig `mapstructure:"zookeeper.latency.avg"`
	ZookeeperLatencyMax                  MetricConfig `mapstructure:"zookeeper.latency.max"`
	ZookeeperLatencyMin                  MetricConfig `mapstructure:"zookeeper.latency.min"`
	ZookeeperPacketCount                 MetricConfig `mapstructure:"zookeeper.packet.count"`
	ZookeeperRequestActive               MetricConfig `mapstructure:"zookeeper.request.active"`
	ZookeeperRuok                        MetricConfig `mapstructure:"zookeeper.ruok"`
	ZookeeperSyncPending                 MetricConfig `mapstructure:"zookeeper.sync.pending"`
	ZookeeperWatchCount                  MetricConfig `mapstructure:"zookeeper.watch.count"`
	ZookeeperZnodeCount                  MetricConfig `mapstructure:"zookeeper.znode.count"`
}

func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		ZookeeperConnectionActive: MetricConfig{
			Enabled: true,
		},
		ZookeeperDataTreeEphemeralNodeCount: MetricConfig{
			Enabled: true,
		},
		ZookeeperDataTreeSize: MetricConfig{
			Enabled: true,
		},
		ZookeeperFileDescriptorLimit: MetricConfig{
			Enabled: true,
		},
		ZookeeperFileDescriptorOpen: MetricConfig{
			Enabled: true,
		},
		ZookeeperFollowerCount: MetricConfig{
			Enabled: true,
		},
		ZookeeperFsyncExceededThresholdCount: MetricConfig{
			Enabled: true,
		},
		ZookeeperLatencyAvg: MetricConfig{
			Enabled: true,
		},
		ZookeeperLatencyMax: MetricConfig{
			Enabled: true,
		},
		ZookeeperLatencyMin: MetricConfig{
			Enabled: true,
		},
		ZookeeperPacketCount: MetricConfig{
			Enabled: true,
		},
		ZookeeperRequestActive: MetricConfig{
			Enabled: true,
		},
		ZookeeperRuok: MetricConfig{
			Enabled: true,
		},
		ZookeeperSyncPending: MetricConfig{
			Enabled: true,
		},
		ZookeeperWatchCount: MetricConfig{
			Enabled: true,
		},
		ZookeeperZnodeCount: MetricConfig{
			Enabled: true,
		},
	}
}

// ResourceAttributeConfig provides common config for a particular resource attribute.
type ResourceAttributeConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (rac *ResourceAttributeConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(rac, confmap.WithErrorUnused())
	if err != nil {
		return err
	}
	rac.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// ResourceAttributesConfig provides config for zookeeper resource attributes.
type ResourceAttributesConfig struct {
	ServerState ResourceAttributeConfig `mapstructure:"server.state"`
	ZkVersion   ResourceAttributeConfig `mapstructure:"zk.version"`
}

func DefaultResourceAttributesConfig() ResourceAttributesConfig {
	return ResourceAttributesConfig{
		ServerState: ResourceAttributeConfig{
			Enabled: true,
		},
		ZkVersion: ResourceAttributeConfig{
			Enabled: true,
		},
	}
}

// MetricsBuilderConfig is a configuration for zookeeper metrics builder.
type MetricsBuilderConfig struct {
	Metrics            MetricsConfig            `mapstructure:"metrics"`
	ResourceAttributes ResourceAttributesConfig `mapstructure:"resource_attributes"`
}

func DefaultMetricsBuilderConfig() MetricsBuilderConfig {
	return MetricsBuilderConfig{
		Metrics:            DefaultMetricsConfig(),
		ResourceAttributes: DefaultResourceAttributesConfig(),
	}
}
