package cortex

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuppressYAMLDiff(t *testing.T) {
	originalYAML := `route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'web.hook'
`
	equivalentYAML := `route:
  group_by:          ['alertname'    ]
  group_wait:      30s
  group_interval: "5m"
  repeat_interval:      1h
  receiver: web.hook
`
	badYAML := `route:
group_wait:      30s
- group_interval: "5m
`
	tests := []struct {
		name             string
		oldValue         string
		newValue         string
		expectedSuppress bool
	}{
		{"original vs original", originalYAML, originalYAML, true},
		{"original vs equivalent", originalYAML, equivalentYAML, true},
		{"original vs empty", originalYAML, "", false},
		{"bad vs bad", badYAML, badYAML, false},
		{"original vs bad", originalYAML, badYAML, false},
		{"boolean vs string equivalent", `some_bool: true`, `some_bool: "true"`, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedSuppress, suppressYAMLDiff("", test.oldValue, test.newValue, nil))
		})
	}
}

func TestSuppressRuleGroupDiff(t *testing.T) {
	tests := []struct {
		name             string
		oldValuePath     string
		newValuePath     string
		expectedSuppress bool
	}{
		{
			"boolean vs string equivalent",
			`./testdata/unquoted_bool.yaml`,
			`./testdata/quoted_bool.yaml`,
			true,
		},
		{
			"RuleGroup with boolean as boolean vs RuleGroup with boolean as string",
			"./testdata/rule_group_with_quoted_boolean.yaml",
			"./testdata/rule_group_with_quoted_boolean.yaml",
			true,
		},
		{
			"RuleGroup with equal rule expressions with different format",
			"./testdata/rule_group_with_single_multi_line_expression.yaml",
			"./testdata/rule_group_with_single_line_expression.yaml",
			true,
		},
		{
			"RuleGroup with multiple equal rule expressions with different formats",
			"./testdata/rule_group_with_multiple_multi_line_expressions.yaml",
			"./testdata/rule_group_with_multiple_single_line_expressions.yaml",
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldValue, err := os.ReadFile(test.oldValuePath)
			require.NoError(t, err)

			newValue, err := os.ReadFile(test.newValuePath)
			require.NoError(t, err)

			require.Equal(t, test.expectedSuppress, suppressRuleGroupDiff("", string(oldValue), string(newValue), nil))
		})
	}
}
