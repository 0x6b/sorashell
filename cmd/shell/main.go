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
		gp.OptionMaxSuggestion(5),

		gp.OptionSuggestionBGColor(gp.Turquoise),
		gp.OptionSuggestionTextColor(gp.Black),
		gp.OptionDescriptionBGColor(gp.LightGray),
		gp.OptionDescriptionTextColor(gp.Black),

		gp.OptionSelectedSuggestionBGColor(gp.DarkGray),
		gp.OptionSelectedSuggestionTextColor(gp.White),
		gp.OptionSelectedDescriptionBGColor(gp.DarkGray),
		gp.OptionSelectedDescriptionTextColor(gp.White),

		gp.OptionPrefixTextColor(gp.Cyan),
		gp.OptionPreviewSuggestionTextColor(gp.Blue),
	).Run()
}
