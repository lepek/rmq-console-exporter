package collectors

import (
	"errors"
	"github.com/oriser/regroup"
	"strconv"
	"strings"
)

type QueueParser struct {
	Cmd			string
	Config		IConfig
	Arguments	[]string
	Parser		*regroup.ReGroup
}

func NewQueueParser(config IConfig) *QueueParser {
	return &QueueParser{
		Cmd: "rabbitmqctl",
		Config: config,
		Arguments: []string{
			"list_queues",
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
	strings.TrimSpace(line)
	matches, err := p.Parser.Groups(line)
	if err != nil {
		var e *regroup.NoMatchFoundError
		if errors.As(err, &e) { return nil, NewNonFatalError(err) }
		return nil, err
	}
	queue, state := matches["name"], matches["state"]
	// If it doesn't go through the filters then we ignore the queue metric
	if !p.Config.filterQueue(queue) {
		return nil, nil
	}
	queueMetrics := NewMetrics()
	for name, value := range matches {
		if name == "name" || name == "state" { continue }
		fValue, err := strconv.ParseFloat(value, 64)
		if err != nil { continue }
		queueMetrics.AddMetric(name, fValue, map[string]string{"queue": queue, "state": state})
	}

	return queueMetrics, nil
}
