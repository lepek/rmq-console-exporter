package collectors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestParserOk struct {
	Cmd			string
	Arguments	[]string
}

func NewTestParserOk() *TestParserOk {
	return &TestParserOk{
		Cmd: "ls",
		Arguments: []string{"/"},
	}
}

func (p *TestParserOk) GetCmd() string {
	return p.Cmd
}

func (p *TestParserOk) GetArguments() []string {
	return p.Arguments
}

func (p *TestParserOk) Parse(line string) (IMetrics, error) {
	queueMetrics := NewMetrics()
	queueMetrics.AddMetric("test_metric", 1.5, map[string]string{"testLabel": "testing labels"})
	return queueMetrics, nil
}

func TestCollectOk(t *testing.T) {
	parser := NewTestParserOk()
	console := NewCmdCollector(parser, 1000)
	results, err := console.Collect()
	assert.Equal(t, nil, err)
	value, _ := results[0].GetMetricValue("test_metric")
	assert.Equal(t, 1.5, value)
}

//********************************************************************************************************************//

type TestParserFail struct {
	Cmd			string
	Arguments	[]string
}

func NewTestParserFail() *TestParserFail {
	return &TestParserFail{
		Cmd: "ls",
		Arguments: []string{"/"},
	}
}

func (p *TestParserFail) GetCmd() string {
	return p.Cmd
}

func (p *TestParserFail) GetArguments() []string {
	return p.Arguments
}

func (p *TestParserFail) Parse(line string) (IMetrics, error) {
	return nil, errors.New("testing a failing parser")
}

func TestCollectFail(t *testing.T) {
	parser := NewTestParserFail()
	console := NewCmdCollector(parser, 1000)
	results, err := console.Collect()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, len(results))
}