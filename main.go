package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"os"
	"rmq-console-exporter/pkg/collectors"
	"rmq-console-exporter/pkg/exporters"
	utils "rmq-console-exporter/pkg/internal_utils"
)

func main() {
	port := flag.Int("port", 2112, "Port to expose metrics")
	prefix := flag.String("prefix", "rmq_", "Metrics prefix")
	timeoutMs := flag.Int("timeout", 600000, "Timeout[Ms] for each collector")
	outputBufferLines := flag.Int("output_buffer", 100000, "Output Buffer[lines]")
	level := flag.String("log_level", "info", "Log Level: debug, info, error, etc")
	qParser := flag.String("queue_parser", "json", "Queue Parser to use: json or tabular")
	configFilePath := flag.String("config_file", "", "Config file (use the flag -create_config to create one)")
	createConfig := flag.Bool("create_config", false, "Lunch the tool to create a config file")
	flag.Parse()

	configureLogLevel(*level)
	log.Infof("Log Level set to %s", log.GetLevel().String())

	if *createConfig {
		err := utils.CreateConfig()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	log.Infof("Collector agent starting...")

	config := loadConfig(*configFilePath)
	executorFactory := collectors.NewExecutorFactory()
	queueParser := queueParserFactory(*qParser, config)
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

func queueParserFactory(strParser string, config collectors.IConfig) collectors.ICmdParser {
	if strParser == "tabular" {
		return collectors.NewQueueParser(config)
	}

	return collectors.NewQueueJSONParser(config)
}

func loadConfig(configFilePath string) collectors.IConfig {
	config, err := collectors.NewConfig(configFilePath)
	if err != nil {
		log.Warningf("error loading config: %v", err)
		return config
	}

	if !config.IsEmpty() {
		config.OnConfigChange(func(e fsnotify.Event) {
			log.Infof("Config file changed: %s", e.Name)
			if err := config.Init(); err != nil {
				log.Fatal(err)
			}
			log.Info("Config reloaded")
		})
		config.WatchConfig()
	}

	return config
}