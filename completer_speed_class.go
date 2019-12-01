package sorashell

import gp "github.com/c-bata/go-prompt"

var speedClassFilterSuggestions = func(word string) []gp.Suggest {
	return filterFunc([]gp.Suggest{
		{Text: "s1.minimum", Description: ""},
		{Text: "s1.slow", Description: ""},
		{Text: "s1.standard", Description: ""},
		{Text: "s1.fast", Description: ""},
		{Text: "s1.4xfast", Description: ""},
	}, word, gp.FilterFuzzy)
}
