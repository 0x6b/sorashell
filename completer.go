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

	if line == "" {
		return []gp.Suggest{}
	}

	if endsWithPipeOrRedirect(line) {
		return []gp.Suggest{}
	}

	// return from hard corded Commands as atm don't have a way to find top-level commands from API definition
	if isFirstCommand(line) {
		s := gp.FilterFuzzy(Commands, line, true)
		sort.Slice(s, func(i, j int) bool {
			return s[i].Text < s[j].Text
		})
		return s
	}

	commands, flags := splitToCommandsAndFlags(line)
	methods, found := s.searchMethods(commands)

	if len(flags) == 0 { // command completion
		if !found {
			return []gp.Suggest{}
		}
		if len(methods) == 1 {
			return suggestion(methods[0], commands)
		}
		return suggestions(methods, commands)
	}

	// flags completion
	if len(methods) != 1 { // if we don't have specific method we can't provide flags suggestion
		return []gp.Suggest{{
			Text:        "Error",
			Description: "cannot find matching command",
		}}
	}
	if params, found := s.searchParams(commands, flags); found {
		return paramSuggestions(params, d.GetWordBeforeCursorWithSpace())
	}

	return []gp.Suggest{}
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

	if len(found) == 0 {
		return found, false
	}
	return found, true
}

// search parameters for cli definition
func (s *SoracomCompleter) searchParams(commands, flags string) ([]param, bool) {
	found := make([]param, 0)
	parsedFlags := parseFlags(flags)

	for _, method := range s.apiDef.Methods {
		if method.CLI == nil || len(method.CLI) == 0 {
			continue
		}

		for _, cli := range method.CLI {
			if strings.Compare(cli, strings.TrimSpace(commands)) == 0 {
				for _, p := range method.Parameters {
					if !contains(parsedFlags, p.Name) {
						found = append(found, param{
							name:        strings.ReplaceAll(p.Name, "_", "-"),
							required:    p.Required,
							description: p.Description,
							paramType:   p.Type,
							enum:        p.Enum,
						})
					}
				}
			}
		}
	}

	sort.Slice(found, func(i, j int) bool {
		return found[i].name < found[j].name
	})

	if len(found) == 0 {
		return found, false
	}
	return found, true
}

// return one command suggestion.
func suggestion(found sl.APIMethod, commands string) []gp.Suggest {
	cli := pickCliDefForPrefix(found.CLI, commands)
	n := strings.Count(commands, " ")

	// return only text after current commands as suggestion e.g.
	// - input:            "users password d"
	// - match result:     "users password delete"
	// - number of spaces: 2
	// - returns:          "delete"
	return []gp.Suggest{
		{
			Text:        strings.Join(strings.Split(cli, " ")[n:], " "),
			Description: found.Summary,
		},
	}
}

// return command suggestions.
func suggestions(methods []sl.APIMethod, commands string) []gp.Suggest {
	tmp := make(map[string]bool)
	suggestions := make([]gp.Suggest, 0)
	n := strings.Count(commands, " ")

	for _, apiMethod := range methods {
		cli := strings.Split(pickCliDefForPrefix(apiMethod.CLI, commands), " ")[n]
		if !tmp[cli] {
			tmp[cli] = true
			suggestions = append(suggestions, gp.Suggest{
				Text:        cli,
				Description: apiMethod.Summary,
			})
		}
	}

	return suggestions
}

// return flag (name or value) suggestions.
func paramSuggestions(params []param, param string) []gp.Suggest {
	if strings.HasPrefix(param, "--") { // name suggestion
		r := make([]gp.Suggest, 0)
		for _, p := range params {
			required := ""
			if p.required {
				required = "(required) "
			}

			r = append(r, gp.Suggest{
				Text:        "--" + strings.ReplaceAll(p.name, "_", "-"),
				Description: required + p.description,
			})
		}
		return gp.FilterFuzzy(r, param, true)
	} else { // value suggestion
		// get previous word
		// if specific name is found, do more intelligent completion
		// else get type for previous word
		// if type is enum, provide possible values
		// else do nothing
	}
	return []gp.Suggest{}
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
	return len(strings.Split(s, " ")) <= 1
}
