package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/soracom/soracom-shell"
)

func main() {
	fmt.Println("Type `exit` or `Ctrl-D` to exit.")
	prompt.New(
		shell.NewSoracomExecutor("/bin/sh").Execute,
		shell.NewSoracomCompleter("resources/soracom-api.en.yaml").Complete,
		prompt.OptionTitle("SORACOM Shell"),
		prompt.OptionPrefix("SORACOM> "),
		prompt.OptionMaxSuggestion(10),
		prompt.OptionPrefixTextColor(prompt.Cyan),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	).Run()
}
