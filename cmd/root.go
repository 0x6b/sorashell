package cmd

import (
	"fmt"
	"github.com/0x6b/sorashell"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

var profileName string
var coverageType string
var apiKey string
var apiToken string

func init() {
	RootCmd.PersistentFlags().StringVar(&profileName, "profile", "", "Specify profile name")
	RootCmd.PersistentFlags().StringVar(&coverageType, "coverage-type", "", "Specify coverage type, 'g' for Global, 'jp' for Japan")
	RootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Specify API key otherwise soracom-cli performs authentication on behalf of you")
	RootCmd.PersistentFlags().StringVar(&apiToken, "api-token", "", "Specify API token otherwise soracom-cli performs authentication on behalf of you")
}

var RootCmd = &cobra.Command{
	Use:   "sorashell",
	Short: "Interactive shell for SORACOM CLI",
	Long:  "Interactive shell for SORACOM CLI",
	Run: func(cmd *cobra.Command, args []string) {
		worker := sorashell.NewSoracomWorker("/bin/sh", profileName, coverageType, apiKey, apiToken)
		executor := sorashell.NewSoracomExecutor(worker)
		completer := sorashell.NewSoracomCompleter("/soracom-api.en.yaml", worker)

		fmt.Print(` _  _  _      _     _    
(_ / \|_) /\ (_ |_||_ | | 
__)\_/| \/--\__)| ||_ |_|_      Type exit or Ctrl-D to exit
`)
		prompt.New(
			executor.Execute,
			completer.Complete,
			prompt.OptionTitle("SORASHELL"),
			prompt.OptionPrefix("SORASHELL> "),
			prompt.OptionMaxSuggestion(5),

			prompt.OptionSuggestionBGColor(prompt.Turquoise),
			prompt.OptionSuggestionTextColor(prompt.Black),
			prompt.OptionDescriptionBGColor(prompt.LightGray),
			prompt.OptionDescriptionTextColor(prompt.Black),

			prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
			prompt.OptionSelectedSuggestionTextColor(prompt.White),
			prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
			prompt.OptionSelectedDescriptionTextColor(prompt.White),

			prompt.OptionPrefixTextColor(prompt.Cyan),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		).Run()
	},
}
