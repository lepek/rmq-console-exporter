package collectors

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

type CmdCollector struct {
	Parser           ICmdParser
	TimeoutMs        int
	CmdQueueExecutor *Executor
}

func NewCmdCollector(parser ICmdParser, timeoutMs int) *CmdCollector {
	return &CmdCollector{
		Parser: parser,
		TimeoutMs: timeoutMs,
		CmdQueueExecutor: NewExecutor(parser.GetCmd(), parser.GetArguments(), 100000),
	}
}

func (c *CmdCollector) Collect() ([]IMetrics, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeoutMs) * time.Millisecond)
	g, ctxError := errgroup.WithContext(ctxTimeout)
	defer cancel()

	var metrics []IMetrics

	g.Go(func() error {
		for {
			select {
			case line:= <- c.CmdQueueExecutor.OutputCh:
				metric, err := c.Parser.Parse(line)
				if err != nil { return err }
				if metric != nil { metrics = append(metrics, metric) }
			case <-c.CmdQueueExecutor.EndExecutionCh:
				return nil
			case <-ctxTimeout.Done():
				return fmt.Errorf("executor timeout while running [%v %v]", c.CmdQueueExecutor.Command, c.CmdQueueExecutor.Arguments)
			case <-ctxError.Done():
				return ctxError.Err()
			}
		}
	})

	g.Go(func() error {
		if err := c.CmdQueueExecutor.Execute(); err != nil { return err }
		return nil
	})

	return metrics, g.Wait()
}
