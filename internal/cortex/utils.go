package cortex

import (
	"github.com/grafana/cortex-tools/pkg/rules"
	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prometheus/prometheus/promql"
	"gopkg.in/yaml.v3"
	"log"
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

	oldExpr, err := promql.ParseExpr(oldRG.Rules[0].Expr.Value)
	if err != nil {
		log.Printf("[ERROR] Error parsing expression:\n\t%v\n", err)
		return false
	}
	oldRG.Rules[0].Expr.Value = oldExpr.String()

	newExpr, err := promql.ParseExpr(newRG.Rules[0].Expr.Value)
	if err != nil {
		log.Printf("[ERROR] Error parsing expression:\n\t%v\n", err)
		return false
	}
	newRG.Rules[0].Expr.Value = newExpr.String()


	err = rules.CompareGroups(oldRG, newRG)
	if err != nil {
		log.Printf("[DEBUG] Diff error:\n\t%s\n", err.Error())
		return false
	}
	return true
}
