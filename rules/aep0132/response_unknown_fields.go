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

// The resource itself is not included here, but also permitted.
// This is covered in code in the rule itself.
var respAllowedFields = stringset.New(
	"results",
	"max_page_size",         // AEP-132
	"next_page_token",       // AEP-158
	"total_size",            // AEP-132
	"unreachable",           // AEP-217
	"unreachable_locations", // Wrong, but a separate AEP-217 rule catches it.
)

var responseUnknownFields = &lint.FieldRule{
	Name:     lint.NewRuleName(132, "response-unknown-fields"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(f *desc.FieldDescriptor) bool {
		return utils.IsListResponseMessage(f.GetOwner())
	},
	LintField: func(f *desc.FieldDescriptor) []lint.Problem {
		if !respAllowedFields.Contains(f.GetName()) {
			return []lint.Problem{{
				Message:    "List responses should only contain fields explicitly described in AEPs.",
				Descriptor: f,
			}}
		}
		return nil
	},
}
