// Copyright 2019 Google LLC
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

package aep0134

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Update methods should have a proper HTTP pattern.
var httpNameField = &lint.MethodRule{
	Name:     lint.NewRuleName(134, "http-uri-path"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf:   utils.IsUpdateMethod,
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		return utils.LintHTTPURIHasVariable(m, "path")
	},
}
