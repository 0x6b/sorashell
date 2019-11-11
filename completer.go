package shell

import (
	gp "github.com/c-bata/go-prompt"
	sl "github.com/soracom/soracom-cli/generators/lib"
	_ "github.com/soracom/soracom-shell/statik"
	"log"
	"sort"
	"strings"
)

// NewSoracomCompleter returns a SoracomCompleter which is based on  api definition loaded from given apiDefPath.
func NewSoracomCompleter(apiDefPath string) *SoracomCompleter {
	apiDef, err := loadAPIDef(apiDefPath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	return &SoracomCompleter{apiDef}
}

// Complete returns suggestions for given Document.
func (s *SoracomCompleter) Complete(d gp.Document) []gp.Suggest {
	line := d.CurrentLine()

	// return from hard corded Commands as atm don't have a way to find top-level commands from API definition
	if isFirstCommand(line) {
		s := filterFunc(Commands, line)
		sort.Slice(s, func(i, j int) bool {
			return s[i].Text < s[j].Text
		})
		return s
	}

	if endsWithPipeOrRedirect(line) {
		return []gp.Suggest{}
	}

	return s.findSuggestions(line)
}

func (s *SoracomCompleter) findSuggestions(line string) []gp.Suggest {
	commands, flags := splitToCommandsAndFlags(line)
	if len(flags) == 0 {
		return s.commandSuggestion(commands)
	}
	return s.flagSuggestions(line)
}

// return command suggestions.
func (s *SoracomCompleter) commandSuggestion(commands string) []gp.Suggest {
	methods, found := s.searchMethods(commands)
	if !found {
		return []gp.Suggest{}
	}

	if len(methods) == 1 {
		cli := pickCliDefForPrefix(methods[0].CLI, commands)
		n := strings.Count(commands, " ")
		text := strings.Join(strings.Split(cli, " ")[n:], " ")
		//fmt.Printf("\n==========\n- methods: %+v\n- cli: %s\n- n: %d\n- text: '%s'\n", methods, cli, n, text)

		// return only text after current commands as suggestion e.g.
		// - input:            "users password d"
		// - match result:     "users password delete"
		// - number of spaces: 2
		// - returns:          "delete"
		return []gp.Suggest{
			{
				Text:        text,
				Description: methods[0].Summary,
			},
		}
	}

	tmp := make(map[string]bool)
	suggestions := make([]gp.Suggest, 0)
	n := strings.Count(commands, " ")

	// filter duplicates
	for _, apiMethod := range methods {
		// Sometimes two are two candidates which causes oob access.
		//   - query subscribers
		//   - query subscribers traffic_volume
		cli := strings.Split(pickCliDefForPrefix(apiMethod.CLI, commands), " ")
		if len(cli) > n {
			cli := cli[n]
			if !tmp[cli] {
				tmp[cli] = true
				suggestions = append(suggestions, gp.Suggest{
					Text:        cli,
					Description: apiMethod.Summary,
				})
			}
		}
	}

	return suggestions
}

// return flag (name or value) suggestions.
func (s *SoracomCompleter) flagSuggestions(line string) []gp.Suggest {
	commands, flags := splitToCommandsAndFlags(line)
	methods, found := s.searchMethods(commands)
	if !found || len(methods) != 1 {
		return []gp.Suggest{{
			Text:        "Error",
			Description: "cannot find matching command",
		}}
	}
	method := methods[0]

	params := make([]param, 0) // all parameters for the method
	for _, p := range method.Parameters {
		params = append(params, param{
			name:        strings.ReplaceAll(p.Name, "_", "-"),
			required:    p.Required,
			description: p.Description,
			paramType:   p.Type,
			enum:        p.Enum,
		})
	}

	sort.Slice(params, func(i, j int) bool {
		return params[i].name < params[j].name
	})

	parsedFlags := parseFlags(flags)
	flagsArray := strings.Split(flags, " ")
	lastWord := flagsArray[len(flagsArray)-1]
	isEnteringFlag := true

	if len(flagsArray) > 1 {
		if strings.HasPrefix(flagsArray[len(flagsArray)-2], "--") &&
			(strings.HasSuffix(line, " ") || !strings.HasPrefix(lastWord, "--")) {
			isEnteringFlag = false
		}
	}
	if strings.HasSuffix(line, " ") {
		isEnteringFlag = false
	}
	if len(flagsArray)%2 == 0 && !strings.HasPrefix(lastWord, "--") && strings.HasSuffix(line, " ") {
		isEnteringFlag = true
	}

	var lastFlag string
	for i := len(flagsArray) - 1; i >= 0; i-- {
		if strings.HasPrefix(flagsArray[i], "--") {
			lastFlag = strings.ReplaceAll(flagsArray[i], "--", "")
			break
		}
	}
	//fmt.Printf("\n----------\n- line: '%s'\n- parsedFlags: %+v\n- flagsArray: %+v\n- lastFlag: %s\n- lastWord: %s\n- isEnteringFlag: %t\n",
	//	line, parsedFlags, flagsArray, lastFlag, lastWord, isEnteringFlag)

	// provide flag name suggestion if user is entering flag
	if isEnteringFlag {
		r := make([]gp.Suggest, 0)
		for _, p := range params {
			if !contains(parsedFlags, p.name) {
				required := ""
				if p.required {
					required = "(required) "
				}

				r = append(r, gp.Suggest{
					Text:        "--" + strings.ReplaceAll(p.name, "_", "-"),
					Description: required + p.description,
				})
			}
		}
		return filterFunc(r, lastWord)
	}

	if strings.HasPrefix(flagsArray[len(flagsArray)-1], "--") {
		lastWord = ""
	}

	// value suggestion
	// if last flag's value type is enum, provide possible values
	var suggests []gp.Suggest
	for _, p := range params {
		if p.name == lastFlag {
			if len(p.enum) > 0 {
				for _, e := range p.enum {
					suggests = append(suggests, gp.Suggest{
						Text:        e,
						Description: "",
					})
				}
			}
			if len(suggests) > 0 {
				return filterFunc(suggests, lastWord)
			}
		}
	}

	// if specific name is found, do more intelligent completion
	switch lastFlag {
	case "status-filter":
		return statusFilterFunc(lastWord)
	case "speed-class-filter":
		return speedClassFilterFunc(lastWord)
	}

	return suggests
}

// search API methods which has x-soracom-cli definition starts with given term
func (s *SoracomCompleter) searchMethods(term string) ([]sl.APIMethod, bool) {
	found := make([]sl.APIMethod, 0)

	for _, method := range s.apiDef.Methods {
		if method.CLI == nil || len(method.CLI) == 0 {
			continue
		}

		for _, cli := range method.CLI {
			// if term ends with space, try exact match first to search 'data get' out of 'data get',
			// 'data get-entries', 'data get-entry'.
			if strings.HasSuffix(term, " ") {
				if strings.Compare(cli, strings.TrimSpace(term)) == 0 {
					found = append(found, method)
				} else if strings.HasPrefix(cli, term) {
					found = append(found, method)
				}
			} else if strings.HasPrefix(cli, term) {
				found = append(found, method)
			}
		}
	}
	sort.Slice(found, func(i, j int) bool {
		return pickCliDefForPrefix(found[i].CLI, term) < pickCliDefForPrefix(found[j].CLI, term)
	})

	return found, len(found) > 0
}

// parse flags
func parseFlags(f string) []flag {
	values := strings.Split(f, " ")
	results := make([]flag, 0)
	inFlag := false
	name := ""

	for _, value := range values {
		switch {
		case strings.HasPrefix(value, "--"):
			inFlag = true
			name = strings.TrimPrefix(value, "--")
		case inFlag:
			inFlag = false
			results = append(results, flag{name, value})
		}
	}

	if inFlag { // add incomplete flag
		results = append(results, flag{name, ""})
	}

	return results
}

// returns cli definition which starts with given term as sometimes an API has multiple cli definitions e.g.
//   '/subscribers/{imsi}/data' has 'data get' and 'subscribers get-data'.
func pickCliDefForPrefix(methods []string, term string) string {
	for _, cli := range methods {
		if strings.HasPrefix(cli, strings.TrimSpace(term)) {
			return cli
		}
	}
	return ""
}

// Split string before and after "--" and returns as tuple.
// The shell command will accept the global flags (--api-key, --api-token, --coverage-type, --profile) as command line
// args so we can naively think the line can split to commands and args between '--'
//
//   subscribers get-data --imsi 440103216354544 --from 1573052400
//   <-- commands -------><-- flags ----------------------------->
func splitToCommandsAndFlags(line string) (string, string) {
	if strings.Index(line, "--") > 1 {
		return line[:strings.Index(line, "--")], strings.TrimSpace(line[strings.Index(line, "--"):])
	}
	return line, ""
}

// check if given flags contains given string
func contains(flags []flag, param string) bool {
	for _, f := range flags {
		if f.name == strings.ReplaceAll(param, "_", "-") {
			return true
		}
	}
	return false
}

func endsWithPipeOrRedirect(s string) bool {
	t := strings.TrimSpace(s)
	return strings.HasSuffix(t, "|") || strings.HasSuffix(t, ">")
}

func isFirstCommand(s string) bool {
	return strings.TrimSpace(s) == "" || len(strings.Split(s, " ")) <= 1
}

var filterFunc = func(suggestions []gp.Suggest, word string) []gp.Suggest {
	return gp.FilterFuzzy(suggestions, word, false)
}

var statusFilterFunc = func(word string) []gp.Suggest {
	return filterFunc([]gp.Suggest{
		{Text: "active", Description: ""},
		{Text: "inactive", Description: ""},
		{Text: "instock", Description: ""},
		{Text: "ready", Description: ""},
		{Text: "shipped", Description: ""},
		{Text: "suspended", Description: ""},
		{Text: "terminated", Description: ""},
	}, word)
}

var speedClassFilterFunc = func(word string) []gp.Suggest {
	return filterFunc([]gp.Suggest{
		{Text: "s1.minimum", Description: ""},
		{Text: "s1.slow", Description: ""},
		{Text: "s1.standard", Description: ""},
		{Text: "s1.fast", Description: ""},
		{Text: "s1.4xfast", Description: ""},
	}, word)
}
