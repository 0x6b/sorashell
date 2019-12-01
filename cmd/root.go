package cmd

import (
	"fmt"
	shell "github.com/0x6b/sorashell"
	gp "github.com/c-bata/go-prompt"
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
		worker := shell.NewSoracomWorker("/bin/sh", profileName, coverageType, apiKey, apiToken)
		executor := shell.NewSoracomExecutor(worker)
		completer := shell.NewSoracomCompleter("/soracom-api.en.yaml", worker)

		fmt.Print(` _  _  _      _     _    
(_ / \|_) /\ (_ |_||_ | | 
__)\_/| \/--\__)| ||_ |_|_      Type exit or Ctrl-D to exit
`)
		gp.New(
			executor.Execute,
			completer.Complete,
			gp.OptionTitle("SORASHELL"),
			gp.OptionPrefix("SORASHELL> "),
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
	},
}
