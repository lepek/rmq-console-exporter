package collectors

import (
	"errors"
	"github.com/oriser/regroup"
	"strconv"
)

type QueueParser struct {
	Cmd			string
	Arguments	[]string
	Parser		*regroup.ReGroup
}

func NewQueueParser() *QueueParser {
	return &QueueParser{
		Cmd: "rabbitmqctl",
		Arguments: []string{
			"list_queues",
			"--formatter json",
			"name",
			"state",
			"messages_ready",
			"message_bytes_ready",
			"messages_unacknowledged",
			"message_bytes_unacknowledged",
			"memory",
			"consumers",
			"consumer_utilisation",
			"head_message_timestamp",
		},
		Parser: regroup.MustCompile(
			`^(?P<name>[[:ascii:]]+)\t` +
				`(?P<state>[[:alpha:]]+)\t` +
				`(?P<messages_ready>[[:digit:]]+)\t` +
				`(?P<message_bytes_ready>[[:digit:]]+)\t` +
				`(?P<messages_unacknowledged>[[:digit:]]+)\t` +
				`(?P<message_bytes_unacknowledged>[[:digit:]]+)\t` +
				`(?P<memory>[[:digit:]]+)\t` +
				`(?P<consumers>[[:digit:]]+)(?:\t)?` +
				`(?P<consumer_utilisation>[[:digit:]]{1}\.[[:digit:]]*)?(?:\t)?` +
				`(?P<head_message_timestamp>[[:digit:]]*)?`,
		),
	}
}

func (p *QueueParser) GetCmd() string {
	return p.Cmd
}

func (p *QueueParser) GetArguments() []string {
	return p.Arguments
}

func (p *QueueParser) Parse(line string) (*Metrics, error) {
	matches, err := p.Parser.Groups(line)
	if err != nil {
		var e *regroup.NoMatchFoundError
		if errors.As(err, &e) { return nil, NewNonFatalError(err) }
		return nil, err
	}
	queueMetrics := NewMetrics()
	queue, state := matches["name"], matches["state"]
	for name, value := range matches {
		if name == "name" || name == "state" { continue }
		fValue, err := strconv.ParseFloat(value, 64)
		if err != nil { continue }
		queueMetrics.AddMetric(name, fValue, map[string]string{"queue": queue, "state": state})
	}

	return queueMetrics, nil
}
