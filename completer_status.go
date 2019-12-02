package sorashell

import "github.com/c-bata/go-prompt"

func (s *SoracomCompleter) statusFilterSuggestions(word string) []prompt.Suggest {
	return filterFunc([]prompt.Suggest{
		{Text: "active", Description: ""},
		{Text: "inactive", Description: ""},
		{Text: "ready", Description: ""},
		{Text: "suspended", Description: ""},
		{Text: "terminated", Description: ""},
	}, word, prompt.FilterFuzzy)
}
