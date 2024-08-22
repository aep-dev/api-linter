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

package aep0158

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

var responseRepeatedFirstField = &lint.MessageRule{
	Name:     lint.NewRuleName(158, "response-repeated-field"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return isPaginatedResponseMessage(m) && len(m.GetFields()) > 0
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		for _, f := range m.GetFields() {
			if(f.IsRepeated()) {
				return nil;
			}
		}

		return []lint.Problem{{
			Message:    "There does not exist a repeated field for pagination results.",
			Descriptor: m,
		}}
	},
}
