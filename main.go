package main

import (
	"fmt"
	"rmq-console-exporter/pkg/collectors"
	"rmq-console-exporter/pkg/exporters"
)

func main() {
	consoleSource := collectors.NewConsoleSource(collectors.NewConsoleQueueParser(), 10000)
	metrics, err := consoleSource.Collect()
	fmt.Println(err)
	for _, metric := range metrics {
		fmt.Println(metric)
	}
	exporter := exporters.NewPrometheusExporter("rmq_", 2112, consoleSource)
	exporter.Init()
}