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

package aep0134

import (
	"fmt"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var requestPathRequired = &lint.MessageRule{
	Name:   lint.NewRuleName(134, "request-path-required"),
	OnlyIf: utils.IsUpdateRequestMessage,
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		if m.FindFieldByName("path") == nil {
			return []lint.Problem{{
				Message:    fmt.Sprintf("Method %q has no `path` field", m.GetName()),
				Descriptor: m,
			}}
		}
		return nil
	},
}
