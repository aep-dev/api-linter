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

package aep0133

import (
	"fmt"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// The create request message should not have unrecognized fields.
var unknownFields = &lint.MessageRule{
	Name:     lint.NewRuleName(133, "request-unknown-fields"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf:   utils.IsCreateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) (problems []lint.Problem) {
		resourceMsgName := getResourceMsgNameFromReq(m)

		// Rule check: Establish that there are no unexpected fields.
		allowedFields := map[string]*builder.FieldType{
			"parent":        nil, // AEP-133
			"request_id":    nil, // AEP-155
			"validate_only": nil, // AEP-163
			"id":            nil,
		}

		for _, field := range m.GetFields() {
			// Skip the check with the field that is the body.
			if t := field.GetMessageType(); t != nil && t.GetName() == resourceMsgName {
				continue
			}
			// Check the remaining fields.
			if _, ok := allowedFields[string(field.GetName())]; !ok {
				problems = append(problems, lint.Problem{
					Message: fmt.Sprintf(
						"Create RPCs must only contain fields explicitly described in AEPs, not %q.",
						field.GetName(),
					),
					Descriptor: field,
				})
			}
		}

		return
	},
}
