package testdata

// Test data for nodes info
var (
	NodesBody = `
	{
	"_nodes": {
		"total": 72,
		"successful": 72,
		"failed": 0
	},
	"cluster_name": "my-huge-cluster",
	"nodes": {
		"3a6VFkY8SLOI4J6ljALdhQ": {
			"timestamp": 1506955347994,
			"name": "mynode0",
			"transport_address": "10.81.112.34:9301",
			"host": "mynode0",
			"ip": "10.81.112.34:9301",
			"roles": ["data", "ingest"],
			"attributes": {
				"rack_id": "rack1"
			},
			"indices": {
				"docs": {
					"count": 13279235,
					"deleted": 9524289
				},
				"store": {
					"size_in_bytes": 84529297743,
					"throttle_time_in_millis": 0
				},
				"indexing": {
					"index_total": 29257252,
					"index_time_in_millis": 172461064,
					"index_current": 0,
					"index_failed": 0,
					"delete_total": 14036033,
					"delete_time_in_millis": 1503493,
					"delete_current": 0,
					"noop_update_total": 17,
					"is_throttled": false,
					"throttle_time_in_millis": 0
				},
				"get": {
					"total": 391122,
					"time_in_millis": 113114,
					"exists_total": 317081,
					"exists_time_in_millis": 102210,
					"missing_total": 74041,
					"missing_time_in_millis": 10904,
					"current": 0
				},
				"search": {
					"open_contexts": 9,
					"query_total": 107511688,
					"query_time_in_millis": 3987735236,
					"query_current": 7,
					"fetch_total": 28802170,
					"fetch_time_in_millis": 7608623,
					"fetch_current": 0,
					"scroll_total": 0,
					"scroll_time_in_millis": 0,
					"scroll_current": 0,
					"suggest_total": 1274654,
					"suggest_time_in_millis": 67729733,
					"suggest_current": 0
				},
				"merges": {
					"current": 0,
					"current_docs": 0,
					"current_size_in_bytes": 0,
					"total": 25704,
					"total_time_in_millis": 192260617,
					"total_docs": 427561358,
					"total_size_in_bytes": 870547589287,
					"total_stopped_time_in_millis": 0,
					"total_throttled_time_in_millis": 33387138,
					"total_auto_throttle_in_bytes": 26214400
				},
				"refresh": {
					"total": 153183,
					"total_time_in_millis": 24574251,
					"listeners": 0
				},
				"flush": {
					"total": 346,
					"total_time_in_millis": 265096
				},
				"warmer": {
					"current": 0,
					"total": 153509,
					"total_time_in_millis": 275292
				},
				"query_cache": {
					"memory_size_in_bytes": 789142361,
					"total_count": 13777344024,
					"hit_count": 1627829341,
					"miss_count": 12149514683,
					"cache_size": 59654,
					"cache_count": 2550865,
					"evictions": 2491211
				},
				"fielddata": {
					"memory_size_in_bytes": 1000712,
					"evictions": 0
				},
				"completion": {
					"size_in_bytes": 3427345025
				},
				"segments": {
					"count": 163,
					"memory_in_bytes": 3566884507,
					"terms_memory_in_bytes": 3540996182,
					"stored_fields_memory_in_bytes": 7910312,
					"term_vectors_memory_in_bytes": 0,
					"norms_memory_in_bytes": 4083904,
					"points_memory_in_bytes": 1280333,
					"doc_values_memory_in_bytes": 12613776,
					"index_writer_memory_in_bytes": 29039966,
					"version_map_memory_in_bytes": 17928,
					"fixed_bit_set_memory_in_bytes": 2857584,
					"max_unsafe_auto_id_timestamp": -1,
					"file_sizes": {}
				},
				"translog": {
					"operations": 189425,
					"size_in_bytes": 795815372
				},
				"request_cache": {
					"memory_size_in_bytes": 7348747,
					"evictions": 0,
					"hit_count": 324629,
					"miss_count": 15589277
				},
				"recovery": {
					"current_as_source": 0,
					"current_as_target": 0,
					"throttle_time_in_millis": 2419159
				}
			},
			"os": {
				"timestamp": 1506955347338,
				"cpu": {
					"percent": 52,
					"load_average": {
						"1m": 19.12,
						"5m": 19.06,
						"15m": 19.37
					}
				},
				"mem": {
					"total_in_bytes": 135083540480,
					"free_in_bytes": 12070580224,
					"used_in_bytes": 123012960256,
					"free_percent": 9,
					"used_percent": 91
				},
				"swap": {
					"total_in_bytes": 0,
					"free_in_bytes": 0,
					"used_in_bytes": 0
				}
			},
			"process": {
				"timestamp": 1506955347338,
				"open_file_descriptors": 2569,
				"max_file_descriptors": 1048576,
				"cpu": {
					"percent": 28,
					"total_in_millis": 6708003540
				},
				"mem": {
					"total_virtual_in_bytes": 140780392448
				}
			},
			"jvm": {
				"timestamp": 1506955347340,
				"uptime_in_millis": 940628583,
				"mem": {
					"heap_used_in_bytes": 11701828904,
					"heap_used_percent": 36,
					"heap_committed_in_bytes": 31968002048,
					"heap_max_in_bytes": 31968002048,
					"non_heap_used_in_bytes": 176416344,
					"non_heap_committed_in_bytes": 185540608,
					"pools": {
						"young": {
							"used_in_bytes": 1473179416,
							"max_in_bytes": 1954217984,
							"peak_used_in_bytes": 1954217984,
							"peak_max_in_bytes": 1954217984
						},
						"survivor": {
							"used_in_bytes": 81838144,
							"max_in_bytes": 244252672,
							"peak_used_in_bytes": 244252672,
							"peak_max_in_bytes": 244252672
						},
						"old": {
							"used_in_bytes": 10146811344,
							"max_in_bytes": 29769531392,
							"peak_used_in_bytes": 22717405456,
							"peak_max_in_bytes": 29769531392
						}
					}
				},
				"threads": {
					"count": 323,
					"peak_count": 344
				},
				"gc": {
					"collectors": {
						"young": {
							"collection_count": 599316,
							"collection_time_in_millis": 25847320
						},
						"old": {
							"collection_count": 376,
							"collection_time_in_millis": 152045
						}
					}
				},
				"buffer_pools": {
					"direct": {
						"count": 293,
						"used_in_bytes": 1086758728,
						"total_capacity_in_bytes": 1086758727
					},
					"mapped": {
						"count": 1427,
						"used_in_bytes": 83320823327,
						"total_capacity_in_bytes": 83320823327
					}
				},
				"classes": {
					"current_loaded_count": 14110,
					"total_loaded_count": 14549,
					"total_unloaded_count": 439
				}
			},
			"thread_pool": {
				"bulk": {
					"threads": 2,
					"queue": 0,
					"active": 0,
					"rejected": 19773,
					"largest": 2,
					"completed": 21022726
				},
				"fetch_shard_started": {
					"threads": 1,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 20,
					"completed": 61
				},
				"fetch_shard_store": {
					"threads": 1,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 30,
					"completed": 7266
				},
				"flush": {
					"threads": 1,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 3,
					"completed": 361
				},
				"force_merge": {
					"threads": 0,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 0,
					"completed": 0
				},
				"generic": {
					"threads": 21,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 21,
					"completed": 1328615
				},
				"get": {
					"threads": 0,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 0,
					"completed": 0
				},
				"index": {
					"threads": 0,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 0,
					"completed": 0
				},
				"listener": {
					"threads": 0,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 0,
					"completed": 0
				},
				"management": {
					"threads": 5,
					"queue": 0,
					"active": 4,
					"rejected": 0,
					"largest": 5,
					"completed": 20443377
				},
				"refresh": {
					"threads": 2,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 2,
					"completed": 61886
				},
				"search": {
					"threads": 40,
					"queue": 0,
					"active": 11,
					"rejected": 805683,
					"largest": 40,
					"completed": 140330444
				},
				"snapshot": {
					"threads": 0,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 0,
					"completed": 0
				},
				"warmer": {
					"threads": 5,
					"queue": 0,
					"active": 0,
					"rejected": 0,
					"largest": 5,
					"completed": 4661871
				}
			},
			"fs": {
				"timestamp": 1506955347340,
				"total": {
					"total_in_bytes": 944856870912,
					"free_in_bytes": 615839916032,
					"available_in_bytes": 567820169216,
					"spins": "true"
				},
				"data": [{
					"path": "/var/lib/elasticsearch/nodes/0",
					"mount": "/var/lib/elasticsearch (/dev/mapper/VolGroup01-var--lib--docker)",
					"type": "ext4",
					"total_in_bytes": 944856870912,
					"free_in_bytes": 615839916032,
					"available_in_bytes": 567820169216,
					"spins": "true"
				}],
				"io_stats": {
					"devices": [{
						"device_name": "dm-2",
						"operations": 529792318,
						"read_operations": 38345976,
						"write_operations": 491446342,
						"read_kilobytes": 836695368,
						"write_kilobytes": 4587166320
					}],
					"total": {
						"operations": 529792318,
						"read_operations": 38345976,
						"write_operations": 491446342,
						"read_kilobytes": 836695368,
						"write_kilobytes": 4587166320
					}
				}
			},
			"transport": {
				"server_open": 849,
				"rx_count": 366595608,
				"rx_size_in_bytes": 1239745834235,
				"tx_count": 366595764,
				"tx_size_in_bytes": 1773813734027
			},
			"http": {
				"current_open": 18,
				"total_opened": 26152
			},
			"breakers": {
				"request": {
					"limit_size_in_bytes": 19180801228,
					"limit_size": "17.8gb",
					"estimated_size_in_bytes": 210944,
					"estimated_size": "206kb",
					"overhead": 1.0,
					"tripped": 0
				},
				"fielddata": {
					"limit_size_in_bytes": 19180801228,
					"limit_size": "17.8gb",
					"estimated_size_in_bytes": 1000712,
					"estimated_size": "977.2kb",
					"overhead": 1.03,
					"tripped": 0
				},
				"in_flight_requests": {
					"limit_size_in_bytes": 31968002048,
					"limit_size": "29.7gb",
					"estimated_size_in_bytes": 151649,
					"estimated_size": "148kb",
					"overhead": 1.0,
					"tripped": 0
				},
				"parent": {
					"limit_size_in_bytes": 22377601433,
					"limit_size": "20.8gb",
					"estimated_size_in_bytes": 1363305,
					"estimated_size": "1.3mb",
					"overhead": 1.0,
					"tripped": 0
				}
			},
			"script": {
				"compilations": 1,
				"cache_evictions": 0
			},
			"discovery": {
				"cluster_state_queue": {
					"total": 0,
					"pending": 0,
					"committed": 0
				}
			},
			"ingest": {
				"total": {
					"count": 0,
					"time_in_millis": 0,
					"current": 0,
					"failed": 0
				},
				"pipelines": {}
			}
		}
	}
}
	`
)
