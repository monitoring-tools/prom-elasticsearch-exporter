package testdata

import (
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
)

// Test data for cluster health info
var (
	ClusterHealthIndicesBody = `
{
	"cluster_name": "my-huge-cluster",
	"status": "green",
	"timed_out": false,
	"number_of_nodes": 89,
	"number_of_data_nodes": 86,
	"active_primary_shards": 45,
	"active_shards": 174,
	"relocating_shards": 0,
	"initializing_shards": 0,
	"unassigned_shards": 0,
	"delayed_unassigned_shards": 0,
	"number_of_pending_tasks": 0,
	"number_of_in_flight_fetch": 0,
	"task_max_waiting_in_queue_millis": 0,
	"active_shards_percent_as_number": 100.0,
	"indices": {
		"some-index": {
			"status": "green",
			"number_of_shards": 42,
			"number_of_replicas": 3,
			"active_primary_shards": 42,
			"active_shards": 168,
			"relocating_shards": 0,
			"initializing_shards": 0,
			"unassigned_shards": 0
		}
	}
}
`

	ClusterHealthIndices = &model.ClusterHealth{
		ClusterName:                 "my-huge-cluster",
		Status:                      "green",
		NumberOfNodes:               89,
		NumberOfDataNodes:           86,
		ActivePrimaryShards:         45,
		ActiveShards:                174,
		ActiveShardsPercentAsNumber: 100.0,
		Indices: map[string]model.ClusterHealthIndex{
			"some-index": {
				Status:              "green",
				NumberOfShards:      42,
				NumberOfReplicas:    3,
				ActivePrimaryShards: 42,
				ActiveShards:        168,
			},
		},
	}
)
