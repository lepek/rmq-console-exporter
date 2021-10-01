package collectors

type ExecutorFactory struct {}

func NewExecutorFactory() *ExecutorFactory {
	return &ExecutorFactory{}
}

func (f *ExecutorFactory) NewExecutor(command string, arguments []string, outputBuffer int) *Executor {
	return &Executor{
		command: command,
		arguments: arguments,
		outputCh: make(chan string, outputBuffer),
		endExecutionCh: make(chan struct{}, 1),
	}
}
