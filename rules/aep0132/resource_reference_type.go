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

package aep0132

import (
	aepapi "buf.build/gen/go/aep/api/protocolbuffers/go/aep/api"
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// List methods should reference the target resource via `child_type` or the
// parent directly via `type`.
var resourceReferenceType = &lint.MethodRule{
	Name:     lint.NewRuleName(132, "resource-reference-type"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		p := m.GetInputType().FindFieldByName("parent")

		var resource *aepapi.ResourceDescriptor
		resourceField := utils.GetListResourceMessage(m)
		if resourceField != nil {
			resource = utils.GetResource(resourceField)
		}
		return utils.IsListMethod(m) && p != nil && utils.GetResourceReference(p) != nil && resource != nil
	},
	LintMethod: func(m *desc.MethodDescriptor) []lint.Problem {
		// The first repeated message field must be the paginated resource.
		repeated := utils.GetRepeatedMessageFields(m.GetOutputType())
		resource := utils.GetResource(repeated[0].GetMessageType())

		parent := m.GetInputType().FindFieldByName("parent")
		ref := utils.GetResourceReference(parent)

		// Check resource reference matches the child resource type.
		// In AEP format, use resource_reference_child_type to reference the child resource.
		// For backwards compatibility, resource_reference (type) is also supported.
		childTypes := ref.GetChildType()
		types := ref.GetType()

		if len(childTypes) > 0 {
			// AEP format with resource_reference_child_type
			if resource.GetType() != childTypes[0] {
				return []lint.Problem{{
					Message:    "List should use `resource_reference_child_type` to reference the paginated resource.",
					Descriptor: parent,
					Location:   locations.FieldResourceReference(parent),
				}}
			}
		} else if len(types) > 0 {
			// AEP format with resource_reference only (backwards compatibility)
			if resource.GetType() != types[0] {
				return []lint.Problem{{
					Message:    "List should use `resource_reference_child_type` to reference the paginated resource.",
					Descriptor: parent,
					Location:   locations.FieldResourceReference(parent),
				}}
			}
		}

		return nil
	},
}
