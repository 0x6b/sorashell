package sorashell

import "github.com/c-bata/go-prompt"

var speedClassFilterSuggestions = func(word string) []prompt.Suggest {
	return filterFunc([]prompt.Suggest{
		{Text: "s1.minimum", Description: "plan01s, plan01s - Low Data Volume, plan-D, plan-K"},
		{Text: "s1.slow", Description: "plan01s, plan01s - Low Data Volume, plan-D, plan-K"},
		{Text: "s1.standard", Description: "plan01s, plan01s - Low Data Volume, plan-D, plan-K"},
		{Text: "s1.fast", Description: "plan01s, plan01s - Low Data Volume, plan-D, plan-K"},
		{Text: "s1.4xfast", Description: "plan01s, plan01s - Low Data Volume"},
		{Text: "t1.standard", Description: "plan-KM1"},
		{Text: "u1.slow", Description: "plan-DU"},
		{Text: "u1.standard", Description: "plan-DU"},
	}, word, prompt.FilterFuzzy)
}
