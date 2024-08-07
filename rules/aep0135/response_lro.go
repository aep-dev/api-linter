// Copyright 2020 Google LLC
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

package aep0135

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var responseLRO = &lint.MethodRule{
	Name:     lint.NewRuleName(135, "response-lro"),
	RuleType: lint.NewRuleType(lint.ShouldRule),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return utils.IsDeleteMethod(m) && utils.IsDeclarativeFriendlyMethod(m)
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		if !utils.IsOperation(m.GetOutputType()) {
			return []lint.Problem{{
				Message:    "Declarative-friendly delete methods should use an LRO.",
				Descriptor: m,
				Location:   locations.MethodResponseType(m),
				Suggestion: "google.longrunning.Operation",
			}}
		}
		return nil
	},
}
