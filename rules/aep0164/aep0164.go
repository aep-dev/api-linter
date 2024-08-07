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

// Package aep0164 contains rules defined in https://aep.dev/164.
package aep0164

import (
	"regexp"

	"github.com/aep-dev/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules adds all of the AEP-164 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		164,
		httpBody,
		httpMethod,
		httpURISuffix,
		requestMessageName,
		requestNameBehavior,
		requestNameField,
		requestNameReference,
		requestUnknownFields,
		resourceExpireTimeField,
		responseLRO,
		responseMessageName,
	)
}

var (
	undeleteMethodRegexp     = regexp.MustCompile("^Undelete(?:[A-Z]|$)")
	undeleteReqMessageRegexp = regexp.MustCompile("^Undelete[A-Za-z0-9]*Request$")
	undeleteURINameRegexp    = regexp.MustCompile(`{name=[a-zA-Z/*]+}:undelete$`)
)

// Returns true if this is a AEP-164 Undelete method, false otherwise.
func isUndeleteMethod(m *desc.MethodDescriptor) bool {
	return undeleteMethodRegexp.MatchString(m.GetName())
}

// Returns true if this is an AEP-164 Undelete request message, false otherwise.
func isUndeleteRequestMessage(m *desc.MessageDescriptor) bool {
	return undeleteReqMessageRegexp.MatchString(m.GetName())
}
