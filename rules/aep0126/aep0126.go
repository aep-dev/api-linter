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

// Package aep0126 contains rules defined in https://aep.dev/126.
package aep0126

import "github.com/aep-dev/api-linter/lint"

// AddRules adds all of the AEP-126 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		126,
		unspecified,
	)
}
