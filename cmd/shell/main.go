package main

import (
	"github.com/0x6b/sorashell/cmd"
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
