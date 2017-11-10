# Prometheus exporter for ElasticSearch
[![Build Status](https://travis-ci.org/monitoring-tools/prom-elasticsearch-exporter.svg?branch=master)](https://travis-ci.org/monitoring-tools/prom-elasticsearch-exporter)

### Installation

#### Pre-build binaries

You can find here: https://github.com/monitoring-tools/prom-elasticsearch-exporter/releases

#### Docker

```bash
docker pull montools/prom-elasticsearch-exporter
docker run --rm -p 9108:9108 montools/prom-elasticsearch-exporter:latest
```

### Build

#### 1. Clone repo

```bash
mkdir -p $GOPATH/src/github.com/monitoring-tools && cd $GOPATH/src/github.com/monitoring-tools
git clone https://github.com/monitoring-tools/prom-elasticsearch-exporter.git
cd prom-elasticsearch-exporter
```

#### 2. Build binary

```bash
make build
./prom-elasticsearch-exporter --es.uri="http://localhost:9200/"
```

#### 3. Build docker image

```bash
make docker
docker run -d -p 9108:9108 --name=prom-elasticsearch-exporter monitoring-tools/prom-elasticsearch-exporter:1.0.0 --es.uri="http://localhost.iddc:9200/"
```

### Configuration

```bash
prom-elasticsearch-exporter --help
```

| Argument              | Description |
| --------              | ----------- |
| web.listen-address    | Address to listen on for web interface and telemetry. Default - :9108 |
| web.telemetry-path    | Path under which to expose metrics. Default - /metrics |
| es.uri                | ElasticSearch URI. You can provide multiple hosts: --es.uri=host1 --es.uri=host2. If you're using multiple hosts and --es.all=true, metrics will be fetched from first responded node`.
| es.all                | If true - export stats for all nodes in the cluster. Default - false
| es.timeout            | Timeout for trying to get stats from ElasticSearch. Default - 5s |
| es.ca                 | Path to PEM file that contains trusted CAs for the Elasticsearch connection.
| es.client-private-key | Path to PEM file that contains the private key for client auth when connecting to Elasticsearch.
| es.client-cert        | Path to PEM file that contains the corresponding cert for the private key to connect to Elasticsearch.

### Grafana dashboards

To use this dashboards you need to set up following Prometheus [aggregation rules](examples/prometheus.rules).

Dashboards:
- [Cluster](examples/grafana-dashboard-cluster.json)
- [Index](examples/grafana-dashboard-index.json)
- [Node](examples/grafana-dashboard-node.json)

## License

This software is distributed under the terms of the [Apache License v2.0](LICENSE)

## Contributing

Contributions are very welcome.
Please fork the project on GitHub and open pull requests for any proposed changes.
