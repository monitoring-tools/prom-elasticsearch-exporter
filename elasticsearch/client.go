package elasticsearch

import (
	"encoding/json"
	"net/http"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/model"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
)

type clusterHealthLevel string

// Cluster health levels
const (
	LevelCluster clusterHealthLevel = "cluster"
	LevelIndices clusterHealthLevel = "indices"
	LevelShards  clusterHealthLevel = "shards"
)

// IClient is an ElasticSearch client interface
type IClient interface {
	ClusterHealth(level clusterHealthLevel) (*model.ClusterHealth, error)
	Aliases() (model.Aliases, error)
	Indices() (*model.Indices, error)
	Nodes(fetchAllNodesInfo bool) (*model.Nodes, error)
	Recovery() (model.Recovery, error)
}

var (
	_ IClient = &Client{}
)

// NewClient returns new client
func NewClient(httpClient httpclient.IClient) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// Client is an ElasticSearch client implementation
type Client struct {
	httpClient httpclient.IClient
}

// ClusterHealth returns ES cluster health info
func (c *Client) ClusterHealth(level clusterHealthLevel) (*model.ClusterHealth, error) {
	var v model.ClusterHealth
	if err := c.makeRequest("/_cluster/health?level="+string(level), &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Aliases returns ES index aliases info
func (c *Client) Aliases() (model.Aliases, error) {
	var v model.Aliases
	if err := c.makeRequest("/_aliases", &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Indices returns ES indices info
func (c *Client) Indices() (*model.Indices, error) {
	var v model.Indices
	if err := c.makeRequest("/_stats", &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Nodes returns ES nodes info
func (c *Client) Nodes(fetchAllNodesInfo bool) (*model.Nodes, error) {
	path := "/_nodes/_local/stats"
	if fetchAllNodesInfo {
		path = "/_nodes/stats"
	}

	var v model.Nodes
	if err := c.makeRequest(path, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Recovery returns ES state with currently active recovery operations
func (c *Client) Recovery() (model.Recovery, error) {
	var v model.Recovery
	path := "/_recovery?active_only=true"

	if err := c.makeRequest(path, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// makeRequest sends request and encodes it to given struct
func (c *Client) makeRequest(path string, v interface{}) error {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
