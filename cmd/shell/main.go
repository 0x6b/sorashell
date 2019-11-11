package main

import (
	"fmt"
	gp "github.com/c-bata/go-prompt"
	"github.com/soracom/soracom-shell"
)

func main() {
	executor := shell.NewSoracomExecutor("/bin/sh")
	completer := shell.NewSoracomCompleter("/soracom-api.en.yaml")

	fmt.Print(` _  _  _      _ _       _     _    
(_ / \|_) /\ / / \|\/| (_ |_||_ | | 
__)\_/| \/--\\_\_/|  | __)| ||_ |_|_
         Type exit or Ctrl-D to exit
`)
	gp.New(
		executor.Execute,
		completer.Complete,
		gp.OptionTitle("SORACOM Shell"),
		gp.OptionPrefix("SORACOM> "),
		gp.OptionMaxSuggestion(10),
		gp.OptionPrefixTextColor(gp.Cyan),
		gp.OptionPreviewSuggestionTextColor(gp.Blue),
		gp.OptionSelectedSuggestionBGColor(gp.LightGray),
		gp.OptionSuggestionBGColor(gp.DarkGray),
	).Run()
}
