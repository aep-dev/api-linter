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

// Package aep0144 contains rules defined in https://aep.dev/144.
package aep0144

import (
	"regexp"

	"github.com/aep-dev/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AEP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		144,
		httpBody,
		httpMethod,
	)
}

var addRemoveMethodRegexp = regexp.MustCompile("^(?:Add|Remove)(?:[A-Z]|$)")

// Returns true if this is an AEP-144 Add/Remove method, false otherwise.
func isAddRemoveMethod(m *desc.MethodDescriptor) bool {
	return addRemoveMethodRegexp.MatchString(m.GetName())
}
