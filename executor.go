package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SoracomExecutor executes given string with the shell.
type SoracomExecutor struct {
	// shell which executes a command
	shell string
}

// NewSoracomExecutor returns a SoracomExecutor which executes commands with shell.
func NewSoracomExecutor(shell string) *SoracomExecutor {
	return &SoracomExecutor{shell}
}

// Execute executes given string in the shell.
func (e *SoracomExecutor) Execute(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	var cmd *exec.Cmd
	if strings.HasPrefix(s, "!cd") {
		dir := strings.Split(s, " ")[1]
		if err := os.Chdir(dir); err != nil {
			fmt.Println("failed to change directory: " + dir)
		}
		return
	} else if strings.HasPrefix(s, "!") {
		cmd = exec.Command("/bin/sh", "-c", strings.TrimPrefix(s, "!"))
	} else {
		cmd = exec.Command("/bin/sh", "-c", "soracom "+s)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err.Error())
	}
	return
}
