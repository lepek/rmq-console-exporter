package collectors

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestExecutor struct {
	outputCh		chan string
	endExecutionCh	chan struct{}
}

func (e *TestExecutor) Output() <-chan string {
	return e.outputCh
}

func (e *TestExecutor) Execute(ctx context.Context) error {
	defer close(e.outputCh)

	output := []string{
		"Timeout: 60.0 seconds ...",
		"Listing queues for vhost / ...",
		"name	state	messages_ready	message_bytes_ready	messages_unacknowledged	message_bytes_unacknowledged	memory	consumers	consumer_utilisation	head_message_timestamp",
		"delegate_encryption_test_3579441e-1f41-4455-90e4-04c3228f1305.tenant_3667d578-644d-4930-b965-4f7bd45ee537.dev	running	1	288	0	0	34764	0		1630920836",
		"SemanticsSystemTestQueue	running	0	0	0	0	34668	0",
		"DeletedExchangeSystemTest-nrm-eb89e741-aa61-4fc9-a350-a07d54c07e1c.dev	running	0	0	0	0	34764	0",
	}
	currLine := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			e.outputCh <- output[currLine]
			currLine++
			if currLine == len(output) {
				return nil
			}
		}
	}
}

type TestExecutorFactory struct {}

func NewTestExecutorFactory() *TestExecutorFactory {
	return &TestExecutorFactory{}
}

func (f *TestExecutorFactory) NewExecutor(command string, arguments []string, outputBuffer int) IExecutor {
	return &TestExecutor{
		outputCh: make(chan string, 100),
		endExecutionCh: make(chan struct{}, 1),
	}
}

//============== TEST ================ //
func TestCollectOk(t *testing.T) {
	console := NewCmdCollector(NewQueueParser(&TrueFilterConfig{}), NewTestExecutorFactory(), 1000000, 1000000)
	results, err := console.Collect()

	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(results))

	value, err := results[0].GetMetricValue("message_bytes_ready")
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(288), value)

	value, err = results[2].GetMetricValue("memory")
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(34764), value)
}

//********************************************************************************************************************//

type TestParserFail struct {
	Cmd			string
	Arguments	[]string
}

func NewTestParserFail() *TestParserFail {
	return &TestParserFail{}
}

func (p *TestParserFail) GetCmd() string {
	return ""
}

func (p *TestParserFail) GetArguments() []string {
	return []string{}
}

func (p *TestParserFail) Parse(line string) (*Metrics, error) {
	return nil, errors.New("testing a failing parser")
}

//============== TEST ================ //
func TestCollectFail(t *testing.T) {
	parser := NewTestParserFail()
	console := NewCmdCollector(parser, NewTestExecutorFactory(), 1000000, 100000)
	results, err := console.Collect()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, len(results))
}

//********************************************************************************************************************//

type TestJSONExecutor struct {
	outputCh		chan string
	endExecutionCh	chan struct{}
}

func (e *TestJSONExecutor) Output() <-chan string {
	return e.outputCh
}

func (e *TestJSONExecutor) Execute(ctx context.Context) error {
	defer close(e.outputCh)

	output := []string{
		`[`,
		`{"name":"10_128_4_241:5672.sage-xds-service.super_after157k_haneeshp.LATEST.perf","state":"running","messages_ready":0,"message_bytes_ready":0,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":55692,"consumers":1,"consumer_utilisation":1.0,"head_message_timestamp":""}`,
		`,{"name":"10_128_4_241:5672.sage-xds-service.performance_v20210920_0921_omsperf.LATEST.perf","state":"running","messages_ready":0,"message_bytes_ready":0,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":55692,"consumers":1,"consumer_utilisation":1.0,"head_message_timestamp":""}`,
		`,{"name":"10_128_4_241:5672.sage-xds-service.35k_sp_20210924_7.LATEST.perf","state":"running","messages_ready":0,"message_bytes_ready":0,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":55692,"consumers":1,"consumer_utilisation":1.0,"head_message_timestamp":""}`,
		`,{"name":"10_128_4_241:5672.sage-xds-service.super_0110_gpcc.LATEST.perf","state":"running","messages_ready":0,"message_bytes_ready":0,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":55556,"consumers":0,"consumer_utilisation":"","head_message_timestamp":""}`,
		`,{"name":"10_128_4_242:5672.sage-xds-service.wfs_0617_jl.LATEST.perf","state":"running","messages_ready":0,"message_bytes_ready":0,"messages_unacknowledged":0,"message_bytes_unacknowledged":0,"memory":55756,"consumers":1,"consumer_utilisation":1.0,"head_message_timestamp":1630920836}`,
		`]`,
		`{"command_executed":"rabbitmqctl list_queues --formatter json name state messages_ready","command_runtime":0.5655179}`,
	}
	currLine := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			e.outputCh <- output[currLine]
			currLine++
			if currLine == len(output) {
				return nil
			}
		}
	}
}

type TestExecutorJSONFactory struct {}

func NewTestExecutorJSONFactory() IExecutorFactory {
	return &TestExecutorJSONFactory{}
}

func (f *TestExecutorJSONFactory) NewExecutor(command string, arguments []string, outputBuffer int) IExecutor {
	return &TestJSONExecutor{
		outputCh: make(chan string, 100),
		endExecutionCh: make(chan struct{}, 1),
	}
}

//============== TEST ================ //
func TestJsonCollectOk(t *testing.T) {
	console := NewCmdCollector(NewQueueJSONParser(&TrueFilterConfig{}), NewTestExecutorJSONFactory(), 1000000, 1000000)
	results, err := console.Collect()

	assert.Nil(t, err)
	assert.Equal(t, 6, len(results))

	value, err := results[0].GetMetricValue("consumer_utilisation")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, value)

	value, err = results[4].GetMetricValue("head_message_timestamp")
	assert.Nil(t, err)
	assert.Equal(t, float64(1630920836), value)

	value, err = results[5].GetMetricValue("command_runtime")
	assert.Nil(t, err)
	assert.Equal(t, 0.5655179, value)
	labels, err := results[5].GetLabels("command_runtime")
	assert.Nil(t, err)
	assert.Equal(t, `rabbitmqctl list_queues --formatter json name state messages_ready`, labels["command_executed"])
}