package nodes

import (
	"log"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Labels lists for different kind of metrics
var (
	labelsNode       = []string{"cluster", "host", "node"}
	labelsThreadPool = append(labelsNode, "type")
	labelsBreaker    = append(labelsNode, "breaker")
	labelsFilesystem = append(labelsNode, "mount", "path")
	labelsJVMGC      = append(labelsNode, "gc")
)

// Label values extractors for different kind of metrics
var (
	labelValuesNode = func(cluster string, node model.Node) []string {
		return []string{cluster, node.Host, node.Name}
	}
	labelValuesThreadPool = func(cluster string, node model.Node, pool string) []string {
		return append(labelValuesNode(cluster, node), pool)
	}
	labelValuesFilesystem = func(cluster string, node model.Node, mount string, path string) []string {
		return append(labelValuesNode(cluster, node), mount, path)
	}
)

type nodeMetric struct {
	*metrics.Metric
	Value func(node model.Node) float64
}

type gcCollectionMetric struct {
	*metrics.Metric
	Value func(gcStats model.NodeJVMGCCollector) float64
}

type breakerMetric struct {
	*metrics.Metric
	Value  func(breakerStats model.Breaker) float64
	Labels func(cluster string, node model.Node, breaker string) []string
}

type threadPoolMetric struct {
	*metrics.Metric
	Value func(threadPoolStats model.ThreadPool) float64
}

type filesystemMetric struct {
	*metrics.Metric
	Value func(fsStats model.NodeFSData) float64
}

// Collector is an node metrics collector
type Collector struct {
	esClient                 elasticsearch.Client
	exportMetricsForAllNodes bool

	nodeMetrics         []*nodeMetric
	gcCollectionMetrics []*gcCollectionMetric
	breakerMetrics      []*breakerMetric
	threadPoolMetrics   []*threadPoolMetric
	filesystemMetrics   []*filesystemMetric
}

func newNodeIndexMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.Node) float64) *nodeMetric {
	return &nodeMetric{
		Metric: metrics.New(t, "node_indices", name, help, labelsNode),
		Value:  valueExtractor,
	}
}

func newFSGauge(name, help string, valueExtractor func(model.NodeFSData) float64) *filesystemMetric {
	return &filesystemMetric{
		Metric: metrics.New(prometheus.GaugeValue, "filesystem_data", name, help, labelsFilesystem),
		Value:  valueExtractor,
	}
}

func newThreadPoolMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.ThreadPool) float64) *threadPoolMetric {
	return &threadPoolMetric{
		Metric: metrics.New(t, "thread_pool", name, help, labelsThreadPool),
		Value:  valueExtractor,
	}
}

func newJVMGCMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.NodeJVMGCCollector) float64) *gcCollectionMetric {
	return &gcCollectionMetric{
		Metric: metrics.New(t, "jvm_gc", name, help, labelsJVMGC),
		Value:  valueExtractor,
	}
}

func newJVMMemoryMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.Node) float64) *nodeMetric {
	return &nodeMetric{
		Metric: metrics.New(t, "jvm_memory", name, help, labelsNode),
		Value:  valueExtractor,
	}
}

func newProcessMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.Node) float64) *nodeMetric {
	return &nodeMetric{
		Metric: metrics.New(t, "process", name, help, labelsNode),
		Value:  valueExtractor,
	}
}

func newTransportMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.Node) float64) *nodeMetric {
	return &nodeMetric{
		Metric: metrics.New(t, "transport", name, help, labelsNode),
		Value:  valueExtractor,
	}
}

func newBreakerMetric(t prometheus.ValueType, name, help string, valueExtractor func(model.Breaker) float64) *breakerMetric {
	return &breakerMetric{
		Metric: metrics.New(t, "breakers", name, help, labelsBreaker),
		Value:  valueExtractor,
	}
}

// NewCollector returns new nodes metrics collector
func NewCollector(esClient elasticsearch.Client, exportMetricsForAllNodes bool) *Collector {
	return &Collector{
		esClient:                 esClient,
		exportMetricsForAllNodes: exportMetricsForAllNodes,

		nodeMetrics: []*nodeMetric{
			newNodeIndexMetric(
				prometheus.GaugeValue, "fielddata_memory_size_bytes", "Field data cache memory usage in bytes",
				func(n model.Node) float64 { return float64(n.Indices.FieldData.MemorySize) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "fielddata_evictions", "Evictions from field data",
				func(n model.Node) float64 { return float64(n.Indices.FieldData.Evictions) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "filter_cache_memory_size_bytes", "Filter cache memory usage in bytes",
				func(n model.Node) float64 { return float64(n.Indices.FilterCache.MemorySize) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "filter_cache_evictions", "Evictions from filter cache",
				func(n model.Node) float64 { return float64(n.Indices.FilterCache.Evictions) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "query_cache_memory_size_bytes", "Query cache memory usage in bytes",
				func(n model.Node) float64 { return float64(n.Indices.QueryCache.MemorySize) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "query_cache_evictions", "Evictions from query cache",
				func(n model.Node) float64 { return float64(n.Indices.QueryCache.Evictions) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "request_cache_memory_size_bytes", "Request cache memory usage in bytes",
				func(n model.Node) float64 { return float64(n.Indices.RequestCache.MemorySize) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "request_cache_evictions", "Evictions from request cache",
				func(n model.Node) float64 { return float64(n.Indices.RequestCache.Evictions) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "translog_operations", "Total translog operations",
				func(n model.Node) float64 { return float64(n.Indices.Translog.Operations) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "translog_size_in_bytes", "Total translog size in bytes",
				func(n model.Node) float64 { return float64(n.Indices.Translog.Size) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "get_time_seconds", "Total get time in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Get.Time / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "get_total", "Total get",
				func(n model.Node) float64 { return float64(n.Indices.Get.Total) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "get_missing_time_seconds", "Total time of get missing in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Get.MissingTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "get_missing_total", "Total get missing",
				func(n model.Node) float64 { return float64(n.Indices.Get.MissingTotal) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "get_exists_time_seconds", "Total time get exists in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Get.ExistsTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "get_exists_total", "Total get exists operations",
				func(n model.Node) float64 { return float64(n.Indices.Get.ExistsTotal) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "refresh_time_seconds_total", "Total refreshes",
				func(n model.Node) float64 { return float64(n.Indices.Refresh.TotalTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "refresh_total", "Total time spent refreshing in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Refresh.Total) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "search_query_time_seconds", "Total search query time in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Search.QueryTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "search_query_total", "Total number of queries",
				func(n model.Node) float64 { return float64(n.Indices.Search.QueryTotal) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "search_fetch_time_seconds", "Total search fetch time in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Search.FetchTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "search_fetch_total", "Total number of fetches",
				func(n model.Node) float64 { return float64(n.Indices.Search.FetchTotal) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "docs", "Count of documents on this node",
				func(n model.Node) float64 { return float64(n.Indices.Docs.Count) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "docs_deleted", "Count of deleted documents on this node",
				func(n model.Node) float64 { return float64(n.Indices.Docs.Deleted) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "store_size_bytes", "Current size of stored index data in bytes",
				func(n model.Node) float64 { return float64(n.Indices.Store.Size) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "store_throttle_time_seconds_total", "Throttle time for index store in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Store.ThrottleTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "segments_memory_bytes", "Current memory size of segments in bytes",
				func(n model.Node) float64 { return float64(n.Indices.Segments.Memory) },
			),
			newNodeIndexMetric(
				prometheus.GaugeValue, "segments_count", "Count of index segments on this node",
				func(n model.Node) float64 { return float64(n.Indices.Segments.Count) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "flush_total", "Total flushes",
				func(n model.Node) float64 { return float64(n.Indices.Flush.Total) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "flush_time_seconds", "Cumulative flush time in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Flush.Time / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "indexing_index_time_seconds_total", "Cumulative index time in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Indexing.IndexTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "indexing_index_total", "Total index calls",
				func(n model.Node) float64 { return float64(n.Indices.Indexing.IndexTotal) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "indexing_delete_time_seconds_total", "Total time indexing delete in seconds",
				func(n model.Node) float64 { return float64(n.Indices.Indexing.DeleteTime / 1000) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "indexing_delete_total", "Total indexing deletes",
				func(n model.Node) float64 { return float64(n.Indices.Indexing.DeleteTotal) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "merges_total", "Total merges",
				func(node model.Node) float64 { return float64(node.Indices.Merges.Total) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "merges_docs_total", "Cumulative docs merged",
				func(node model.Node) float64 { return float64(node.Indices.Merges.TotalDocs) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "merges_total_size_bytes_total", "Total merge size in bytes",
				func(node model.Node) float64 { return float64(node.Indices.Merges.TotalSize) },
			),
			newNodeIndexMetric(
				prometheus.CounterValue, "merges_total_time_seconds_total", "Total time spent merging in seconds",
				func(node model.Node) float64 { return float64(node.Indices.Merges.TotalTime / 1000) },
			),

			newJVMMemoryMetric(
				prometheus.GaugeValue, "heap_used_bytes", "JVM memory currently used by heap",
				func(n model.Node) float64 { return float64(n.JVM.Mem.HeapUsed) },
			),
			newJVMMemoryMetric(
				prometheus.GaugeValue, "non_heap_used_bytes", "JVM memory currently used by area",
				func(node model.Node) float64 { return float64(node.JVM.Mem.NonHeapUsed) },
			),
			newJVMMemoryMetric(
				prometheus.GaugeValue, "heap_max_bytes", "JVM memory max",
				func(node model.Node) float64 { return float64(node.JVM.Mem.HeapMax) },
			),
			newJVMMemoryMetric(
				prometheus.GaugeValue, "heap_committed_bytes", "JVM memory currently committed by area",
				func(node model.Node) float64 { return float64(node.JVM.Mem.HeapCommitted) },
			),
			newJVMMemoryMetric(
				prometheus.GaugeValue, "non_heap_committed_bytes", "JVM memory currently committed by area",
				func(node model.Node) float64 { return float64(node.JVM.Mem.NonHeapCommitted) },
			),
			newProcessMetric(
				prometheus.GaugeValue, "cpu_percent", "Percent CPU used by process",
				func(node model.Node) float64 { return float64(node.Process.CPU.Percent) },
			),
			newProcessMetric(
				prometheus.GaugeValue, "mem_resident_size_bytes", "Resident memory in use by process in bytes",
				func(node model.Node) float64 { return float64(node.Process.Memory.Resident) },
			),
			newProcessMetric(
				prometheus.GaugeValue, "mem_share_size_bytes", "Shared memory in use by process in bytes",
				func(node model.Node) float64 { return float64(node.Process.Memory.Share) },
			),
			newProcessMetric(
				prometheus.GaugeValue, "mem_virtual_size_bytes", "Total virtual memory used in bytes",
				func(node model.Node) float64 { return float64(node.Process.Memory.TotalVirtual) },
			),
			newProcessMetric(
				prometheus.GaugeValue, "open_files_count", "Open file descriptors",
				func(node model.Node) float64 { return float64(node.Process.OpenFD) },
			),
			newProcessMetric(
				prometheus.CounterValue, "cpu_time_total_seconds_sum", "Total process CPU time in seconds",
				func(node model.Node) float64 { return float64(node.Process.CPU.Total / 1000) },
			),
			newProcessMetric(
				prometheus.CounterValue, "cpu_time_system_seconds_sum", "Process system CPU time in seconds",
				func(node model.Node) float64 { return float64(node.Process.CPU.Sys / 1000) },
			),
			newProcessMetric(
				prometheus.CounterValue, "cpu_time_user_seconds_sum", "Process CPU time in seconds",
				func(node model.Node) float64 { return float64(node.Process.CPU.User / 1000) },
			),
			newTransportMetric(
				prometheus.CounterValue, "rx_packets_total", "Count of packets received",
				func(node model.Node) float64 { return float64(node.Transport.RxCount) },
			),
			newTransportMetric(
				prometheus.CounterValue, "rx_size_bytes_total", "Total number of bytes received",
				func(node model.Node) float64 { return float64(node.Transport.RxSize) },
			),
			newTransportMetric(
				prometheus.CounterValue, "tx_packets_total", "Count of packets sent",
				func(node model.Node) float64 { return float64(node.Transport.TxCount) },
			),
			newTransportMetric(
				prometheus.CounterValue, "tx_size_bytes_total", "Total number of bytes sent",
				func(node model.Node) float64 { return float64(node.Transport.TxSize) },
			),
		},
		gcCollectionMetrics: []*gcCollectionMetric{
			newJVMGCMetric(
				prometheus.CounterValue, "collection_seconds_count", "Count of JVM GC runs",
				func(gc model.NodeJVMGCCollector) float64 { return float64(gc.CollectionCount) },
			),
			newJVMGCMetric(
				prometheus.CounterValue, "collection_seconds_sum", "GC run time in seconds",
				func(gc model.NodeJVMGCCollector) float64 { return float64(gc.CollectionTime / 1000) },
			),
		},
		breakerMetrics: []*breakerMetric{
			newBreakerMetric(
				prometheus.GaugeValue, "estimated_size_bytes", "Estimated size in bytes of breaker",
				func(breakerStats model.Breaker) float64 { return float64(breakerStats.EstimatedSize) },
			),
			newBreakerMetric(
				prometheus.GaugeValue, "limit_size_bytes", "Limit size in bytes for breaker",
				func(breakerStats model.Breaker) float64 { return float64(breakerStats.LimitSize) },
			),
			newBreakerMetric(
				prometheus.GaugeValue, "tripped", "tripped for breaker",
				func(breakerStats model.Breaker) float64 { return float64(breakerStats.Tripped) },
			),
		},
		threadPoolMetrics: []*threadPoolMetric{
			newThreadPoolMetric(
				prometheus.CounterValue, "completed_count", "Thread Pool operations completed",
				func(threadPoolStats model.ThreadPool) float64 { return float64(threadPoolStats.Completed) },
			),
			newThreadPoolMetric(
				prometheus.CounterValue, "rejected_count", "Thread Pool operations rejected",
				func(threadPoolStats model.ThreadPool) float64 { return float64(threadPoolStats.Rejected) },
			),
			newThreadPoolMetric(
				prometheus.GaugeValue, "active_count", "Thread Pool threads active",
				func(threadPoolStats model.ThreadPool) float64 { return float64(threadPoolStats.Active) },
			),
			newThreadPoolMetric(
				prometheus.GaugeValue, "largest_count", "Thread Pool largest threads count",
				func(threadPoolStats model.ThreadPool) float64 { return float64(threadPoolStats.Largest) },
			),
			newThreadPoolMetric(
				prometheus.GaugeValue, "queue_count", "Thread Pool operations queued",
				func(threadPoolStats model.ThreadPool) float64 { return float64(threadPoolStats.Queue) },
			),
			newThreadPoolMetric(
				prometheus.GaugeValue, "threads_count", "Thread Pool current threads count",
				func(threadPoolStats model.ThreadPool) float64 { return float64(threadPoolStats.Threads) },
			),
		},
		filesystemMetrics: []*filesystemMetric{
			newFSGauge(
				"available_bytes", "Available space on block device in bytes",
				func(fs model.NodeFSData) float64 { return float64(fs.Available) },
			),
			newFSGauge(
				"free_bytes", "Free space on block device in bytes",
				func(fs model.NodeFSData) float64 { return float64(fs.Free) },
			),
			newFSGauge(
				"size_bytes", "Size of block device in bytes",
				func(fs model.NodeFSData) float64 { return float64(fs.Total) },
			),
		},
	}
}

// Describe implements prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.nodeMetrics {
		ch <- metric.Desc()
	}
	for _, metric := range c.breakerMetrics {
		ch <- metric.Desc()
	}
	for _, metric := range c.gcCollectionMetrics {
		ch <- metric.Desc()
	}
	for _, metric := range c.threadPoolMetrics {
		ch <- metric.Desc()
	}
	for _, metric := range c.filesystemMetrics {
		ch <- metric.Desc()
	}
}

// Collect writes data to metrics channel
func (c *Collector) Collect(clusterName string, ch chan<- prometheus.Metric) {
	nodeStats, err := c.esClient.Nodes(c.exportMetricsForAllNodes)
	if err != nil {
		log.Println("ERROR: failed to fetch nodes stats: ", err)
		return
	}

	for _, node := range nodeStats.Nodes {
		for _, metric := range c.nodeMetrics {
			ch <- prometheus.MustNewConstMetric(
				metric.Desc(),
				metric.Type(),
				metric.Value(node),
				labelValuesNode(clusterName, node)...,
			)
		}

		// GC Stats
		for collector, gcStats := range node.JVM.GC.Collectors {
			for _, metric := range c.gcCollectionMetrics {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc(),
					metric.Type(),
					metric.Value(gcStats),
					append(labelValuesNode(clusterName, node), collector)...,
				)
			}
		}

		// Breaker stats
		for breaker, bstats := range node.Breakers {
			for _, metric := range c.breakerMetrics {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc(),
					metric.Type(),
					metric.Value(bstats),
					append(labelValuesNode(clusterName, node), breaker)...,
				)
			}
		}

		// Thread Pool stats
		for pool, pstats := range node.ThreadPool {
			for _, metric := range c.threadPoolMetrics {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc(),
					metric.Type(),
					metric.Value(pstats),
					labelValuesThreadPool(clusterName, node, pool)...,
				)
			}
		}

		// File System Stats
		for _, fsStats := range node.FS.Data {
			for _, metric := range c.filesystemMetrics {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc(),
					metric.Type(),
					metric.Value(fsStats),
					labelValuesFilesystem(clusterName, node, fsStats.Mount, fsStats.Path)...,
				)
			}
		}
	}
}
