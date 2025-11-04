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

		// Check resource reference matches the created resource type.
		// In AEP format, use resource_reference_child_type to reference the created resource.
		// For backwards compatibility, resource_reference (type) is also supported.
		childTypes := ref.GetChildType()
		types := ref.GetType()

		if len(childTypes) > 0 {
			// AEP format with resource_reference_child_type
			if resource.GetType() != childTypes[0] {
				return []lint.Problem{{
					Message:    "Create should use `resource_reference_child_type` to reference the created resource.",
					Descriptor: parent,
					Location:   locations.FieldResourceReference(parent),
				}}
			}
		} else if len(types) > 0 {
			// AEP format with resource_reference only (backwards compatibility)
			if resource.GetType() != types[0] {
				return []lint.Problem{{
					Message:    "Create should use `resource_reference_child_type` to reference the created resource.",
					Descriptor: parent,
					Location:   locations.FieldResourceReference(parent),
				}}
			}
		}

		return nil
	},
}
