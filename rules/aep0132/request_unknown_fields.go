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
	"bitbucket.org/creachadair/stringset"
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var allowedFields = stringset.New(
	"parent",       // AEP-132
	"max_page_size",    // AEP-158
	"page_token",   // AEP-158
	"skip",         // AEP-158
	"filter",       // AEP-132
	"order_by",     // AEP-132
	"show_deleted", // AEP-135
	"request_id",   // AEP-155
	"read_mask",    // AEP-157
	"view",         // AEP-157
)

// List methods should not have unrecognized fields.
var unknownFields = &lint.FieldRule{
	Name:     lint.NewRuleName(132, "request-unknown-fields"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.IsListRequestMessage(f.GetOwner())
	},
	LintField: func(field *desc.FieldDescriptor) []lint.Problem {
		if !allowedFields.Contains(field.GetName()) {
			return []lint.Problem{{
				Message:    "List RPCs should only contain fields explicitly described in AEPs.",
				Descriptor: field,
			}}
		}

		return nil
	},
}
