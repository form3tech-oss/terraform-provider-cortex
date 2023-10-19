package cortex

import (
	"fmt"
	"log"

	"github.com/grafana/cortex-tools/pkg/rules"
	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prometheus/prometheus/promql/parser"
	"gopkg.in/yaml.v3"
)

func formatYAML(input string) (string, error) {
	var rg interface{}
	err := yaml.Unmarshal([]byte(input), &rg)
	if err != nil {
		return "", err
	}
	out, err := yaml.Marshal(rg)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func suppressYAMLDiff(_, old, new string, _ *schema.ResourceData) bool {
	olds, err := formatYAML(old)
	if err != nil {
		return false
	}
	news, err := formatYAML(new)
	if err != nil {
		return false
	}
	return olds == news
}

func suppressRuleGroupDiff(_, old, new string, _ *schema.ResourceData) bool {
	log.Println("[DEBUG] DiffSuppressFunc")
	oldRG := rwrulefmt.RuleGroup{}
	err := yaml.Unmarshal([]byte(old), &oldRG)
	if err != nil {
		log.Printf("[DEBUG] Error parsing old:\n\t%s\n", err)
		log.Printf("[DEBUG] Old value\n%s\n", old)
	}
	newRG := rwrulefmt.RuleGroup{}
	err = yaml.Unmarshal([]byte(new), &newRG)
	if err != nil {
		log.Printf("[DEBUG] Error parsing new:\n\t%s\n", err)
		log.Printf("[DEBUG] New value\n%s\n", old)
	}

	newRG, err = normaliseRuleGroupRuleExpressions(newRG)
	if err != nil {
		log.Printf("[ERROR] Error normalising old rule group expression:\n\t%v\n", err)
		return false
	}

	oldRG, err = normaliseRuleGroupRuleExpressions(oldRG)
	if err != nil {
		log.Printf("[ERROR] Error normalising old rule group expression:\n\t%v\n", err)
		return false
	}

	err = rules.CompareGroups(oldRG, newRG)
	if err != nil {
		log.Printf("[DEBUG] Diff error:\n\t%s\n", err.Error())
		return false
	}
	return true
}

func normaliseRuleGroupRuleExpressions(rg rwrulefmt.RuleGroup) (rwrulefmt.RuleGroup, error) {
	for i := 0; i<len(rg.RuleGroup.Rules); i++ {
		expr, err := parser.ParseExpr(rg.RuleGroup.Rules[i].Expr.Value)
		if err != nil {
			return rg, fmt.Errorf("normalise rule group expression: %w", err)
		}
		rg.RuleGroup.Rules[i].Expr.Value = expr.String()
	}
	return rg, nil
}
