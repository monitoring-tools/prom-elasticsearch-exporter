package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "elasticsearch"
)

// New creates new metric
func New(t prometheus.ValueType, subsystem, name, help string, labels []string) *Metric {
	return &Metric{
		_type: t,
		desc:  NewDesc(subsystem, name, help, labels, nil),
	}
}

// NewWithConstLabels returns new metric with const labels
func NewWithConstLabels(t prometheus.ValueType, subsystem, name, help string, labels []string, constLabels prometheus.Labels) *Metric {
	return &Metric{
		_type: t,
		desc:  NewDesc(subsystem, name, help, labels, constLabels),
	}
}

// NewDesc returns new metric description in terms of Prometheus client
func NewDesc(subsystem, name, help string, labels []string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		help,
		labels, constLabels,
	)
}

// Metric is a base metric struct with common fields
type Metric struct {
	_type prometheus.ValueType
	desc  *prometheus.Desc
}

// Type returns Prometheus metric type
func (m *Metric) Type() prometheus.ValueType {
	return m._type
}

// Desc returns Prometheus metric desc
func (m *Metric) Desc() *prometheus.Desc {
	return m.desc
}
