package exporters

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"reflect"
	"rmq-console-exporter/pkg/collectors"
	"strconv"
)

type PrometheusMetrics struct {
	MessagesReady				*prometheus.Desc
	MessageBytesReady			*prometheus.Desc
	MessagesUnacknowledged		*prometheus.Desc
	MessageBytesUnacknowledged	*prometheus.Desc
	Memory						*prometheus.Desc
	Consumers					*prometheus.Desc
	ConsumerUtilisation			*prometheus.Desc
}

func (pm *PrometheusMetrics) GetDescByName(name string) *prometheus.Desc {
	r := reflect.ValueOf(pm)
	f := reflect.Indirect(r).FieldByName(name).Interface()
	desc, ok := f.(*prometheus.Desc)
	if !ok { return nil }
	return desc
}

type PrometheusExporter struct {
	MetricsDesc		*PrometheusMetrics
	Port			int
	RMQCollector	collectors.ICollector
	MetricLabels	[]string
}

func NewPrometheusExporter(prefix string, port int, collector collectors.ICollector) *PrometheusExporter {
	labels := []string{"queue", "state"}
	return &PrometheusExporter{
		MetricsDesc: createPrometheusMetrics(prefix, labels),
		Port: port,
		RMQCollector: collector,
		MetricLabels: labels,
	}
}

func (p *PrometheusExporter) Describe(ch chan<- *prometheus.Desc) {
	metricsDesc := structs.Values(p.MetricsDesc)
	for _, metricDesc := range metricsDesc {
		metricDesc, ok := metricDesc.(*prometheus.Desc)
		if ok { ch <- metricDesc }
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
	metrics, err := p.RMQCollector.Collect()
	if err != nil {
		log.Errorf("Metrics collection has failed: ", err)
	}
	for _, metric := range metrics {
		metricNameDescPairs := structs.Map(p.MetricsDesc)
		for metricName, metricDesc := range metricNameDescPairs {
			metricDesc := metricDesc.(*prometheus.Desc)
			metricValue := metric.GetMetricByName(metricName)
			labels := p.buildLabels(metric.GetLabelsForMetric(metricName))

			constMetric, err := prometheus.NewConstMetric(metricDesc, prometheus.GaugeValue, metricValue, labels...)
			if err != nil {
				log.Errorf("Error building metric for %s: %v", metricName, err)
				continue
			}
			ch <- constMetric
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

func createPrometheusMetrics(prefix string, labels []string) *PrometheusMetrics {
	pMetrics := &PrometheusMetrics{}

	pMetrics.MessagesReady = prometheus.NewDesc(
		prefix + "messages_ready",
		"Number of messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics.MessageBytesReady = prometheus.NewDesc(
		prefix + "message_bytes_ready",
		"Like message_bytes but counting only those messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics.MessagesUnacknowledged = prometheus.NewDesc(
		prefix + "messages_unacknowledged",
		"Like message_bytes but counting only those messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics.MessageBytesUnacknowledged = prometheus.NewDesc(
		prefix + "message_bytes_unacknowledged",
		"Like message_bytes but counting only those messages delivered to clients but not yet acknowledged.",
		labels,
		nil,
	)

	pMetrics.Memory = prometheus.NewDesc(
		prefix + "memory",
		"Bytes of memory allocated by the runtime for the queue, including stack, heap and internal structures.",
		labels,
		nil,
	)

	pMetrics.Consumers = prometheus.NewDesc(
		prefix + "consumers",
		"Number of consumers.",
		labels,
		nil,
	)

	pMetrics.ConsumerUtilisation = prometheus.NewDesc(
		prefix + "consumer_utilisation",
		"Fraction of the time (between 0.0 and 1.0) that the queue is able to immediately deliver messages to " +
			"consumers. This can be less than 1.0 if consumers are limited by network congestion or prefetch count.",
		labels,
		nil,
	)

	return pMetrics
}