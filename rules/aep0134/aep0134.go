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

// Package aep0134 contains rules defined in https://aep.dev/134.
package aep0134

import (
	"strings"

	"github.com/aep-dev/api-linter/lint"
	"github.com/stoewer/go-strcase"
)

// AddRules accepts a register function and registers each of
// this AEP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		134,
		allowMissing,
		httpBody,
		httpMethod,
		httpNameField,
		methodSignature,
		responseMessageName,
		requestMaskField,
		requestMaskRequired,
		requestMessageName,
		requestPathRequired,
		requestRequiredFields,
		requestResourceField,
		requestResourceRequired,
		responseLRO,
		synonyms,
		unknownFields,
	)
}

func extractResource(reqName string) string {
	// Strips "Update" from the beginning and "Request" from the end.
	return reqName[6 : len(reqName)-7]
}

func fieldNameFromResource(resource string) string {
	return strings.ToLower(strcase.SnakeCase(resource))
}
