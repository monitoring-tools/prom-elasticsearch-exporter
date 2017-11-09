package aliases

import (
	"log"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector is a metrics collection for ElasticSearch indices aliases
type Collector struct {
	esClient elasticsearch.IClient

	metrics []*metrics.Metric
}

// NewCollector returns new metrics collector for index aliases
func NewCollector(esClient elasticsearch.IClient) *Collector {
	return &Collector{
		esClient: esClient,
		metrics: []*metrics.Metric{
			metrics.New(
				prometheus.GaugeValue,
				"indices",
				"alias",
				"Static metric with index name, alias and const value = 1",
				[]string{"cluster", "index", "alias"},
			),
		},
	}
}

// Describe implements prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.Desc()
	}
}

// Collect writes data to metrics channel
func (c *Collector) Collect(clusterName string, ch chan<- prometheus.Metric) {
	indices, err := c.esClient.Aliases()
	if err != nil {
		log.Println("ERROR: failed to fetch aliases: ", err)
		return
	}

	for _, metric := range c.metrics {
		for indexName, aliases := range indices {
			for aliasName := range aliases.Aliases {
				ch <- prometheus.MustNewConstMetric(metric.Desc(), metric.Type(), 1.0, clusterName, indexName, aliasName)
			}
		}
	}
}
