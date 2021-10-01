package collectors

import (
	"context"
	"errors"
	"github.com/prometheus/common/log"
	"golang.org/x/sync/errgroup"
	"time"
)

type ICmdParser interface {
	GetCmd() string
	GetArguments() []string
	Parse(string) (*Metrics, error)
}

type IExecutor interface {
	Output() <-chan string
	Execute(ctx context.Context) error
}

type IExecutorFactory interface {
	NewExecutor(command string, arguments []string, outputBuffer int) IExecutor
}

type NonFatalError struct {
	err error
}

func NewNonFatalError(err error) *NonFatalError {
	return &NonFatalError{err}
}

func (e NonFatalError) Error() string {
	return e.err.Error()
}

type CmdCollector struct {
	Parser           	ICmdParser
	TimeoutMs        	int
	OutputBuffer		int
	ExecutorFactory 	IExecutorFactory
	ActiveExecutor		IExecutor
}

func NewCmdCollector(parser ICmdParser, executorFactory IExecutorFactory, timeoutMs int, outputBuffer int) *CmdCollector {
	return &CmdCollector{
		Parser: parser,
		TimeoutMs: timeoutMs,
		ExecutorFactory: executorFactory,
		OutputBuffer: outputBuffer,
		ActiveExecutor: nil, // Just for documentation, no need to initialize
	}
}

func (c *CmdCollector) Collect() ([]Metrics, error) {
	c.ActiveExecutor = c.ExecutorFactory.NewExecutor(c.Parser.GetCmd(), c.Parser.GetArguments(), c.OutputBuffer)
	defer c.closeActiveExecutor()

	log.Info("Starting collection of metrics from console")
	g, ctxError := errgroup.WithContext(context.Background())
	ctxTimeout, cancel := context.WithTimeout(ctxError, time.Duration(c.TimeoutMs) * time.Millisecond)

	defer cancel()

	var metrics []Metrics

	// Parsing command output
	g.Go(func() error {
		var nonFatalError *NonFatalError
		log.Info("Starting line listener")
		defer func() {
			log.Info("Shutting down line listener")
			cancel()
		}()
		for {
			select {
			case line, ok := <-c.ActiveExecutor.Output():
				if !ok {
					log.Info("Command execution finished")
					return nil
				}
				log.Info(line)
				metric, err := c.Parser.Parse(line)
				if err != nil && !errors.As(err, &nonFatalError) { return err }
				if metric != nil { metrics = append(metrics, *metric) }
			case <-ctxError.Done():
				return ctxError.Err()
			}
		}
	})

	// Executing command
	g.Go(func() error {
		log.Info("Starting command executor")
		defer func() {
			log.Info("Shutting down line command executor")
		}()
		if err := c.ActiveExecutor.Execute(ctxTimeout); err != nil {
			log.Errorf("Error while executing command: %v", err)
			return err
		}
		return nil
	})

	return metrics, g.Wait()
}

func (c *CmdCollector) closeActiveExecutor() {
	c.ActiveExecutor = nil
}