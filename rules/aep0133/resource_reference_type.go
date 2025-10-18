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

package aep0133

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// Create methods should reference the target resource via `child_type` or the
// parent directly via `type`.
var resourceReferenceType = &lint.MethodRule{
	Name:     lint.NewRuleName(133, "resource-reference-type"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		ot := utils.GetResponseType(m)
		// Unresolvable response_type for an Operation results in nil here.
		resource := utils.GetResource(ot)
		p := m.GetInputType().FindFieldByName("parent")
		return utils.IsCreateMethod(m) && p != nil && utils.GetResourceReference(p) != nil && resource != nil
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// Return type of the RPC.
		ot := utils.GetResponseType(m)
		resource := utils.GetResource(ot)
		parent := m.GetInputType().FindFieldByName("parent")
		ref := utils.GetResourceReference(parent)

		// In AEP format, resource_reference is just a string. When used in Create methods,
		// it should match the created resource type. The old Google API format distinguishes
		// between `type` and `child_type`, but AEP format just uses the string value.
		// If child_type is set (Google API format), check it. Otherwise, check the type field
		// and treat it as an implicit child_type reference.
		if ref.GetChildType() != "" {
			// Google API format with explicit child_type
			if resource.GetType() != ref.GetChildType() {
				return []lint.Problem{{
					Message:    "Create should use a `child_type` reference to the created resource.",
					Descriptor: parent,
					Location:   locations.FieldResourceReference(parent),
				}}
			}
		} else if ref.GetType() != "" {
			// AEP format or Google API format with only type set
			// In AEP format, this should match the created resource type
			if resource.GetType() != ref.GetType() {
				return []lint.Problem{{
					Message:    "Create should use a `child_type` reference to the created resource.",
					Descriptor: parent,
					Location:   locations.FieldResourceReference(parent),
				}}
			}
		}

		return nil
	},
}
