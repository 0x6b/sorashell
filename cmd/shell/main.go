package main

import (
	"github.com/soracom/soracom-shell/cmd"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	err := cmd.RootCmd.Execute()
	if err != nil {
		return -1
	}
	return 0
}
