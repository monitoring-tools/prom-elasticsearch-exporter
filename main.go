package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/collector"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/encryption"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/httpclient/decorator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const usage = `prom-elasticsearch-exporter - tool for exporting ElasticSearch metrics for Prometheus

Usage:

  prom-elasticsearch-exporter [commands|flags]

The commands & flags are:

  version                   print the version to stdout

  --web.listen-address      address to listen on for web interface and telemetry. Default - :9108
  --web.telemetry-path      path under which to expose metrics. Default - /metrics
  --es.timeout              timeout for trying to get stats from ElasticSearch. Default - 5s
  --es.uri                  ElasticSearch URI. You can provide multiple hosts: --es.uri=host1 --es.uri=host2
  --es.all                  export stats for all nodes in the cluster. Default - false
  --es.ca                   path to PEM file that conains trusted CAs for the ElasticSearch connection
  --es.client-private-key   path to PEM file that conains the private key for client auth when connecting to ElasticSearch
  --es.client-cert          path to PEM file that conains the corresponding cert for the private key to connect to Elasticsearch
`

// Variables passed through ldflags
var (
	version   string
	goVersion string
	gitBranch string
)

func main() {
	var (
		listenAddress      = flag.String("web.listen-address", ":9108", "Address to listen on for web interface and telemetry.")
		metricsPath        = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		esTimeout          = flag.Duration("es.timeout", 5*time.Second, "Timeout for trying to get stats from ElasticSearch.")
		esAllNodes         = flag.Bool("es.all", false, "Export stats for all nodes in the cluster.")
		esCA               = flag.String("es.ca", "", "Path to PEM file that conains trusted CAs for the ElasticSearch connection.")
		esClientPrivateKey = flag.String("es.client-private-key", "", "Path to PEM file that conains the private key for client auth when connecting to ElasticSearch.")
		esClientCert       = flag.String("es.client-cert", "", "Path to PEM file that conains the corresponding cert for the private key to connect to ElasticSearch.")

		esURI stringSliceValue
	)

	flag.Var(&esURI, "es.uri", "HTTP API address of an Elasticsearch node.")

	flag.Usage = func() { printUsage() }
	flag.Parse()

	if len(esURI) == 0 {
		esURI = stringSliceValue([]string{"http://localhost:9200"})
	}

	args := flag.Args()

	if len(args) > 0 {
		switch args[0] {
		case "version":
			if version == "" {
				println("n/a")
			}
			fmt.Println(version)
			return
		case "help":
			printUsage()
		}
	}

	// returns nil if not provided and falls back to simple TCP.
	tlsConfig := encryption.CreateTLSConfig(*esCA, *esClientCert, *esClientPrivateKey)

	httpClient := &http.Client{
		Timeout: *esTimeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	decoratedClient := httpclient.Decorate(
		httpClient,
		decorator.BaseURLDecorator([]string(esURI)),
		decorator.RecoverDecorator(), // better to place it last to recover panics from decorators too
	)

	prometheus.MustRegister(collector.NewCompositeCollector(
		elasticsearch.NewClient(decoratedClient),
		*esAllNodes,
		version, goVersion, gitBranch,
	))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", IndexHandler(*metricsPath))

	log.Println("Listening on:", formatListenAddr(*listenAddress))

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Println("Unable to start HTTP server:", err)
	}
}

// IndexHandler returns a http handler with the correct metricsPath
func IndexHandler(metricsPath string) http.HandlerFunc {
	indexHTML := `
<html>
	<head>
		<title>Elasticsearch Exporter</title>
	</head>
	<body>
		<h1>Elasticsearch Exporter</h1>
		<p>
			<a href='%s'>Metrics</a>
		</p>
	</body>
</html>
`
	index := []byte(fmt.Sprintf(strings.TrimSpace(indexHTML), metricsPath))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	}
}

// stringSliceValue is a value used for array of string flags
type stringSliceValue []string

// String returns stringified representation of string slice
func (s *stringSliceValue) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set sets the flag value
func (s *stringSliceValue) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// formatListenAddr returns formatted UNIX addr
func formatListenAddr(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) == 2 && parts[0] == "" {
		addr = fmt.Sprintf("localhost:%s", parts[1])
	}
	return "http://" + addr
}

func printUsage() {
	fmt.Println(usage)
	os.Exit(0)
}
