package collectors

import (
	"context"
	"errors"
	"fmt"
	"github.com/oriser/regroup"
	"golang.org/x/sync/errgroup"
	"time"
)

type ConsoleSource struct {
	Parser				IConsoleParser
	TimeoutMs			int
	CmdQueueExecutor	*Executor
}

func NewConsoleSource(parser IConsoleParser, timeoutMs int) *ConsoleSource {
	return &ConsoleSource{
		Parser: parser,
		TimeoutMs: timeoutMs,
		CmdQueueExecutor: NewExecutor(parser.GetCmd(), parser.GetArguments(), 100000),
	}
}

func (c *ConsoleSource) Collect() ([]IMetrics, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeoutMs) * time.Millisecond)
	g, ctxError := errgroup.WithContext(ctxTimeout)
	defer cancel()

	var metrics []IMetrics

	g.Go(func() error {
		for {
			select {
			case line:= <- c.CmdQueueExecutor.OutputCh:
				queueMetrics := c.Parser.GetNewContainer()
				if err := c.Parser.GetParser().MatchToTarget(line, queueMetrics); err != nil {
					var e *regroup.NoMatchFoundError
					if errors.As(err, &e) { continue }
					return err
				}
				metrics = append(metrics, queueMetrics)
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
