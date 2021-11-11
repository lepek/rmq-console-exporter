package exporters

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/tevino/abool"
	"net/http"
	"strconv"
)

type IMetrics interface {
	GetMetricValue(name string) (float64, error)
	GetLabels(name string) (map[string]string, error)
}

type ICollector interface {
	Collect() ([]IMetrics, error)
}

var (
	/*
	 *  Since the collection operation is time very consuming,
	 *  we will allow only one collector running at any time to avoid multiple errors:
	 *  - Multiple threads sending the same metric will trigger a duplicated metric error on Prometheus library
	 *  - Multiple threads hammering the RMQ server for metrics
	 *  - One scrape starting when the previous one hasn't finished yet
	 */
	isRunning = abool.New()
)

type PrometheusExporter struct {
	MetricsDesc		map[string]*prometheus.Desc
	Port			int
	RMQCollector	[]ICollector
	MetricLabels	[]string
}

func NewPrometheusExporter(prefix string, port int, collector []ICollector) *PrometheusExporter {
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
	return http.ListenAndServe(fmt.Sprintf(":" + strconv.Itoa(p.Port)), nil)
}

func (p *PrometheusExporter) Collect(ch chan<- prometheus.Metric) {
	if isRunning.IsSet() {
		log.Error("A collection is running, skipping new collection...")
		return
	}
	isRunning.Set()
	defer func() {
		isRunning.UnSet()
		log.Info("Metrics collection finished")
	}()

	log.Info("Starting metrics collection")
	for _, collector := range p.RMQCollector {
		var metrics []IMetrics
		var err error
		metrics, err = collector.Collect()
		if err != nil {
			log.Errorf("Metrics collection has failed for collector %v: %v", collector, err)
			continue
		}
		log.Infof("Metrics collected from >> %v << objects. Starting building metrics...", len(metrics))
		for metricName, pDesc := range p.MetricsDesc {
			for _, queueMetrics := range metrics {
				if queueMetricValue, err := queueMetrics.GetMetricValue(metricName); err == nil {
					metricLabels, err := queueMetrics.GetLabels(metricName)
					var labels []string
					if err == nil { labels = p.buildLabels(metricLabels) }
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
	// TBD: Redesign how different metrics can have different labels
	if len(labelPairs) == 1 {
		for _, label := range labelPairs {
			return []string{label}
		}
	}

	labels := make([]string, len(p.MetricLabels))
	for i, labelName := range p.MetricLabels {
		labels[i] = labelPairs[labelName]
	}
	return labels
}

func createPrometheusMetrics(prefix string, labels []string) map[string]*prometheus.Desc {
	pMetrics := make(map[string]*prometheus.Desc)

	pMetrics["messages_ready"] = prometheus.NewDesc(
		prefix + "messages_ready",
		"Number of messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics["message_bytes_ready"] = prometheus.NewDesc(
		prefix + "message_bytes_ready",
		"Like message_bytes but counting only those messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics["messages_unacknowledged"] = prometheus.NewDesc(
		prefix + "messages_unacknowledged",
		"Like message_bytes but counting only those messages ready to be delivered to clients.",
		labels,
		nil,
	)

	pMetrics["message_bytes_unacknowledged"] = prometheus.NewDesc(
		prefix + "message_bytes_unacknowledged",
		"Like message_bytes but counting only those messages delivered to clients but not yet acknowledged.",
		labels,
		nil,
	)

	pMetrics["memory"] = prometheus.NewDesc(
		prefix + "memory",
		"Bytes of memory allocated by the runtime for the queue, including stack, heap and internal structures.",
		labels,
		nil,
	)

	pMetrics["consumers"] = prometheus.NewDesc(
		prefix + "consumers",
		"Number of consumers.",
		labels,
		nil,
	)

	pMetrics["consumer_utilisation"] = prometheus.NewDesc(
		prefix + "consumer_utilisation",
		"Fraction of the time (between 0.0 and 1.0) that the queue is able to immediately deliver messages to " +
			"consumers. This can be less than 1.0 if consumers are limited by network congestion or prefetch count.",
		labels,
		nil,
	)

	pMetrics["head_message_timestamp"] = prometheus.NewDesc(
		prefix + "head_message_timestamp",
		"The timestamp property of the first message in the queue, if present. " +
			"Timestamps of messages only appear when they are in the paged-in state.",
		labels,
		nil,
	)

	pMetrics["command_runtime"] = prometheus.NewDesc(
		prefix + "command_runtime_seconds",
		"Runtime of the command executed to collect the metrics.",
		[]string{"command_executed"},
		nil,
	)

	return pMetrics
}