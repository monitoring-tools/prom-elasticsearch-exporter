package tasks

import (
	"log"
	"strings"
	"time"

	"github.com/monitoring-tools/prom-elasticsearch-exporter/elasticsearch"
	"github.com/monitoring-tools/prom-elasticsearch-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Labels lists for different kind of metrics
var (
	labelsTasks = []string{"action", "node", "cluster"}
)

// Label values extractors for different kind of metrics
var (
	labelValuesTasks = func(action, host, cluster string) []string {
		return []string{action, host, cluster}
	}
)

type taskMetric struct {
	*metrics.Metric
	Value func(t float64) float64
}

// Collector is an tasks metrics collector
type Collector struct {
	esClient   elasticsearch.Client
	taskMetric *taskMetric
	maxTasks   map[string]max
}

type max struct {
	time      float64
	taskIndex int
}

func newTaskMetric(t prometheus.ValueType, name, help string, valueExtractor func(t float64) float64) *taskMetric {
	return &taskMetric{
		Metric: metrics.New(t, "", name, help, labelsTasks),
		Value:  valueExtractor,
	}
}

// NewCollector returns new nodes metrics collector
func NewCollector(esClient elasticsearch.Client) *Collector {
	return &Collector{
		esClient: esClient,
		taskMetric: newTaskMetric(
			prometheus.GaugeValue, "task_group_duration_seconds", "Task group running duration in seconds",
			func(t float64) float64 { return t },
		),
		maxTasks: make(map[string]max),
	}
}

// Describe implements prometheus.Collector interface
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.taskMetric.Desc()
}

// Collect writes data to metrics channel
func (c *Collector) Collect(clusterName string, ch chan<- prometheus.Metric) {
	tasks, err := c.esClient.Tasks()
	if err != nil {
		log.Println("ERROR: failed to fetch nodes stats: ", err)
		return
	}

	for x := range c.maxTasks {
		delete(c.maxTasks, x)
	}

	// Select only max running time for particular task and node
	for i, task := range *tasks {
		t, err := time.ParseDuration(strings.Replace(task.RunningTime, "micros", "us", 1))
		if err != nil {
			log.Println("ERROR: failed to parse time: ", err)
			continue
		}

		taskTime := float64(t) / float64(time.Second)
		key := task.Action + "|" + task.Node
		maxItem, has := c.maxTasks[key]
		if !has || taskTime > maxItem.time {
			c.maxTasks[key] = max{time: taskTime, taskIndex: i}
		}
	}

	for _, max := range c.maxTasks {
		task := (*tasks)[max.taskIndex]
		ch <- prometheus.MustNewConstMetric(
			c.taskMetric.Desc(),
			c.taskMetric.Type(),
			c.taskMetric.Value(max.time),
			labelValuesTasks(task.Action, task.Node, clusterName)...,
		)
	}
}
