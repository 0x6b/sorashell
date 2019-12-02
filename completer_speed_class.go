package sorashell

import "github.com/c-bata/go-prompt"

var speedClassFilterSuggestions = func(word string) []prompt.Suggest {
	return filterFunc([]prompt.Suggest{
		{Text: "s1.minimum", Description: ""},
		{Text: "s1.slow", Description: ""},
		{Text: "s1.standard", Description: ""},
		{Text: "s1.fast", Description: ""},
		{Text: "s1.4xfast", Description: ""},
	}, word, prompt.FilterFuzzy)
}
