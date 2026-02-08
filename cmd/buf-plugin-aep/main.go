package main

import (
	"fmt"
	"log"
	"strings"

	"buf.build/go/bufplugin/check"
	"github.com/aep-dev/api-linter/lint"
	lint_v2 "github.com/aep-dev/api-linter/lint/v2"
	"github.com/aep-dev/api-linter/rules"
	rules_v2 "github.com/aep-dev/api-linter/rules/v2"
)

const (
	aepCategoryID     = "AEP"
	aepCoreCategoryID = "AEP_CORE"
)

func main() {
	spec, err := newSpec()
	if err != nil {
		log.Fatalln(err)
	}
	// AEP rules cannot be run in parallel as there is thread-unsafe code in
	// this repository that causes concurrent read and write access to a map.
	check.Main(spec, check.MainWithParallelism(1))
}

func newSpec() (*check.Spec, error) {
	ruleRegistryV1 := lint.NewRuleRegistry()
	if err := rules.Add(ruleRegistryV1); err != nil {
		return nil, err
	}
	ruleRegistryV2 := lint_v2.NewRuleRegistry()
	if err := rules_v2.Add(ruleRegistryV2); err != nil {
		return nil, err
	}

	ruleSpecs := make([]*check.RuleSpec, 0, len(ruleRegistryV1)+len(ruleRegistryV2))
	
	// Add V1 Rules
	for _, protoRule := range ruleRegistryV1 {
		ruleSpec, err := newRuleSpecV1(protoRule)
		if err != nil {
			return nil, err
		}
		ruleSpecs = append(ruleSpecs, ruleSpec)
	}

	// Add V2 Rules
	for _, protoRule := range ruleRegistryV2 {
		ruleSpec, err := newRuleSpecV2(protoRule)
		if err != nil {
			return nil, err
		}
		ruleSpecs = append(ruleSpecs, ruleSpec)
	}

	return &check.Spec{
		Rules: ruleSpecs,
		Categories: []*check.CategorySpec{
			{
				ID:      aepCategoryID,
				Purpose: "Checks all API Enhancement proposals as specified at https://aep.dev.",
			},
			{
				ID:      aepCoreCategoryID,
				Purpose: "Checks all core API Enhancement proposals as specified at https://aep.dev.",
			},
		},
		Before: before,
	}, nil
}

func newRuleSpecV1(protoRule lint.ProtoRule) (*check.RuleSpec, error) {
	ruleName := protoRule.GetName()
	if !ruleName.IsValid() {
		return nil, fmt.Errorf("lint.RuleName is invalid: %q", ruleName)
	}
	return createRuleSpec(string(ruleName), fmt.Sprintf("Checks AEP rule %s.", ruleName), newRuleHandlerV1(protoRule))
}

func newRuleSpecV2(protoRule lint_v2.ProtoRule) (*check.RuleSpec, error) {
	ruleName := protoRule.GetName()
	if !ruleName.IsValid() {
		return nil, fmt.Errorf("lint.RuleName is invalid: %q", ruleName)
	}
	return createRuleSpec(string(ruleName), fmt.Sprintf("Checks AEP rule %s.", ruleName), newRuleHandlerV2(protoRule))
}

func createRuleSpec(ruleName string, purpose string, handler check.RuleHandler) (*check.RuleSpec, error) {
	split := strings.Split(ruleName, "::")
	if len(split) != 3 {
		return nil, fmt.Errorf("unknown lint.RuleName format, expected three parts split by '::' : %q", ruleName)
	}
	categoryIDs := []string{aepCategoryID}
	switch extraCategoryID := split[0]; extraCategoryID {
	case "core":
		categoryIDs = append(categoryIDs, aepCoreCategoryID)
	default:
		return nil, fmt.Errorf("unknown lint.RuleName format: unknown category %q : %q", extraCategoryID, ruleName)
	}

	// The allowed characters for RuleName are a-z, 0-9, -.
	// The separator :: is also allowed.
	// We do a translation of these into valid check.Rule IDs.
	ruleID := "AEP_" + strings.Join(split[1:3], "_")
	ruleID = strings.ReplaceAll(ruleID, "-", "_")
	ruleID = strings.ToUpper(ruleID)

	return &check.RuleSpec{
		ID:          ruleID,
		CategoryIDs: categoryIDs,
		Default:     true,
		Purpose:     purpose,
		Type:        check.RuleTypeLint,
		Handler:     handler,
	}, nil
}
