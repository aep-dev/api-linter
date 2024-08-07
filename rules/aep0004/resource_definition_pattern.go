// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aep0004

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourceDefinitionPatterns = &lint.FileRule{
	Name:     lint.NewRuleName(4, "resource-definition-pattern"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf:   hasResourceDefinitionAnnotation,
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		var problems []lint.Problem
		resources := utils.GetResourceDefinitions(f)

		for ndx, resource := range resources {
			loc := locations.FileResourceDefinition(f, ndx)
			probs := lintResourcePattern(resource, f, loc)
			problems = append(problems, probs...)
		}
		return problems
	},
}
