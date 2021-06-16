package main

import (
	"flag"
	"github.com/prometheus/common/log"
	"rmq-console-exporter/pkg/collectors"
	"rmq-console-exporter/pkg/exporters"
	"rmq-console-exporter/pkg/parsers"
)

func main() {
	port := flag.Int("port", 2112, "Port to expose metrics")
	prefix := flag.String("prefix", "rmq_", "Metrics prefix")
	timeoutMs := flag.Int("timeout", 30000, "Timeout[Ms] for each collector")
	flag.Parse()

	var rmqCollectors []collectors.ICollector
	queueCollector := collectors.NewCmdCollector(parsers.NewQueueParser(), *timeoutMs)
	rmqCollectors = append(rmqCollectors, queueCollector)

	exporter := exporters.NewPrometheusExporter(*prefix, *port, rmqCollectors)
	log.Fatal(exporter.Init())
}