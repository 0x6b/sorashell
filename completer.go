package shell

import (
	"github.com/c-bata/go-prompt"
	"github.com/soracom/soracom-cli/generators/lib"
	"log"
	"sort"
	"strings"
)

// NewSoracomCompleter returns a SoracomCompleter which is based on  api definition loaded from given apiDefPath.
func NewSoracomCompleter(apiDefPath string) *SoracomCompleter {
	apiDef, err := lib.LoadAPIDef(apiDefPath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	return &SoracomCompleter{apiDef}
}

// Complete returns suggestions for given Document.
func (s *SoracomCompleter) Complete(d prompt.Document) []prompt.Suggest {
	line := d.CurrentLine()

	if line == "" {
		return []prompt.Suggest{}
	}

	if endsWithPipeOrRedirect(line) {
		return []prompt.Suggest{}
	}

	// return from hard corded Commands as atm don't have a way to find top-level commands from API definition
	if isFirstCommand(line) {
		s := prompt.FilterFuzzy(Commands, line, true)
		sort.Slice(s, func(i, j int) bool {
			return s[i].Text < s[j].Text
		})
		return s
	}

	commands, flags := splitToCommandsAndFlags(line)
	found := s.searchMethods(commands)

	if len(flags) == 0 { // command completion
		switch l := len(found); {
		case l == 1:
			return suggestion(found[0], commands)
		case l > 1:
			return suggestions(found, commands)
		default:
			return []prompt.Suggest{}
		}
	} else { // flags completion
		if len(found) != 1 { // if we don't have specific command we can't provide args suggestion
			return []prompt.Suggest{{
				Text:        "Error",
				Description: "cannot find matching command",
			}}
		}

		//parsedArgs := parseFlags(flags)
		params := s.searchParams(commands)

		r := make([]prompt.Suggest, 0)

		for _, p := range params {
			r = append(r, prompt.Suggest{
				Text:        "--" + strings.ReplaceAll(p.name, "_", "-"),
				Description: p.description,
			})
		}
		return prompt.FilterFuzzy(r, d.GetWordBeforeCursorWithSpace(), true)
	}
}

// search API methods which has x-soracom-cli definition starts with given term
func (s *SoracomCompleter) searchMethods(term string) []lib.APIMethod {
	found := make([]lib.APIMethod, 0)

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

	return found
}

// search parameters for cli definition
func (s *SoracomCompleter) searchParams(term string) []param {
	found := make([]param, 0)

	for _, method := range s.apiDef.Methods {
		if method.CLI == nil || len(method.CLI) == 0 {
			continue
		}

		for _, cli := range method.CLI {
			if strings.Compare(cli, strings.TrimSpace(term)) == 0 {
				for _, p := range method.Parameters {
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

	sort.Slice(found, func(i, j int) bool {
		return found[i].name < found[j].name
	})

	return found
}

// return one command suggestion.
func suggestion(found lib.APIMethod, commands string) []prompt.Suggest {
	cli := pickCliDefForPrefix(found.CLI, commands)
	n := strings.Count(commands, " ")

	// return only text after current commands as suggestion e.g.
	// - input:            "users password d"
	// - match result:     "users password delete"
	// - number of spaces: 2
	// - returns:          "delete"
	return []prompt.Suggest{
		{
			Text:        strings.Join(strings.Split(cli, " ")[n:], " "),
			Description: found.Summary,
		},
	}
}

// return command suggestions.
func suggestions(found []lib.APIMethod, commands string) []prompt.Suggest {
	tmp := make(map[string]bool)
	suggestions := make([]prompt.Suggest, 0)
	n := strings.Count(commands, " ")

	for _, apiMethod := range found {
		cli := strings.Split(pickCliDefForPrefix(apiMethod.CLI, commands), " ")[n]
		if !tmp[cli] {
			tmp[cli] = true
			suggestions = append(suggestions, prompt.Suggest{
				Text:        cli,
				Description: apiMethod.Summary,
			})
		}
	}

	return suggestions
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

func endsWithPipeOrRedirect(s string) bool {
	t := strings.TrimSpace(s)
	return strings.HasSuffix(t, "|") || strings.HasSuffix(t, ">")
}

func isFirstCommand(s string) bool {
	return len(strings.Split(s, " ")) <= 1
}
