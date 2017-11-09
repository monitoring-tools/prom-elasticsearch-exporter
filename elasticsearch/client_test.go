package elasticsearch

import (
	"reflect"
	"testing"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch/testdata"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
)

func TestClient_Aliases_Ok(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_aliases").WillReturn(200, testdata.AliasesBody)

	esClient := NewClient(mockHTTPClient)
	got, err := esClient.Aliases()

	if err != nil {
		t.Fatalf("Error on getting ES aliases: %s", err)
	}
	if !reflect.DeepEqual(testdata.Aliases, got) {
		t.Fatalf("Structs are not equal: want %+v, got %+v", testdata.Aliases, got)
	}
}

func TestClient_Aliases_Error(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_aliases").WillReturn(500, ``)

	esClient := NewClient(mockHTTPClient)
	_, err := esClient.Aliases()

	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
}

func TestClient_ClusterHealth_Ok(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_cluster/health?level=indices").WillReturn(200, testdata.ClusterHealthIndicesBody)

	esClient := NewClient(mockHTTPClient)
	got, err := esClient.ClusterHealth(LevelIndices)

	if err != nil {
		t.Fatalf("Error on getting ES cluster health: %s", err)
	}

	if !reflect.DeepEqual(testdata.ClusterHealthIndices, got) {
		t.Fatalf("Structs are not equal: want %+v, got %+v", testdata.ClusterHealthIndices, got)
	}
}

func TestClient_ClusterHealth_Error(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_cluster/health?level=cluster").WillReturn(500, ``)

	esClient := NewClient(mockHTTPClient)
	_, err := esClient.ClusterHealth(LevelCluster)

	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
}

func TestClient_NodesAll_Ok(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_nodes/stats").WillReturn(200, testdata.NodesBody)

	esClient := NewClient(mockHTTPClient)
	nodes, err := esClient.Nodes(true)

	if err != nil {
		t.Fatalf("Error on getting ES nodes stats: %s", err)
	}

	if len(nodes.Nodes) != 1 {
		t.Fatalf("Unexpected nodes count, wat 1, get: %d", len(nodes.Nodes))
	}
}

func TestClient_NodesAll_Error(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_nodes/stats").WillReturn(500, ``)

	esClient := NewClient(mockHTTPClient)
	_, err := esClient.Nodes(true)

	if err == nil {
		t.Fatalf("Error expected, got nil")
	}
}

func TestClient_NodesSelf_Ok(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_nodes/_local/stats").WillReturn(200, testdata.NodesBody)

	esClient := NewClient(mockHTTPClient)
	_, err := esClient.Nodes(false)

	if err != nil {
		t.Fatalf("Error on getting ES nodes stats: %s", err)
	}
}

func TestClient_Recovery_Ok(t *testing.T) {
	mockHTTPClient := httpclient.NewClient()
	mockHTTPClient.Get("/_recovery?active_only=true").WillReturn(200, testdata.RecoveryBody)

	esClient := NewClient(mockHTTPClient)
	got, err := esClient.Recovery()

	if err != nil {
		t.Fatalf("Error on getting ES recovery state: %s", err)
	}

	if !reflect.DeepEqual(testdata.Recovery, got) {
		t.Fatalf("Structs are not equal: want %+v, got %+v", testdata.Recovery, got)
	}
}
