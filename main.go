package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"rmq-console-exporter/pkg/collectors"
	"rmq-console-exporter/pkg/exporters"
)

func main() {
	port := flag.Int("port", 2112, "Port to expose metrics")
	prefix := flag.String("prefix", "rmq_", "Metrics prefix")
	timeoutMs := flag.Int("timeout", 60000, "Timeout[Ms] for each collector")
	outputBufferLines := flag.Int("output_buffer", 100000, "Output Buffer[lines]")
	level := flag.String("log_level", "info", "Log Level: debug, info, error, etc")
	qParser := flag.String("queue_parser", "json", "Queue Parser to use: json or tabular")
	flag.Parse()

	configureLogLevel(*level)
	log.Infof("Log Level set to %s", log.GetLevel().String())

	log.Infof("Collector agent starting...")

	executorFactory := collectors.NewExecutorFactory()
	queueParser := queueParserFactory(*qParser)
	queueCollector := collectors.NewCmdCollector(queueParser, executorFactory, *timeoutMs, *outputBufferLines)

	var rmqCollectors []exporters.ICollector
	rmqCollectors = append(rmqCollectors, queueCollector)
	exporter := exporters.NewPrometheusExporter(*prefix, *port, rmqCollectors)

	log.Infof("Collector agent running")
	log.Fatal(exporter.Init())
}

func configureLogLevel(strLogLevel string) {
	logLevel, err := log.ParseLevel(strLogLevel)
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
}

func queueParserFactory(strParser string) collectors.ICmdParser {
	if strParser == "tabular" {
		return collectors.NewQueueParser()
	}

	return collectors.NewQueueJSONParser()
}