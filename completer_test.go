package shell

import (
	gp "github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SoracomCompleterTestSuite struct {
	suite.Suite
	completer *SoracomCompleter
}

func (suite *SoracomCompleterTestSuite) SetupTest() {
	suite.completer = NewSoracomCompleter("/soracom-api.en.yaml")
}

func (suite *SoracomCompleterTestSuite) TestSplitToCommandsAndFlags() {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "subscribers get --imsi",
			expected: []string{
				"subscribers get ",
				"--imsi",
			},
		},
		{
			input: "subscribers get --",
			expected: []string{
				"subscribers get ",
				"--",
			},
		},
	}

	for _, t := range tests {
		commands, args := splitToCommandsAndFlags(t.input)
		suite.Equal(t.expected, []string{commands, args})
	}
}

func (suite *SoracomCompleterTestSuite) TestParseFlags() {
	tests := []struct {
		input    string
		expected []flag
	}{
		{
			input: "--imsi 44010",
			expected: []flag{{
				name:  "imsi",
				value: "44010",
			}},
		},
		{
			input: "--imsi 44010 --operator-id OPxxxxxx",
			expected: []flag{
				{
					name:  "imsi",
					value: "44010",
				},
				{
					name:  "operator-id",
					value: "OPxxxxxx",
				},
			},
		},
		{
			input: "--imsi",
			expected: []flag{
				{
					name:  "imsi",
					value: "",
				},
			},
		},
	}

	for _, t := range tests {
		r := parseFlags(t.input)
		suite.Equal(t.expected, r)
	}
}

func (suite *SoracomCompleterTestSuite) TestSearchParams() {
	tests := []struct {
		input    string
		expected []param
	}{
		{
			input: "users password configured",
			expected: []param{
				{
					name:        "operator-id",
					required:    true,
					description: "operator_id",
					paramType:   "string",
					enum:        []string(nil),
				},
				{
					name:        "user-name",
					required:    true,
					description: "user_name",
					paramType:   "string",
					enum:        []string(nil),
				},
			},
		},
		{
			input: "users password configured --operator-id OPXXX --",
			expected: []param{
				{
					name:        "user-name",
					required:    true,
					description: "user_name",
					paramType:   "string",
					enum:        []string(nil),
				},
			},
		},
		{
			input: "groups delete-config",
			expected: []param{
				{
					name:        "group-id",
					required:    true,
					description: "Target group.",
					paramType:   "string",
					enum:        []string(nil),
				},
				{
					name:        "name",
					required:    true,
					description: "Parameter name to be deleted. (This will be part of a URL path, so it needs to be percent-encoded. In JavaScript, specify the name after it has been encoded using encodeURIComponent().)",
					paramType:   "string",
					enum:        []string(nil),
				},
				{
					name:        "namespace",
					required:    true,
					description: "Namespace of target parameters.",
					paramType:   "string",
					enum: []string{
						"SoracomAir",
						"SoracomBeam",
						"SoracomEndorse",
						"SoracomFunk",
						"SoracomFunnel",
						"SoracomHarvest",
						"SoracomHarvestFiles",
						"SoracomKrypton",
						"UnifiedEndpoint",
					},
				},
			},
		},
		{
			input: "groups delete-config --group-id xxx --namespace SoracomAir",
			expected: []param{
				{
					name:        "name",
					required:    true,
					description: "Parameter name to be deleted. (This will be part of a URL path, so it needs to be percent-encoded. In JavaScript, specify the name after it has been encoded using encodeURIComponent().)",
					paramType:   "string",
					enum:        []string(nil),
				},
			},
		},
	}

	for _, t := range tests {
		commands, flags := splitToCommandsAndFlags(t.input)
		methods, _ := suite.completer.searchMethods(commands)
		r, _ := suite.completer.searchParams(methods[0], flags)
		suite.Equal(t.expected, r)
	}
}

func (suite *SoracomCompleterTestSuite) TestComplete() {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "s",
			expected: []string{"audit-logs", "bills", "coupons", "credentials", "devices", "event-handlers", "files", "gadgets", "groups", "logs", "lora-devices", "lora-gateways", "lora-network-sets", "orders", "payment-history", "payment-methods", "payment-statements", "port-mappings", "products", "roles", "sandbox", "shipping-addresses", "sigfox-devices", "stats", "subscribers", "test", "users", "version", "volume-discounts"},
		},
		{
			input:    "subscribers",
			expected: []string{"subscribers"},
		},
		{
			input:    "users d",
			expected: []string{"default-permissions", "delete", "detach-role"},
		},
		{
			input:    "users password",
			expected: []string{"password"},
		},
		{
			input:    "groups delete-config",
			expected: []string{"delete-config"},
		},
		{
			input:    "groups delete-config ",
			expected: []string{""},
		},
		{
			input:    "users password ",
			expected: []string{"configured", "create", "delete", "update"},
		},
		{
			input:    "subscribers list |",
			expected: []string{},
		},
		{
			input:    "subscribers list >",
			expected: []string{},
		},
		{
			input:    "!",
			expected: []string{},
		},
		{
			input:    "subscribers list --speed-class-filter ",
			expected: []string{"s1.minimum", "s1.slow", "s1.standard", "s1.fast", "s1.4xfast"},
		},
		{
			input:    "subscribers list --tag-value-match-mode ",
			expected: []string{"exact", "prefix"},
		},
	}

	for _, t := range tests {
		d := gp.Document{Text: t.input}

		complete := suite.completer.Complete(d)
		r := toStringSlice(complete, func(d gp.Suggest) string {
			return d.Text
		})

		suite.Equal(t.expected, r)
	}
}

func (suite *SoracomCompleterTestSuite) TestIncompleteCommand() {
	tests := []string{
		"subscribers --imsi",
		"subscribers l --imsi",
	}

	for _, t := range tests {
		d := gp.Document{Text: t}

		complete := suite.completer.Complete(d)
		r := toStringSlice(complete, func(d gp.Suggest) string {
			return d.Text
		})

		// { Text: "Error", Description: "cannot find matching command" }
		suite.Equal("Error", r[0])
	}
}

func TestSoracomCompleterTestSuite(t *testing.T) {
	suite.Run(t, new(SoracomCompleterTestSuite))
}

func toStringSlice(suggests []gp.Suggest, f func(s gp.Suggest) string) []string {
	r := make([]string, len(suggests))
	for i, v := range suggests {
		r[i] = f(v)
	}
	return r
}
