package cmd

import (
	"fmt"
	gp "github.com/c-bata/go-prompt"
	shell "github.com/soracom/soracom-shell"
	"github.com/spf13/cobra"
)

var specifiedProfileName string
var specifiedCoverageType string
var providedAPIKey string
var providedAPIToken string

func init() {
	RootCmd.PersistentFlags().StringVar(&specifiedProfileName, "profile", "", "Specify profile name")
	RootCmd.PersistentFlags().StringVar(&specifiedCoverageType, "coverage-type", "", "Specify coverage type, 'g' for Global, 'jp' for Japan")
	RootCmd.PersistentFlags().StringVar(&providedAPIKey, "api-key", "", "Specify API key otherwise soracom-cli performs authentication on behalf of you")
	RootCmd.PersistentFlags().StringVar(&providedAPIToken, "api-token", "", "Specify API token otherwise soracom-cli performs authentication on behalf of you")
}

var RootCmd = &cobra.Command{
	Use:   "soracom-shell",
	Short: "Interactive shell for SORACOM CLI",
	Long:  "Interactive shell for SORACOM CLI",
	Run: func(cmd *cobra.Command, args []string) {
		executor := shell.NewSoracomExecutor("/bin/sh", specifiedProfileName, specifiedCoverageType, providedAPIKey, providedAPIToken)
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
	},
}
