package shell

import "fmt"

// NewSoracomExecutor returns a SoracomExecutor which executes commands with shell.
func NewSoracomExecutor(worker *SoracomWorker) *SoracomExecutor {
	return &SoracomExecutor{
		worker,
	}
}

// Execute executes given string in the shell.
func (e *SoracomExecutor) Execute(s string) {
	fmt.Printf("%s", e.worker.Execute(s))
}
