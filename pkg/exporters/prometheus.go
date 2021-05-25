package exporters

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"rmq-console-exporter/pkg/collectors"
	"strconv"
)

type PrometheusExporter struct {
	MetricsDesc		map[string]*prometheus.Desc
	Port			int
	RMQCollector	[]collectors.ICollector
	MetricLabels	[]string
}

func NewPrometheusExporter(prefix string, port int, collector []collectors.ICollector) *PrometheusExporter {
	labels := []string{"queue", "state"}
	return &PrometheusExporter{
		MetricsDesc: createPrometheusMetrics(prefix, labels),
		Port: port,
		RMQCollector: collector,
		MetricLabels: labels,
	}
}

func (p *PrometheusExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metricDesc := range p.MetricsDesc {
		ch <- metricDesc
	}
}

func (p *PrometheusExporter) Init() error {
	prometheus.MustRegister(p)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":" + strconv.Itoa(p.Port)), nil); err != nil {
		return err
	}
	return nil
}

func (p *PrometheusExporter) Collect(ch chan<- prometheus.Metric) {
	for _, collector := range p.RMQCollector {
		metrics, err := collector.Collect()
		if err != nil {
			log.Errorf("Metrics collection has failed for collector %v: %v", collector, err)
			continue
		}
		for metricName, pDesc := range p.MetricsDesc {
			for _, queueMetrics := range metrics {
				if queueMetricValue, err := queueMetrics.GetMetricValue(metricName); err == nil {
					metricLabels, err := queueMetrics.GetLabels(metricName)
					if err != nil { continue } // No labels
					labels := p.buildLabels(metricLabels)
					constMetric, err := prometheus.NewConstMetric(pDesc, prometheus.GaugeValue, queueMetricValue, labels...)
					if err != nil {
						log.Errorf("Error building metric for %s: %v", metricName, err)
						continue
					}
					ch <- constMetric
				}
			}
		}
	}
}

func (p *PrometheusExporter) buildLabels(labelPairs map[string]string) []string {
	labels := make([]string, len(p.MetricLabels))
	for i, labelName := range p.MetricLabels {
		labels[i] = labelPairs[labelName]
	}
	return labels
}

func createPrometheusMetrics(prefix string, labels []string) map[string]*prometheus.Desc {
	pMetrics := make(map[string]*prometheus.Desc)

	pMetrics["MessagesReady"] = prometheus.NewDesc(
		prefix + "messages_ready",
		"Number of messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics["MessageBytesReady"] = prometheus.NewDesc(
		prefix + "message_bytes_ready",
		"Like message_bytes but counting only those messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics["MessagesUnacknowledged"] = prometheus.NewDesc(
		prefix + "messages_unacknowledged",
		"Like message_bytes but counting only those messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics["MessageBytesUnacknowledged"] = prometheus.NewDesc(
		prefix + "message_bytes_unacknowledged",
		"Like message_bytes but counting only those messages delivered to clients but not yet acknowledged.",
		labels,
		nil,
	)

	pMetrics["Memory"] = prometheus.NewDesc(
		prefix + "memory",
		"Bytes of memory allocated by the runtime for the queue, including stack, heap and internal structures.",
		labels,
		nil,
	)

	pMetrics["Consumers"] = prometheus.NewDesc(
		prefix + "consumers",
		"Number of consumers.",
		labels,
		nil,
	)

	pMetrics["ConsumerUtilisation"] = prometheus.NewDesc(
		prefix + "consumer_utilisation",
		"Fraction of the time (between 0.0 and 1.0) that the queue is able to immediately deliver messages to " +
			"consumers. This can be less than 1.0 if consumers are limited by network congestion or prefetch count.",
		labels,
		nil,
	)

	return pMetrics
}