package testdata

import "github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"

// Test data for recovery info
var (
	RecoveryBody = `
{
	"my_awsome_index": {
		"shards": [{
			"id": 9,
			"type": "PEER",
			"stage": "INDEX",
			"primary": false,
			"start_time_in_millis": 1504261908035,
			"total_time_in_millis": 114573,
			"source": {
				"id": "JhO_zXsHRhKARQm9I20mUw",
				"host": "10.36.8.103",
				"transport_address": "10.36.8.103:9301",
				"ip": "10.36.8.103",
				"name": "old-host"
			},
			"target": {
				"id": "a3ZUS9r2RHyn0b4UeTSX0g",
				"host": "10.129.104.32",
				"transport_address": "10.129.104.32:9300",
				"ip": "10.129.104.32",
				"name": "new-host"
			},
			"index": {
				"size": {
					"total_in_bytes": 27181211342,
					"reused_in_bytes": 0,
					"recovered_in_bytes": 15359042167,
					"percent": "56.5%"
				},
				"files": {
					"total": 513,
					"reused": 0,
					"recovered": 504,
					"percent": "98.2%"
				},
				"total_time_in_millis": 114569,
				"source_throttle_time_in_millis": 7771,
				"target_throttle_time_in_millis": 7100
			},
			"translog": {
				"recovered": 0,
				"total": 10664,
				"percent": "0.0%",
				"total_on_start": 10081,
				"total_time_in_millis": 0
			},
			"verify_index": {
				"check_index_time_in_millis": 0,
				"total_time_in_millis": 0
			}
		}]
	}
}`

	Recovery = model.Recovery{
		"my_awsome_index": model.IndexRecovery{
			Shards: []model.ShardRecovery{
				{
					ID:                9,
					Type:              "PEER",
					Stage:             "INDEX",
					StartTimeInMillis: 1504261908035,
					TotalTimeInMillis: 114573,
					Source: model.RecoveryDestination{
						ID:               "JhO_zXsHRhKARQm9I20mUw",
						Host:             "10.36.8.103",
						TransportAddress: "10.36.8.103:9301",
						IP:               "10.36.8.103",
						Name:             "old-host",
					},
					Target: model.RecoveryDestination{
						ID:               "a3ZUS9r2RHyn0b4UeTSX0g",
						Host:             "10.129.104.32",
						TransportAddress: "10.129.104.32:9300",
						IP:               "10.129.104.32",
						Name:             "new-host",
					},
					Index: model.IndexRecoveryState{
						Size: model.IndexRecoverySize{
							TotalInBytes:     27181211342,
							RecoveredInBytes: 15359042167,
							Percent:          "56.5%",
						},
						Files: model.IndexRecoveryFiles{
							Total:     513,
							Recovered: 504,
							Percent:   "98.2%",
						},
						TotalTimeInMillis:          114569,
						SourceThrottleTimeInMillis: 7771,
						TargetThrottleTimeInMillis: 7100,
					},
					Translog: model.IndexRecoveryTranslog{
						Total:        10664,
						Percent:      "0.0%",
						TotalOnStart: 10081,
					},
				},
			},
		},
	}
)
