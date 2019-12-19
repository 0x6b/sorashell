package sorashell

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// NewSoracomWorker returns a SoracomExecutor which executes commands with shell.
func NewSoracomWorker(shell, profileName, coverageType, apiKey, apiToken string) *SoracomWorker {
	c := make(chan string, 1)

	sc := &SoracomWorker{
		shell,
		profileName,
		coverageType,
		apiKey,
		apiToken,
		c,
	}
	go sc.run(c)

	return sc
}

// Execute executes given string after simple parse
func (w *SoracomWorker) Execute(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	if strings.HasPrefix(s, "!cd") {
		dir := strings.Split(s, " ")[1]
		if err := os.Chdir(dir); err != nil {
			fmt.Println("failed to change directory: " + dir)
		}
		return ""
	}

	w.command <- s
	return <-w.command
}

func (w *SoracomWorker) run(ch chan string) {
	for {
		select {
		case s := <-ch:
			if s == "quit" || s == "exit" {
				fmt.Println("Bye!")
				os.Exit(0)
			}

			command := ""
			if strings.HasPrefix(s, "!") {
				command = strings.TrimPrefix(s, "!")
			} else {
				command = "soracom "
				if w.profileName != "" {
					command = fmt.Sprintf("%s --profile %s ", command, w.profileName)
				}
				if w.coverageType != "" {
					command = fmt.Sprintf("%s --coverage-type %s ", command, w.coverageType)
				}
				if w.apiKey != "" {
					command = fmt.Sprintf("%s --api-key %s ", command, w.apiKey)
				}
				if w.apiToken != "" {
					command = fmt.Sprintf("%s --api-token %s ", command, w.apiToken)
				}
				command = fmt.Sprintf("%s %s", command, s)
			}

			cmd := exec.Command("/bin/sh", "-c", command)
			stdout, err := cmd.StdoutPipe()

			if err != nil {
				ch <- fmt.Sprintf("Error while setting up %s", command)
			}
			if err := cmd.Start(); err != nil {
				ch <- fmt.Sprintf("Error while starting %s", command)
			}
			result, err := ioutil.ReadAll(stdout)
			if err != nil {
				ch <- fmt.Sprintf("Error while reading %s", command)
			}

			ch <- string(result)
		}
	}
}
