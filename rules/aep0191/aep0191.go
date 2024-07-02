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

// Package aep0191 contains rules defined in https://aep.dev/191.
package aep0191

import (
	"regexp"

	"github.com/aep-dev/api-linter/lint"
)

// AddRules adds all of the AEP-191 rules to the provided registry.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		191,
		filename,
		fileLayout,
		protoPkg,
		syntax,
	)
}

var (
	versionRegexp        = regexp.MustCompile(`^v[0-9]+(p[0-9]+)?((alpha|beta)[0-9]*)?$`)
	validCharacterRegexp = regexp.MustCompile(`^[a-z0-9\\_\\/]*$`)
)
