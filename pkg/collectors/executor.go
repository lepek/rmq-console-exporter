package collectors

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

type Executor struct {
	Command			string
	Arguments		[]string
	OutputCh		chan string
	OutputBuffer	int
	EndExecutionCh	chan struct{}
}

func NewExecutor(command string, arguments []string, outputBuffer int) *Executor {
	return &Executor{
		Command: command,
		Arguments: arguments,
		OutputCh: make(chan string, outputBuffer),
		OutputBuffer: outputBuffer,
		EndExecutionCh: make(chan struct{}),
	}
}

func (e *Executor) Execute() error {
	cmd := exec.Command(e.Command, e.Arguments...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("%v - Command Execute: cmd.StdoutPipe() failed: %s\n", err, string(stderr.Bytes()))
	}

	scanner := bufio.NewScanner(output)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			e.OutputCh <- scanner.Text()
		}
	}()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%v - Command Execute: cmd.Start(): %v", err, string(stderr.Bytes()))
	}

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%v - Command Execute: cmd.Wait(): %v", err, string(stderr.Bytes()))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%v - Command Execute: scan process output error: %s", err, string(stderr.Bytes()))
	}

	e.EndExecutionCh <-struct{}{}
	return nil
}
