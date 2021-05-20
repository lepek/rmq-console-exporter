package collectors

import "github.com/oriser/regroup"

type ICollector interface {
	Collect() ([]IMetrics, error)
}

type IMetrics interface {
	GetMetricByName(name string) float64
	GetLabelsForMetric(name string) map[string]string
}

type IConsoleParser interface {
	GetCmd() string
	GetArguments() []string
	GetParser()	*regroup.ReGroup
	GetNewContainer() IMetrics
}