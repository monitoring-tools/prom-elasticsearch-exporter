package collector

import (
	"log"
	"sync"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/collector/aliases"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/collector/clusterhealth"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/collector/indices"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/collector/internal"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/collector/nodes"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/collector/recovery"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/prometheus/client_golang/prometheus"
)

// ICollector is a metrics collector interface
type ICollector interface {
	Describe(ch chan<- *prometheus.Desc)
	Collect(clusterName string, ch chan<- prometheus.Metric)
}

// CompositeCollector collects all ES metrics: cluster, nodes, indices.
// Implements prometheus.Collector
type CompositeCollector struct {
	esClient   elasticsearch.IClient
	collectors []ICollector
}

// NewCompositeCollector creates new composite collector
func NewCompositeCollector(esClient elasticsearch.IClient, exportMetricsForAllNodes bool, appVersion, goVersion, gitBranch string) *CompositeCollector {
	collectors := []ICollector{
		internal.NewCollector(appVersion, goVersion, gitBranch),
		clusterhealth.NewCollector(esClient),
		nodes.NewCollector(esClient, exportMetricsForAllNodes),
		aliases.NewCollector(esClient),
		indices.NewCollector(esClient),
		recovery.NewCollector(esClient),
	}

	return &CompositeCollector{
		esClient:   esClient,
		collectors: collectors,
	}
}

// Describe sends the super-set of all possible descriptors of metrics
func (c *CompositeCollector) Describe(ch chan<- *prometheus.Desc) {
	c.forEachCollector(func(collector ICollector) {
		collector.Describe(ch)
	})
}

// Collect is called by the Prometheus registry when collecting metrics
func (c *CompositeCollector) Collect(ch chan<- prometheus.Metric) {
	clusterHealth, err := c.esClient.ClusterHealth(elasticsearch.LevelCluster)
	if err != nil {
		log.Println("ERROR: can't fetch cluster name: ", err)
		return
	}

	c.forEachCollector(func(collector ICollector) {
		collector.Collect(clusterHealth.ClusterName, ch)
	})
}

// runCollectorsConcurrently runs given
func (c *CompositeCollector) forEachCollector(fn func(ICollector)) {
	var group sync.WaitGroup
	group.Add(len(c.collectors))

	for _, collector := range c.collectors {
		go func(collector ICollector) {
			fn(collector)
			group.Done()
		}(collector)
	}

	group.Wait()
}
