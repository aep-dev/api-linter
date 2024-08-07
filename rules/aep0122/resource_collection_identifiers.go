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

package aep0122

import (
	"regexp"
	"strings"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

var firstCharRegexp = regexp.MustCompile(`^[a-z]`)

var resourceCollectionIdentifiers = &lint.MessageRule{
	Name: lint.NewRuleName(122, "resource-collection-identifiers"),
	OnlyIf: func(m *desc.MessageDescriptor) bool {
		return utils.GetResource(m) != nil
	},
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		var problems []lint.Problem
		resource := utils.GetResource(m)
		for _, p := range resource.GetPattern() {
			if !firstCharRegexp.MatchString(p) {
				return append(problems, lint.Problem{
					Message:    "Resource patterns must start with a lowercase letter.",
					Descriptor: m,
					Location:   locations.MessageResource(m),
				})
			}

			segs := strings.Split(p, "/")
			for _, seg := range segs {
				if HasUpper(seg) || strings.Contains(seg, "_") {
					problems = append(problems, lint.Problem{
						Message:    "Resource patterns must use kebab-case for collection identifiers.",
						Descriptor: m,
						Location:   locations.MessageResource(m),
					})
				}
			}
		}

		return problems
	},
}
