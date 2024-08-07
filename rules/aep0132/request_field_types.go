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

package aep0132

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var knownFields = map[string]func(*desc.FieldDescriptor) []lint.Problem{
	"filter":       utils.LintSingularStringField,
	"order_by":     utils.LintSingularStringField,
	"show_deleted": utils.LintSingularBoolField,
}

// List fields should have the correct type.
var requestFieldTypes = &lint.FieldRule{
	Name:     lint.NewRuleName(132, "request-field-types"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.IsListRequestMessage(f.GetOwner()) && knownFields[f.GetName()] != nil
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		return knownFields[f.GetName()](f)
	},
}
