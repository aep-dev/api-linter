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

package aep0133

import (
	"fmt"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var requestIDField = &lint.MessageRule{
	Name:     lint.NewRuleName(133, "request-id-field"),
	OnlyIf:   utils.IsCreateRequestMessage,
	RuleType: lint.NewRuleType(lint.ShouldRule),
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		idField := "id"
		if field := m.FindFieldByName(idField); field == nil || utils.GetTypeName(field) != "string" || field.IsRepeated() {
			return []lint.Problem{{
				Message:    fmt.Sprintf("create methods should contain a singular `string %s` field.", idField),
				Descriptor: m,
			}}
		}
		return nil
	},
}
