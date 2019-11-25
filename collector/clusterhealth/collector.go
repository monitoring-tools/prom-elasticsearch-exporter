package clusterhealth

import (
	"log"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Cluster statuses
var clusterStatuses = map[string]float64{
	"green":  1,
	"yellow": 2,
	"red":    3,
}

var (
	labelsClusterHealth = []string{"cluster"}
	labelsIndexHealth   = []string{"cluster", "index"}

	subsystemClusterHealth = "cluster_health"
	subsystemIndexHealth   = "cluster_health_index"
)

type clusterHealthMetric struct {
	*metrics.Metric
	Value func(clusterHealth *model.ClusterHealth) float64
}

type indexHealthMetric struct {
	*metrics.Metric
	Value func(indexHealth model.ClusterHealthIndex) float64
}

type clusterHealthStatusMetric struct {
	*metrics.Metric
	Value  func(clusterHealth *model.ClusterHealth) float64
	Labels func(clusterName, color string) []string
}

// Collector is an cluster health metrics collector
type Collector struct {
	esClient elasticsearch.Client

	metrics        []*clusterHealthMetric
	statusMetric   *clusterHealthStatusMetric
	indicesMetrics []*indexHealthMetric
}

func newClusterHealthMetric(name, help string, valueExtractor func(*model.ClusterHealth) float64) *clusterHealthMetric {
	return &clusterHealthMetric{
		Metric: metrics.New(prometheus.GaugeValue, subsystemClusterHealth, name, help, labelsClusterHealth),
		Value:  valueExtractor,
	}
}

func newIndexHealthMetric(name, help string, valueExtractor func(model.ClusterHealthIndex) float64) *indexHealthMetric {
	return &indexHealthMetric{
		Metric: metrics.New(prometheus.GaugeValue, subsystemIndexHealth, name, help, labelsIndexHealth),
		Value:  valueExtractor,
	}
}

// NewCollector returns new cluster health collector
func NewCollector(esClient elasticsearch.Client) *Collector {
	return &Collector{
		esClient: esClient,

		metrics: []*clusterHealthMetric{
			newClusterHealthMetric(
				"active_primary_shards", "The number of primary shards in your cluster. This is an aggregate total across all indices.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.ActivePrimaryShards) },
			),
			newClusterHealthMetric(
				"active_shards", "Aggregate total of all shards across all indices, which includes replica shards.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.ActiveShards) },
			),
			newClusterHealthMetric(
				"delayed_unassigned_shards", "Shards delayed to reduce reallocation overhead",
				func(clusterHealth *model.ClusterHealth) float64 {
					return float64(clusterHealth.DelayedUnassignedShards)
				},
			),
			newClusterHealthMetric(
				"initializing_shards", "Count of shards that are being freshly created.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.InitializingShards) },
			),
			newClusterHealthMetric(
				"number_of_data_nodes", "Number of data nodes in the cluster.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.NumberOfDataNodes) },
			),
			newClusterHealthMetric(
				"number_of_in_flight_fetch", "The number of ongoing shard info requests.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.NumberOfInFlightFetch) },
			),
			newClusterHealthMetric(
				"number_of_nodes", "Number of nodes in the cluster.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.NumberOfNodes) },
			),
			newClusterHealthMetric(
				"number_of_pending_tasks", "Cluster level changes which have not yet been executed",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.NumberOfPendingTasks) },
			),
			newClusterHealthMetric(
				"relocating_shards", "The number of shards that are currently moving from one node to another node.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.RelocatingShards) },
			),
			newClusterHealthMetric(
				"timed_out", "Number of cluster health checks timed out",
				func(clusterHealth *model.ClusterHealth) float64 {
					if clusterHealth.TimedOut {
						return 1
					}
					return 0
				},
			),
			newClusterHealthMetric(
				"unassigned_shards", "The number of shards that exist in the cluster state, but cannot be found in the cluster itself.",
				func(clusterHealth *model.ClusterHealth) float64 { return float64(clusterHealth.UnassignedShards) },
			),
		},
		statusMetric: &clusterHealthStatusMetric{
			Metric: metrics.New(prometheus.GaugeValue, subsystemClusterHealth, "status", "Cluster status. 1 = green, 2 = yellow, 3 = red", labelsClusterHealth),
			Value: func(clusterHealth *model.ClusterHealth) float64 {
				return clusterStatuses[clusterHealth.Status]
			},
		},
		indicesMetrics: []*indexHealthMetric{
			newIndexHealthMetric(
				"status", "Index status. 1 = green, 2 = yellow, 3 = red",
				func(i model.ClusterHealthIndex) float64 { return clusterStatuses[i.Status] },
			),
			newIndexHealthMetric(
				"number_of_shards", "The number of shards that used by index",
				func(i model.ClusterHealthIndex) float64 { return float64(i.NumberOfShards) },
			),
			newIndexHealthMetric(
				"number_of_replicas", "The number of replicas of index",
				func(i model.ClusterHealthIndex) float64 { return float64(i.NumberOfReplicas) },
			),
			newIndexHealthMetric(
				"active_primary_shards", "The number of active primary shards of index",
				func(i model.ClusterHealthIndex) float64 { return float64(i.ActivePrimaryShards) },
			),
			newIndexHealthMetric(
				"active_shards", "The number of active shards of index",
				func(i model.ClusterHealthIndex) float64 { return float64(i.ActiveShards) },
			),
			newIndexHealthMetric(
				"relocating_shards", "The number of relocating shards of index",
				func(i model.ClusterHealthIndex) float64 { return float64(i.RelocatingShards) },
			),
			newIndexHealthMetric(
				"initializing_shards", "The number of initializing shards of index",
				func(i model.ClusterHealthIndex) float64 { return float64(i.InitializingShards) },
			),
			newIndexHealthMetric(
				"unassigned_shards", "The number of unassigned shards of index",
				func(i model.ClusterHealthIndex) float64 { return float64(i.UnassignedShards) },
			),
		},
	}
}

// Describe implements prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.Desc()
	}

	for _, metric := range c.indicesMetrics {
		ch <- metric.Desc()
	}

	ch <- c.statusMetric.Desc()
}

// Collect writes data to metrics channel
func (c *Collector) Collect(clusterName string, ch chan<- prometheus.Metric) {
	resp, err := c.esClient.ClusterHealth(elasticsearch.LevelIndices)
	if err != nil {
		log.Println("ERROR: failed to fetch cluster health: ", err)
		return
	}

	for _, metric := range c.metrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc(),
			metric.Type(),
			metric.Value(resp),
			clusterName,
		)
	}

	for name, index := range resp.Indices {
		for _, metric := range c.indicesMetrics {
			ch <- prometheus.MustNewConstMetric(
				metric.Desc(),
				metric.Type(),
				metric.Value(index),
				clusterName, name,
			)
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.statusMetric.Desc(),
		c.statusMetric.Type(),
		c.statusMetric.Value(resp),
		clusterName,
	)
}
