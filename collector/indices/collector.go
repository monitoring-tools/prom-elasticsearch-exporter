package indices

import (
	"log"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labelsIndex = []string{"cluster", "index"}
)

// Collector is a metrics collection with ElasticSearch indices stats
type Collector struct {
	esClient elasticsearch.IClient

	totalMetrics     []*indexMetric
	primariesMetrics []*indexMetric
}

type indexMetric struct {
	*metrics.Metric

	Value       func(model.IndexSummary) float64
	LabelValues func(cluster, index string) []string
}

type indexMetricTemplate struct {
	Type           prometheus.ValueType
	Name           string
	Help           string
	ValueExtractor func(model.IndexSummary) float64
}

func newIndexMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.IndexSummary) float64) *indexMetricTemplate {
	return &indexMetricTemplate{
		Type:           t,
		Name:           name,
		Help:           help,
		ValueExtractor: valueExtractor,
	}
}

// NewCollector returns new metrics collection for indices metrics
func NewCollector(esClient elasticsearch.IClient) *Collector {
	var indexMetricTemplates = []*indexMetricTemplate{
		newIndexMetric(
			prometheus.GaugeValue, "docs_count", "Docs count",
			func(i model.IndexSummary) float64 { return float64(i.Docs.Count) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "docs_deleted", "Docs deleted",
			func(i model.IndexSummary) float64 { return float64(i.Docs.Deleted) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "store_size_bytes", "The size of the store for shards",
			func(i model.IndexSummary) float64 { return float64(i.Store.SizeInBytes) },
		),
		newIndexMetric(
			prometheus.CounterValue, "search_query_time_seconds", "Total search query time in seconds",
			func(i model.IndexSummary) float64 { return float64(i.Search.QueryTimeInMillis / 1000) },
		),
		newIndexMetric(
			prometheus.CounterValue, "search_query_total", "Total number of search queries",
			func(i model.IndexSummary) float64 { return float64(i.Search.QueryTotal) },
		),
		newIndexMetric(
			prometheus.CounterValue, "search_fetch_time_seconds", "Total search fetch time in seconds",
			func(i model.IndexSummary) float64 { return float64(i.Search.FetchTimeInMillis / 1000) },
		),
		newIndexMetric(
			prometheus.CounterValue, "search_fetch_total", "Total number of fetches",
			func(i model.IndexSummary) float64 { return float64(i.Search.FetchTotal) },
		),
		newIndexMetric(
			prometheus.CounterValue, "indexing_index_total", "Total index calls",
			func(i model.IndexSummary) float64 { return float64(i.Indexing.IndexTotal) },
		),
		newIndexMetric(
			prometheus.CounterValue, "indexing_index_seconds_total", "Cumulative indexing time in seconds",
			func(i model.IndexSummary) float64 { return float64(i.Indexing.IndexTimeInMillis / 1000) },
		),
		newIndexMetric(
			prometheus.CounterValue, "indexing_throttle_seconds_total", "Cumulative throttle time in seconds",
			func(i model.IndexSummary) float64 { return float64(i.Indexing.ThrottleTimeInMillis / 1000) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "segments_count", "Number of segments",
			func(i model.IndexSummary) float64 { return float64(i.Segments.Count) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "segments_memory_bytes", "Segments memory in bytes",
			func(i model.IndexSummary) float64 { return float64(i.Segments.MemoryInBytes) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "query_cache_memory_size_bytes", "Query cache memory usage in bytes",
			func(i model.IndexSummary) float64 { return float64(i.QueryCache.MemorySizeInBytes) },
		),
		newIndexMetric(
			prometheus.CounterValue, "query_cache_evictions", "Total evictions number from query cache",
			func(i model.IndexSummary) float64 { return float64(i.QueryCache.Evictions) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "request_cache_memory_size_bytes", "Request cache memory usage in bytes",
			func(i model.IndexSummary) float64 { return float64(i.RequestCache.MemorySizeInBytes) },
		),
		newIndexMetric(
			prometheus.CounterValue, "request_cache_evictions", "Total evictions number from request cache",
			func(i model.IndexSummary) float64 { return float64(i.RequestCache.Evictions) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "fielddata_memory_size_bytes", "Fielddata memory usage in bytes",
			func(i model.IndexSummary) float64 { return float64(i.Fielddata.MemorySizeInBytes) },
		),
		newIndexMetric(
			prometheus.CounterValue, "fielddata_evictions", "Total evictions number from fielddata",
			func(i model.IndexSummary) float64 { return float64(i.Fielddata.Evictions) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "segments_index_writer_memory_size_bytes", "Index writer memory usage",
			func(i model.IndexSummary) float64 { return float64(i.Segments.IndexWriterMemoryInBytes) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "merges_size_bytes", "Merges total size in bytes",
			func(i model.IndexSummary) float64 { return float64(i.Merges.TotalSizeInBytes) },
		),
		newIndexMetric(
			prometheus.CounterValue, "refresh_total", "Total refresh calls",
			func(i model.IndexSummary) float64 { return float64(i.Refresh.Total) },
		),
		newIndexMetric(
			prometheus.CounterValue, "refresh_time_seconds", "Total refresh time in seconds",
			func(i model.IndexSummary) float64 { return float64(i.Refresh.TotalTimeInMillis / 1000) },
		),
		newIndexMetric(
			prometheus.CounterValue, "translog_operations", "Total translog operations",
			func(i model.IndexSummary) float64 { return float64(i.Translog.Operations) },
		),
		newIndexMetric(
			prometheus.GaugeValue, "translog_size_in_bytes", "Transolog size in bytes",
			func(i model.IndexSummary) float64 { return float64(i.Translog.SizeInBytes) },
		),
	}

	var subsystem = "index"

	labelValuesExtractor := func(cluster, index string) []string {
		return []string{cluster, index}
	}

	primariesMetrics := make([]*indexMetric, len(indexMetricTemplates))
	totalMetrics := make([]*indexMetric, len(indexMetricTemplates))

	for i, m := range indexMetricTemplates {
		primariesMetrics[i] = &indexMetric{
			Metric:      metrics.New(m.Type, subsystem, "primaries_"+m.Name, m.Help, labelsIndex),
			Value:       m.ValueExtractor,
			LabelValues: labelValuesExtractor,
		}
		totalMetrics[i] = &indexMetric{
			Metric:      metrics.New(m.Type, subsystem, "total_"+m.Name, m.Help, labelsIndex),
			Value:       m.ValueExtractor,
			LabelValues: labelValuesExtractor,
		}
	}

	return &Collector{
		esClient:         esClient,
		primariesMetrics: primariesMetrics,
		totalMetrics:     totalMetrics,
	}
}

// Describe implements prometheus.Collector interface
func (i *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range i.primariesMetrics {
		ch <- metric.Desc()
	}
	for _, metric := range i.totalMetrics {
		ch <- metric.Desc()
	}
}

// Collect writes data to metrics channel
func (i *Collector) Collect(clusterName string, ch chan<- prometheus.Metric) {
	res, err := i.esClient.Indices()
	if err != nil {
		log.Println("ERROR: failed to fetch indices stats: ", err)
		return
	}

	for indexName, index := range res.Indices {
		for _, metric := range i.primariesMetrics {
			ch <- prometheus.MustNewConstMetric(
				metric.Desc(),
				metric.Type(),
				metric.Value(index.Primaries),
				metric.LabelValues(clusterName, indexName)...,
			)
		}
		for _, metric := range i.totalMetrics {
			ch <- prometheus.MustNewConstMetric(
				metric.Desc(),
				metric.Type(),
				metric.Value(index.Total),
				metric.LabelValues(clusterName, indexName)...,
			)
		}
	}
}
