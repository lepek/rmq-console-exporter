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
	console := NewCmdCollector(NewQueueParser(), NewTestExecutorFactory(), 1000000, 1000000)
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