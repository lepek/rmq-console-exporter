package collectors

import (
	"github.com/oriser/regroup"
)

type ConsoleQueueParser struct {
	Cmd			string
	Arguments	[]string
	Parser		*regroup.ReGroup
}

func NewConsoleQueueParser() *ConsoleQueueParser {
	return &ConsoleQueueParser{
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

func (p *ConsoleQueueParser) GetCmd() string {
	return p.Cmd
}

func (p *ConsoleQueueParser) GetArguments() []string {
	return p.Arguments
}

func (p *ConsoleQueueParser) GetParser() *regroup.ReGroup {
	return p.Parser
}

func (p *ConsoleQueueParser) GetNewContainer() IMetrics {
	return &RMQQueueMetrics{}
}
