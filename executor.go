package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// NewSoracomExecutor returns a SoracomExecutor which executes commands with shell.
func NewSoracomExecutor(shell, specifiedProfileName, specifiedCoverageType, providedAPIKey, providedAPIToken string) *SoracomExecutor {
	return &SoracomExecutor{
		shell,
		specifiedProfileName,
		specifiedCoverageType,
		providedAPIKey,
		providedAPIToken,
	}
}

// Execute executes given string in the shell.
func (e *SoracomExecutor) Execute(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	if strings.HasPrefix(s, "!cd") {
		dir := strings.Split(s, " ")[1]
		if err := os.Chdir(dir); err != nil {
			fmt.Println("failed to change directory: " + dir)
		}
		return
	}

	var cmd *exec.Cmd
	if strings.HasPrefix(s, "!") {
		cmd = exec.Command("/bin/sh", "-c", strings.TrimPrefix(s, "!"))
	} else {
		command := "soracom "
		if e.specifiedProfileName != "" {
			command = fmt.Sprintf("%s --profile %s ", command, e.specifiedProfileName)
		}
		if e.specifiedCoverageType != "" {
			command = fmt.Sprintf("%s --coverage-type %s ", command, e.specifiedCoverageType)
		}
		if e.providedAPIKey != "" {
			command = fmt.Sprintf("%s --api-key %s ", command, e.providedAPIKey)
		}
		if e.providedAPIToken != "" {
			command = fmt.Sprintf("%s --api-token %s ", command, e.providedAPIToken)
		}
		cmd = exec.Command("/bin/sh", "-c", command+s)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err.Error())
	}
	return
}
