package elasticsearch

import (
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
)

// ClientMock is a client mock implementation
type ClientMock struct {
	ClusterHealthCallback func(level clusterHealthLevel) (*model.ClusterHealth, error)
	AliasesCallback       func() (model.Aliases, error)
	IndicesCallback       func() (*model.Indices, error)
	NodesCallback         func(fetchAllNodesInfo bool) (*model.Nodes, error)
}
