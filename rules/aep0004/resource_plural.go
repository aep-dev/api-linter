// Copyright 2023 Google LLC
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
	"fmt"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var resourcePlural = &lint.MessageRule{
	Name:     lint.NewRuleName(4, "resource-plural"),
	OnlyIf:   hasResourceAnnotation,
	RuleType: lint.NewRuleType(lint.MustRule),
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		r := utils.GetResource(m)
		l := locations.MessageResource(m)
		p := r.GetPlural()
		pLower := utils.ToKebabCase(p)
		if p == "" {
			return []lint.Problem{{
				Message:    "Resources should declare plural.",
				Descriptor: m,
				Location:   l,
			}}
		}
		if pLower != p {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Resource plural should be lowerCamelCase: %q", pLower,
				),
				Descriptor: m,
				Location:   l,
			}}
		}
		return nil
	},
}
