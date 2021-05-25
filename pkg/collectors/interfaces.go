package collectors

type ICollector interface {
	Collect() ([]IMetrics, error)
}

type IMetrics interface {
	GetMetricValue(name string) (float64, error)
	GetLabels(name string) (map[string]string, error)
}

type ICmdParser interface {
	GetCmd() string
	GetArguments() []string
	Parse(string) (IMetrics, error)
}