package recovery

import (
	"log"
	"strconv"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Recovery stages mapping
var recoveryStages = map[string]float64{
	"INIT":     1,
	"INDEX":    2,
	"START":    3,
	"TRANSLOG": 4,
	"FINALIZE": 5,
	"DONE":     6,
}

type recoveryMetric struct {
	*metrics.Metric

	Value       func(model.ShardRecovery) float64
	LabelValues func(cluster, index string, shard model.ShardRecovery) []string
}

func newMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.ShardRecovery) float64) *recoveryMetric {
	return &recoveryMetric{
		Metric: metrics.New(
			t,
			"index_recovery",
			name,
			help,
			[]string{"cluster", "index", "shard_id", "src_name", "target_name"},
		),

		Value: valueExtractor,
		LabelValues: func(cluster, index string, shard model.ShardRecovery) []string {
			return []string{cluster, index, strconv.Itoa(int(shard.ID)), shard.Source.Name, shard.Target.Name}
		},
	}
}

var recoveryShardInfo = metrics.NewDesc(
	"index_recovery",
	"info",
	"Index recovery info",
	[]string{"cluster", "index", "shard_id", "type", "is_primary", "src_name", "src_ip", "target_name", "target_ip"},
	nil,
)

// Collector is a metrics collection for ElasticSearch indices aliases
type Collector struct {
	esClient elasticsearch.IClient

	recoveryInfo prometheus.Metric
	metrics      []*recoveryMetric
}

// NewCollector returns new metrics collection for indices aliases
func NewCollector(esClient elasticsearch.IClient) *Collector {
	return &Collector{
		esClient: esClient,
		metrics: []*recoveryMetric{
			newMetric(
				prometheus.GaugeValue, "bytes_total", "Total size of index shard in bytes",
				func(s model.ShardRecovery) float64 {
					return float64(s.Index.Size.TotalInBytes)
				},
			),
			newMetric(
				prometheus.GaugeValue, "bytes_recovered", "Size of recovered data in bytes",
				func(s model.ShardRecovery) float64 {
					return float64(s.Index.Size.RecoveredInBytes)
				},
			),
			newMetric(
				prometheus.GaugeValue, "files_total", "Number of files in index shard",
				func(s model.ShardRecovery) float64 {
					return float64(s.Index.Files.Total)
				},
			),
			newMetric(
				prometheus.GaugeValue, "files_recovered", "Number of recovered files in index shard",
				func(s model.ShardRecovery) float64 {
					return float64(s.Index.Files.Recovered)
				},
			),
			newMetric(
				prometheus.GaugeValue, "translog_total", "Total size of translog",
				func(s model.ShardRecovery) float64 {
					return float64(s.Translog.Total)
				},
			),
			newMetric(
				prometheus.GaugeValue, "translog_recovered", "Total size of recovered translog",
				func(s model.ShardRecovery) float64 {
					return float64(s.Translog.Recovered)
				},
			),
			newMetric(
				prometheus.GaugeValue, "stage", "Index shard recovery stage. 1 = INIT, 2 = INDEX, 3 = START, 4 = TRANSLOG, 5 = FINALIZE, 6 = DONE.",
				func(s model.ShardRecovery) float64 {
					if v, ok := recoveryStages[s.Stage]; ok {
						return v
					}
					return 0
				},
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
	indicesRecovery, err := c.esClient.Recovery()
	if err != nil {
		log.Println("ERROR: failed to fetch recovery stats: ", err)
		return
	}

	for indexName, index := range indicesRecovery {
		for _, shard := range index.Shards {
			isPrimary := "false"
			if shard.Primary {
				isPrimary = "true"
			}

			infoMetric, err := prometheus.NewConstMetric(
				recoveryShardInfo,
				prometheus.GaugeValue,
				1,
				clusterName, indexName, strconv.Itoa(int(shard.ID)), shard.Type, isPrimary, shard.Source.Name, shard.Source.IP, shard.Target.Name, shard.Target.IP,
			)

			if err == nil {
				ch <- infoMetric
			} else {
				log.Println("Can't create recovery info recoveryMetric:", err)
			}

			for _, metric := range c.metrics {
				m, err := prometheus.NewConstMetric(
					metric.Desc(),
					metric.Type(),
					metric.Value(shard),
					metric.LabelValues(clusterName, indexName, shard)...,
				)

				if err != nil {
					log.Println("Can't create recovery recoveryMetric:", err)
					continue
				}

				ch <- m
			}
		}
	}
}
