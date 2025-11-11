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

var pathNeverOptional = &lint.MessageRule{
	Name:     lint.NewRuleName(4, "path-never-optional"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		f := "path"
		if nf := utils.GetResourceNameField(utils.GetResource(m)); nf != "" {
			f = nf
		}
		return utils.IsResource(m) && m.FindFieldByName(f) != nil
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		f := "path"
		if nf := utils.GetResourceNameField(utils.GetResource(m)); nf != "" {
			f = nf
		}
		field := m.FindFieldByName(f)

		if field.IsProto3Optional() {
			return []lint.Problem{{
				Message:    "Resource path fields must never be labeled with proto3_optional",
				Descriptor: field,
				Location:   locations.FieldLabel(field),
				Suggestion: "",
			}}
		}

		return nil
	},
}
