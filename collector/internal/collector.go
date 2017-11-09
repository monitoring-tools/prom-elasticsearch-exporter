package internal

import (
	"github.com/monitoring-tools/prom-elasticsearch-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector is a collector for internal exporter's metrics
type Collector struct {
	buildInfoMetric *metrics.Metric
}

// NewCollector returns new metrics collector for index aliases
func NewCollector(appVersion, goVersion, gitBranch string) *Collector {
	if appVersion == "" {
		appVersion = "n/a"
	}
	if goVersion == "" {
		goVersion = "n/a"
	}
	if gitBranch == "" {
		gitBranch = "n/a"
	}

	return &Collector{
		buildInfoMetric: metrics.NewWithConstLabels(
			prometheus.GaugeValue,
			"exporter",
			"build_info",
			"Static buildInfoMetric with exporter build info",
			nil,
			prometheus.Labels{
				"version":   appVersion,
				"goversion": goVersion,
				"branch":    gitBranch,
			},
		),
	}
}

// Describe implements prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.buildInfoMetric.Desc()
}

// Collect writes data to metrics channel
func (c *Collector) Collect(clusterName string, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.buildInfoMetric.Desc(),
		c.buildInfoMetric.Type(),
		1,
	)
}
