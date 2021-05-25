package parsers

import (
	"errors"
	"github.com/oriser/regroup"
	"rmq-console-exporter/pkg/collectors"
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
			`^(?P<Name>[[:ascii:]]+)\s+` +
				`(?P<State>[[:alpha:]]+)\s+` +
				`(?P<MessagesReady>[[:digit:]]+)\s+` +
				`(?P<MessageBytesReady>[[:digit:]]+)\s+` +
				`(?P<MessagesUnacknowledged>[[:digit:]]+)\s+` +
				`(?P<MessageBytesUnacknowledged>[[:digit:]]+)\s+` +
				`(?P<Memory>[[:digit:]]+)\s+` +
				`(?P<Consumers>[[:digit:]]+)\s+` +
				`(?P<ConsumerUtilisation>[[:digit:]]+)\s+` +
				`(?P<HeadMessageTimestamp>[[:digit:]]*)`,
		),
	}
}

func (p *QueueParser) GetCmd() string {
	return p.Cmd
}

func (p *QueueParser) GetArguments() []string {
	return p.Arguments
}

func (p *QueueParser) Parse(line string) (collectors.IMetrics, error) {
	matches, err := p.Parser.Groups(line)
	if err != nil {
		var e *regroup.NoMatchFoundError
		if errors.As(err, &e) { return nil, nil }
		return nil, err
	}
	queueMetrics := collectors.NewMetrics()
	queue, state := matches["Name"], matches["State"]
	for name, value := range matches {
		if name == "Name" || name == "State" { continue }
		fValue, err := strconv.ParseFloat(value, 64)
		if err != nil { continue }
		queueMetrics.AddMetric(name, fValue, map[string]string{"queue": queue, "state": state})
	}

	return queueMetrics, nil
}
