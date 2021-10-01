package main

import (
	"flag"
	"github.com/prometheus/common/log"
	"rmq-console-exporter/pkg/collectors"
	"rmq-console-exporter/pkg/exporters"
)

func main() {
	port := flag.Int("port", 2112, "Port to expose metrics")
	prefix := flag.String("prefix", "rmq_", "Metrics prefix")
	timeoutMs := flag.Int("timeout", 30000, "Timeout[Ms] for each collector")
	outputBufferLines := flag.Int("output_buffer", 100000, "Output Buffer[lines]")
	flag.Parse()

	var rmqCollectors []exporters.ICollector

	queueParser := collectors.NewQueueParser()
	executorFactory := collectors.NewExecutorFactory()
	queueCollector := collectors.NewCmdCollector(queueParser, executorFactory, *timeoutMs, *outputBufferLines)

	rmqCollectors = append(rmqCollectors, queueCollector)
	exporter := exporters.NewPrometheusExporter(*prefix, *port, rmqCollectors)

	log.Fatal(exporter.Init())
}