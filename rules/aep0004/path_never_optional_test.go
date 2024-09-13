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

package aep0004

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestPathNeverOptional(t *testing.T) {
	for _, test := range []struct {
		name      string
		FieldPath string
		PathField string
		Label     string
		problems  testutils.Problems
	}{
		{"Valid", "path", "", "", testutils.Problems{}},
		{"ValidAlternativePath", "resource", "resource", "", testutils.Problems{}},
		{"InvalidProto3Optional", "path", "", "optional", testutils.Problems{{Message: "never be labeled"}}},
		{"SkipPathFieldDNE", "path", "does_not_exist", "", testutils.Problems{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
						name_field: "{{.PathField}}"
					};

					{{.Label}} string {{.FieldPath}} = 1;
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(pathNeverOptional.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
