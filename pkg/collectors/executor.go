package collectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-cmd/cmd"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

type Executor struct {
	command			string
	arguments		[]string
	outputCh		chan string
	endExecutionCh	chan struct{}
}

// The Output Channel is closed when the execution finishes.
// TBD: This can be more elegant than just exposing the output channel
func (e *Executor) Output() <-chan string {
	return e.outputCh
}

func (e *Executor) Execute(ctx context.Context) error {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(100_000_000)
	log.Infof("Executor started with ID %d", id)

	defer func() {
		close(e.outputCh)
		log.Infof("Executed finished with ID %d", id)
	}()

	options := cmd.Options{Streaming: true}
	log.Infof("Executing %v %v", e.command, strings.Join(e.arguments, " "))
	cmd := cmd.NewCmdOptions(options, e.command, e.arguments...)
	for {
		select {
		// Command has finished and the final status of the execution is received
		case finalStatus := <-cmd.Start():
			if finalStatus.Error == nil && finalStatus.Exit == 0 && finalStatus.Complete {
				if status, err := e.statusToJson(finalStatus); err == nil {
					e.outputCh <- status
				}
				log.Infof( "Command end streaming sucessfully: %v", finalStatus)
				return nil
			}
			log.Errorf("Command failed. Error: %v - Exit Code: %v", finalStatus.Error, finalStatus.Exit)
			return finalStatus.Error
		// Output of the command in execution
		case outputLine := <-cmd.Stdout:
			e.outputCh <- outputLine
		// Errors on StdErr will cancel of the execution
		case outputErr, ok := <-cmd.Stderr:
			if ok {
				log.Errorf("Command returned an error: %v", outputErr)
				cmd.Stop()
				return errors.New(outputErr)
			}
		// Context pass to the execution is cancel for any reason (timeout or external errors)
		case <-ctx.Done():
			cmd.Stop()
			return fmt.Errorf("executor timeout or parser error while running [%v %v]", e.command, e.arguments)
		}
	}
}

func (e *Executor) statusToJson(status cmd.Status) (string, error) {
	statusMap := map[string]interface{}{
		"command_runtime": status.Runtime,
		"command_executed": fmt.Sprintf("%s %s", e.command, strings.Join(e.arguments, " ")),
	}
	statusByte, err := json.Marshal(statusMap)
	if err != nil { return "", err }
	return string(statusByte), nil
}
