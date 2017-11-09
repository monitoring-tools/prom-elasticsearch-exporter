package model

import (
	"encoding/json"
)

// Nodes is a representation of ElasticSearch cluster nodes statistics
type Nodes struct {
	ClusterName string `json:"cluster_name"`
	Nodes       map[string]Node
}

// Node is a representation of ElasticSearch node stats
type Node struct {
	Name             string                `json:"name"`
	Host             string                `json:"host"`
	Timestamp        int64                 `json:"timestamp"`
	TransportAddress string                `json:"transport_address"`
	Hostname         string                `json:"hostname"`
	Indices          NodeIndices           `json:"indices"`
	OS               OS                    `json:"os"`
	Network          Network               `json:"network"`
	FS               FS                    `json:"fs"`
	ThreadPool       map[string]ThreadPool `json:"thread_pool"`
	JVM              NodeJVM               `json:"jvm"`
	Breakers         map[string]Breaker    `json:"breakers"`
	Transport        Transport             `json:"transport"`
	Process          Process               `json:"process"`
}

// Breaker is a representation of statistics about the field data circuit breaker
type Breaker struct {
	EstimatedSize int64   `json:"estimated_size_in_bytes"`
	LimitSize     int64   `json:"limit_size_in_bytes"`
	Overhead      float64 `json:"overhead"`
	Tripped       int64   `json:"tripped"`
}

// NodeJVM is a representation of JVM stats, memory pool information, garbage collection, buffer pools, number of loaded/unloaded classes
type NodeJVM struct {
	BufferPools map[string]NodeJVMBufferPool `json:"buffer_pools"`
	GC          NodeJVMGC                    `json:"gc"`
	Mem         NodeJVMMem                   `json:"mem"`
}

// NodeJVMGC is a representation of JVM GC stats for all kinds of collectors
type NodeJVMGC struct {
	Collectors map[string]NodeJVMGCCollector `json:"collectors"`
}

// NodeJVMGCCollector is a representation of JVM GC collector stats
type NodeJVMGCCollector struct {
	CollectionCount int64 `json:"collection_count"`
	CollectionTime  int64 `json:"collection_time_in_millis"`
}

// NodeJVMBufferPool is a representation of JVM buffer pool stats
type NodeJVMBufferPool struct {
	Count         int64 `json:"count"`
	TotalCapacity int64 `json:"total_capacity_in_bytes"`
	Used          int64 `json:"used_in_bytes"`
}

// NodeJVMMem is a representation of JVM memory usage stats
type NodeJVMMem struct {
	HeapCommitted    int64 `json:"heap_committed_in_bytes"`
	HeapUsed         int64 `json:"heap_used_in_bytes"`
	HeapMax          int64 `json:"heap_max_in_bytes"`
	NonHeapCommitted int64 `json:"non_heap_committed_in_bytes"`
	NonHeapUsed      int64 `json:"non_heap_used_in_bytes"`
}

// Network is a representation of network usage stats
type Network struct {
	TCP NodeTCP `json:"tcp"`
}

// Transport is a representation of transport statistics about sent and received bytes in cluster communication
type Transport struct {
	ServerOpen int64 `json:"server_open"`
	RxCount    int64 `json:"rx_count"`
	RxSize     int64 `json:"rx_size_in_bytes"`
	TxCount    int64 `json:"tx_count"`
	TxSize     int64 `json:"tx_size_in_bytes"`
}

// ThreadPool is a representation of statistics about each thread pool, including current size, queue and rejected tasks
type ThreadPool struct {
	Threads   int64 `json:"threads"`
	Queue     int64 `json:"queue"`
	Active    int64 `json:"active"`
	Rejected  int64 `json:"rejected"`
	Largest   int64 `json:"largest"`
	Completed int64 `json:"completed"`
}

// NodeTCP is a representation of tcp stats on node
type NodeTCP struct {
	ActiveOpens  int64 `json:"active_opens"`
	PassiveOpens int64 `json:"passive_opens"`
	CurrEstab    int64 `json:"curr_estab"`
	InSegs       int64 `json:"in_segs"`
	OutSegs      int64 `json:"out_segs"`
	RetransSegs  int64 `json:"retrans_segs"`
	EstabResets  int64 `json:"estab_resets"`
	AttemptFails int64 `json:"attempt_fails"`
	InErrs       int64 `json:"in_errs"`
	OutRsts      int64 `json:"out_rsts"`
}

// NodeIndices is a representation of indices stats (size, document count, indexing and deletion times, search times, field cache size, merges and flushes)
type NodeIndices struct {
	Docs         NodeIndicesDocs
	Store        NodeIndicesStore
	Indexing     NodeIndicesIndexing
	Merges       NodeIndicesMerges
	Get          NodeIndicesGet
	Search       NodeIndicesSearch
	FieldData    NodeIndicesCache `json:"fielddata"`
	FilterCache  NodeIndicesCache `json:"filter_cache"`
	QueryCache   NodeIndicesCache `json:"query_cache"`
	RequestCache NodeIndicesCache `json:"request_cache"`
	Flush        NodeIndicesFlush
	Segments     NodeIndicesSegments
	Refresh      NodeIndicesRefresh
	Translog     NodeIndicesTranslog
}

// NodeIndicesDocs is a representation of indices docs stats
type NodeIndicesDocs struct {
	Count   int64 `json:"count"`
	Deleted int64 `json:"deleted"`
}

// NodeIndicesRefresh is a representation of index refresh stats
type NodeIndicesRefresh struct {
	Total     int64 `json:"total"`
	TotalTime int64 `json:"total_time_in_millis"`
}

// NodeIndicesTranslog is a representation of index translog stats
type NodeIndicesTranslog struct {
	Operations int64 `json:"operations"`
	Size       int64 `json:"size_in_bytes"`
}

// NodeIndicesSegments is a representation of index segments stats
type NodeIndicesSegments struct {
	Count  int64 `json:"count"`
	Memory int64 `json:"memory_in_bytes"`
}

// NodeIndicesStore is a representation of index store stats
type NodeIndicesStore struct {
	Size         int64 `json:"size_in_bytes"`
	ThrottleTime int64 `json:"throttle_time_in_millis"`
}

// NodeIndicesIndexing is a representation of index indexing stats
type NodeIndicesIndexing struct {
	IndexTotal    int64 `json:"index_total"`
	IndexTime     int64 `json:"index_time_in_millis"`
	IndexCurrent  int64 `json:"index_current"`
	DeleteTotal   int64 `json:"delete_total"`
	DeleteTime    int64 `json:"delete_time_in_millis"`
	DeleteCurrent int64 `json:"delete_current"`
}

// NodeIndicesMerges is a representation of index merging stats
type NodeIndicesMerges struct {
	Current     int64 `json:"current"`
	CurrentDocs int64 `json:"current_docs"`
	CurrentSize int64 `json:"current_size_in_bytes"`
	Total       int64 `json:"total"`
	TotalDocs   int64 `json:"total_docs"`
	TotalSize   int64 `json:"total_size_in_bytes"`
	TotalTime   int64 `json:"total_time_in_millis"`
}

// NodeIndicesGet is a representation of index get stats
type NodeIndicesGet struct {
	Total        int64 `json:"total"`
	Time         int64 `json:"time_in_millis"`
	ExistsTotal  int64 `json:"exists_total"`
	ExistsTime   int64 `json:"exists_time_in_millis"`
	MissingTotal int64 `json:"missing_total"`
	MissingTime  int64 `json:"missing_time_in_millis"`
	Current      int64 `json:"current"`
}

// NodeIndicesSearch is a representation of index search stats
type NodeIndicesSearch struct {
	OpenContext  int64 `json:"open_contexts"`
	QueryTotal   int64 `json:"query_total"`
	QueryTime    int64 `json:"query_time_in_millis"`
	QueryCurrent int64 `json:"query_current"`
	FetchTotal   int64 `json:"fetch_total"`
	FetchTime    int64 `json:"fetch_time_in_millis"`
	FetchCurrent int64 `json:"fetch_current"`
}

// NodeIndicesFlush is a representation of index flush stats
type NodeIndicesFlush struct {
	Total int64 `json:"total"`
	Time  int64 `json:"total_time_in_millis"`
}

// NodeIndicesCache is a representation of index cache stats
type NodeIndicesCache struct {
	Evictions  int64 `json:"evictions"`
	MemorySize int64 `json:"memory_size_in_bytes"`
	CacheCount int64 `json:"cache_count"`
	CacheSize  int64 `json:"cache_size"`
	HitCount   int64 `json:"hit_count"`
	MissCount  int64 `json:"miss_count"`
	TotalCount int64 `json:"total_count"`
}

// OS is a representation of  operating system stats, load average, mem, swap
type OS struct {
	Timestamp int64 `json:"timestamp"`
	Uptime    int64 `json:"uptime_in_millis"`
	// LoadAvg was an array of per-cpu values pre-2.0, and is a string in 2.0
	// Leaving this here in case we want to implement parsing logic later
	LoadAvg json.RawMessage `json:"load_average"`
	CPU     NodeOSCPU       `json:"cpu"`
	Mem     NodeOSMem       `json:"mem"`
	Swap    NodeOSSwap      `json:"swap"`
}

// NodeOSMem is a representation of OS memory stats
type NodeOSMem struct {
	Free       int64 `json:"free_in_bytes"`
	Used       int64 `json:"used_in_bytes"`
	ActualFree int64 `json:"actual_free_in_bytes"`
	ActualUsed int64 `json:"actual_used_in_bytes"`
}

// NodeOSSwap is a representation of swap stats
type NodeOSSwap struct {
	Used int64 `json:"used_in_bytes"`
	Free int64 `json:"free_in_bytes"`
}

// NodeOSCPU is a representation of CPU usage stats
type NodeOSCPU struct {
	Sys   int64 `json:"sys"`
	User  int64 `json:"user"`
	Idle  int64 `json:"idle"`
	Steal int64 `json:"stolen"`
}

// Process is a representation of process statistics, memory consumption, cpu usage, open file descriptors
type Process struct {
	Timestamp int64          `json:"timestamp"`
	OpenFD    int64          `json:"open_file_descriptors"`
	MaxFD     int64          `json:"max_file_descriptors"`
	CPU       NodeProcessCPU `json:"cpu"`
	Memory    NodeProcessMem `json:"mem"`
}

// NodeProcessMem is a representation of memory usage stats
type NodeProcessMem struct {
	Resident     int64 `json:"resident_in_bytes"`
	Share        int64 `json:"share_in_bytes"`
	TotalVirtual int64 `json:"total_virtual_in_bytes"`
}

// NodeProcessCPU is a representation of process CPU usage stats
type NodeProcessCPU struct {
	Percent int64 `json:"percent"`
	Sys     int64 `json:"sys_in_millis"`
	User    int64 `json:"user_in_millis"`
	Total   int64 `json:"total_in_millis"`
}

// NodeHTTP is a representation of HTTP connections stats
type NodeHTTP struct {
	CurrentOpen int64 `json:"current_open"`
	TotalOpen   int64 `json:"total_open"`
}

// FS is a representation of file system information, data path, free disk space, read/write stats
type FS struct {
	Timestamp int64        `json:"timestamp"`
	Data      []NodeFSData `json:"data"`
}

// NodeFSData is a representation of filesystem stats
type NodeFSData struct {
	Path          string `json:"path"`
	Mount         string `json:"mount"`
	Device        string `json:"dev"`
	Total         int64  `json:"total_in_bytes"`
	Free          int64  `json:"free_in_bytes"`
	Available     int64  `json:"available_in_bytes"`
	DiskReads     int64  `json:"disk_reads"`
	DiskWrites    int64  `json:"disk_writes"`
	DiskReadSize  int64  `json:"disk_read_size_in_bytes"`
	DiskWriteSize int64  `json:"disk_write_size_in_bytes"`
}
