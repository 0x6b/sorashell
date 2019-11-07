package main

import (
	"fmt"
	gp "github.com/c-bata/go-prompt"
	"github.com/soracom/soracom-shell"
)

func main() {
	fmt.Println("Type `exit` or `Ctrl-D` to exit.")
	gp.New(
		shell.NewSoracomExecutor("/bin/sh").Execute,
		shell.NewSoracomCompleter("resources/soracom-api.en.yaml").Complete,
		gp.OptionTitle("SORACOM Shell"),
		gp.OptionPrefix("SORACOM> "),
		gp.OptionMaxSuggestion(10),
		gp.OptionPrefixTextColor(gp.Cyan),
		gp.OptionPreviewSuggestionTextColor(gp.Blue),
		gp.OptionSelectedSuggestionBGColor(gp.LightGray),
		gp.OptionSuggestionBGColor(gp.DarkGray),
	).Run()
}
