package collectors

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/prometheus/common/log"
	"math/rand"
	"time"
)

type Executor struct {
	command			string
	arguments		[]string
	outputCh		chan string
	endExecutionCh	chan struct{}
}

func (e *Executor) Output() <-chan string {
	return e.outputCh
}

func (e *Executor) Execute(ctx context.Context) error {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(100000)
	log.Infof("Executor started with ID %d", id)

	defer func() {
		close(e.outputCh)
		log.Infof("Executed finished with ID %d", id)
	}()

	options := cmd.Options{Streaming: true}
	cmd := cmd.NewCmdOptions(options, e.command, e.arguments...)
	for {
		select {
		case finalStatus := <-cmd.Start():
			fmt.Println(finalStatus)
			return nil
		case outputLine := <-cmd.Stdout:
			e.outputCh <- outputLine
		case outputErr := <-cmd.Stderr:
			fmt.Println(outputErr)
			cmd.Stop()
			return errors.New(outputErr)
		case <-ctx.Done():
			cmd.Stop()
			return fmt.Errorf("executor timeout or parser error while running [%v %v]", e.command, e.arguments)
		}
	}
}
