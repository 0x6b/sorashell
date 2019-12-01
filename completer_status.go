package sorashell

import gp "github.com/c-bata/go-prompt"

var statusFilterSuggestions = func(word string) []gp.Suggest {
	return filterFunc([]gp.Suggest{
		{Text: "active", Description: ""},
		{Text: "inactive", Description: ""},
		{Text: "instock", Description: ""},
		{Text: "ready", Description: ""},
		{Text: "shipped", Description: ""},
		{Text: "suspended", Description: ""},
		{Text: "terminated", Description: ""},
	}, word, gp.FilterFuzzy)
}
