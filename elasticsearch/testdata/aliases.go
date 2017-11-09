package testdata

import (
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
)

// Test data for aliases info
var (
	AliasesBody = `{"twitter": {"aliases": {"alias1": {}, "alias2": {}}}}`

	Aliases = model.Aliases{
		"twitter": model.AliasInfo{
			Aliases: map[string]interface{}{
				"alias1": map[string]interface{}{},
				"alias2": map[string]interface{}{},
			},
		},
	}
)
