package shell

import "github.com/soracom/soracom-cli/generators/lib"

// SoracomCompleter returns suggestions for given go-prompt document.
type SoracomCompleter struct {
	// SORACOM CLI API definitions
	apiDef *lib.APIDefinitions
}

type flag struct {
	name  string
	value string
}

type param struct {
	name        string
	required    bool
	description string
	paramType   string
	enum        []string
}
