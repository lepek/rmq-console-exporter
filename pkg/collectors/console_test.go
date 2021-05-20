package collectors

import (
	"github.com/oriser/regroup"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestMetrics struct {
	BcBaseMax		float64	`regroup:"BcBaseMax"`
}

func (m *TestMetrics) GetMetricByName(name string) float64 {
	r := reflect.ValueOf(m)
	f := reflect.Indirect(r).FieldByName(name)
	return f.Float()
}

func (m *TestMetrics) GetLabelsForMetric(name string) map[string]string {
	labelPairs := make(map[string]string, 1)
	labelPairs["user"] = "Martin"
	return labelPairs
}

type TestParser struct {
	Cmd			string
	Arguments	[]string
	Parser		*regroup.ReGroup
}

func NewTestParser() *TestParser {
	return &TestParser{
		Cmd: "sysctl",
		Arguments: []string{"-a"},
		Parser: regroup.MustCompile(`^user.bc_base_max:\s+(?P<BcBaseMax>.+)`),
	}
}

func (p *TestParser) GetCmd() string {
	return p.Cmd
}

func (p *TestParser) GetArguments() []string {
	return p.Arguments
}

func (p *TestParser) GetParser() *regroup.ReGroup {
	return p.Parser
}

func (p *TestParser) GetNewContainer() IMetrics {
	return &TestMetrics{}
}

func TestCollectOk(t *testing.T) {
	parser := NewTestParser()
	console := NewConsoleSource(parser, 10000)
	results, err := console.Collect()
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(99), results[0].GetMetricByName("BcBaseMax"))
}
