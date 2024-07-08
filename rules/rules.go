// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package rules contains implementations of rules that apply to Google APIs.
//
// Rules are sorted into subpackages by the AEP (https://aep.dev/) that
// mandates the rule. Every rule represented in code here must be represented
// in English in a corresponding AEP. Conversely, anything mandated in an AEP
// should have a rule here if it is feasible to enforce in code (sometimes it
// is infeasible, however).
//
// A rule is technically anything with a `GetName()`, `GetURI()“, and
// `Lint(*desc.FileDescriptorProto) []lint.Problem` method, but most rule
// authors will want to use the rule structs provided in the lint package
// (`&lint.MessageRule`, `&lint.FieldRule`, and so on). These run against
// each applicable descriptor in the file (`MessageRule` against every message,
// for example). They also have an `OnlyIf` property that can be used to run
// against a subset of descriptors.
//
// A simple rule therefore looks like this:
//
//	var myRule = &lint.MessageRule{
//	  Name: lint.NewRuleName(1234, "my-rule"),
//	  LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
//	    if isBad(m) {
//	      return []lint.Problem{{
//	        Message: "This message is bad.",
//	        Descriptor: m,
//	      }}
//	    }
//	    return nil
//	  },
//	}
//
// Once a rule is written, it needs to be registered. This involves adding
// the rule to the `AddRules` method for the appropriate AEP package.
// If this is the first rule for a new AEP, then the `rules.go` init() function
// must also be updated to run the `AddRules` function for the new package.
package rules

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/aep0004"
	"github.com/aep-dev/api-linter/rules/aep0121"
	"github.com/aep-dev/api-linter/rules/aep0122"
	"github.com/aep-dev/api-linter/rules/aep0126"
	"github.com/aep-dev/api-linter/rules/aep0127"
	"github.com/aep-dev/api-linter/rules/aep0131"
	"github.com/aep-dev/api-linter/rules/aep0132"
	"github.com/aep-dev/api-linter/rules/aep0133"
	"github.com/aep-dev/api-linter/rules/aep0134"
	"github.com/aep-dev/api-linter/rules/aep0135"
	"github.com/aep-dev/api-linter/rules/aep0136"
	"github.com/aep-dev/api-linter/rules/aep0141"
	"github.com/aep-dev/api-linter/rules/aep0142"
	"github.com/aep-dev/api-linter/rules/aep0144"
	"github.com/aep-dev/api-linter/rules/aep0148"
	"github.com/aep-dev/api-linter/rules/aep0151"
	"github.com/aep-dev/api-linter/rules/aep0155"
	"github.com/aep-dev/api-linter/rules/aep0156"
	"github.com/aep-dev/api-linter/rules/aep0157"
	"github.com/aep-dev/api-linter/rules/aep0158"
	"github.com/aep-dev/api-linter/rules/aep0159"
	"github.com/aep-dev/api-linter/rules/aep0164"
	"github.com/aep-dev/api-linter/rules/aep0191"
	"github.com/aep-dev/api-linter/rules/aep0192"
	"github.com/aep-dev/api-linter/rules/aep0216"
)

type addRulesFuncType func(lint.RuleRegistry) error

var aepAddRulesFuncs = []addRulesFuncType{
	aep0121.AddRules,
	aep0122.AddRules,
	aep0004.AddRules,
	aep0126.AddRules,
	aep0127.AddRules,
	aep0131.AddRules,
	aep0132.AddRules,
	aep0133.AddRules,
	aep0134.AddRules,
	aep0135.AddRules,
	aep0136.AddRules,
	aep0141.AddRules,
	aep0142.AddRules,
	aep0144.AddRules,
	aep0148.AddRules,
	aep0151.AddRules,
	aep0155.AddRules,
	aep0156.AddRules,
	aep0157.AddRules,
	aep0158.AddRules,
	aep0159.AddRules,
	aep0164.AddRules,
	aep0191.AddRules,
	aep0192.AddRules,
	aep0216.AddRules,
}

// Add all rules to the given registry.
func Add(r lint.RuleRegistry) error {
	return addAEPRules(r, aepAddRulesFuncs)
}

func addAEPRules(r lint.RuleRegistry, addRulesFuncs []addRulesFuncType) error {
	for _, addRules := range addRulesFuncs {
		if err := addRules(r); err != nil {
			return err
		}
	}
	return nil
}
