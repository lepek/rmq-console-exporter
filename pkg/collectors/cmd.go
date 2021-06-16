package collectors

import (
	"context"
	"fmt"
	"github.com/prometheus/common/log"
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
	log.Info("Starting collection of metrics from console")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeoutMs) * time.Millisecond)
	g, ctxError := errgroup.WithContext(ctxTimeout)
	defer cancel()

	var metrics []IMetrics

	g.Go(func() error {
		log.Info("Starting line listener")
		defer func() {
			log.Info("Shutting down line listener")
			cancel()
		}()
		for {
			select {
			case line:= <- c.CmdQueueExecutor.OutputCh:
				metric, err := c.Parser.Parse(line)
				if err != nil { return err }
				if metric != nil { metrics = append(metrics, metric) }
			case <-c.CmdQueueExecutor.EndExecutionCh:
				return nil
			case <-ctxTimeout.Done():
				metrics = nil
				return fmt.Errorf("executor timeout while running [%v %v]", c.CmdQueueExecutor.Command, c.CmdQueueExecutor.Arguments)
			case <-ctxError.Done():
				metrics = nil
				return ctxError.Err()
			}
		}
	})

	g.Go(func() error {
		log.Info("Starting command executor")
		defer func() {
			log.Info("Shuting down line command executor")
		}()
		if err := c.CmdQueueExecutor.Execute(ctxTimeout); err != nil {
			log.Errorf("Error while executing command: %v", err)
			return err
		}
		return nil
	})

	return metrics, g.Wait()
}
