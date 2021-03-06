package sorashell

import (
	"github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SoracomCompleterTestSuite struct {
	suite.Suite
	completer *SoracomCompleter
}

func (suite *SoracomCompleterTestSuite) SetupTest() {
	suite.completer = NewSoracomCompleter("/soracom-api.en.yaml", NewSoracomWorker("/bin/sh", "", "", "", ""))
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

func (suite *SoracomCompleterTestSuite) TestComplete() {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: []string{"audit-logs", "bills", "coupons", "credentials", "data", "devices", "event-handlers", "files", "gadgets", "groups", "lagoon", "logs", "lora-devices", "lora-gateways", "lora-network-sets", "operator", "orders", "payer-information", "payment-history", "payment-methods", "payment-statements", "port-mappings", "products", "query", "roles", "sandbox", "shipping-addresses", "sigfox-devices", "stats", "subscribers", "test", "users", "version", "volume-discounts", "vpg"},
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
			expected: []string{"s1.minimum", "s1.slow", "s1.standard", "s1.fast", "s1.4xfast", "t1.standard", "u1.slow", "u1.standard"},
		},
		{
			input:    "subscribers list --tag-value-match-mode ",
			expected: []string{"exact", "prefix"},
		},
	}

	for _, t := range tests {
		d := prompt.Document{Text: t.input}

		complete := suite.completer.Complete(d)
		r := toStringSlice(complete, func(d prompt.Suggest) string {
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
		d := prompt.Document{Text: t}

		complete := suite.completer.Complete(d)
		r := toStringSlice(complete, func(d prompt.Suggest) string {
			return d.Text
		})

		// { Text: "Error", Description: "cannot find matching command" }
		suite.Equal("Error", r[0])
	}
}

func TestSoracomCompleterTestSuite(t *testing.T) {
	suite.Run(t, new(SoracomCompleterTestSuite))
}

func toStringSlice(suggests []prompt.Suggest, f func(s prompt.Suggest) string) []string {
	r := make([]string, len(suggests))
	for i, v := range suggests {
		r[i] = f(v)
	}
	return r
}
