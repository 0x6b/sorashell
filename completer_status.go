package sorashell

import "github.com/c-bata/go-prompt"

var statusFilterSuggestions = func(word string) []prompt.Suggest {
	return filterFunc([]prompt.Suggest{
		{Text: "active", Description: ""},
		{Text: "inactive", Description: ""},
		{Text: "instock", Description: ""},
		{Text: "ready", Description: ""},
		{Text: "shipped", Description: ""},
		{Text: "suspended", Description: ""},
		{Text: "terminated", Description: ""},
	}, word, prompt.FilterFuzzy)
}
