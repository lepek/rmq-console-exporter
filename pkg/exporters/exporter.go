package exporters

type Exporter interface {
	Init() error
	Collect() error
}