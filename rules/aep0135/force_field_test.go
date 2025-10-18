// Copyright 2023 Google LLC
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

package aep0135

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestForceField(t *testing.T) {
	tests := []struct {
		name      string
		Reference string
		BoolField string
		problems  testutils.Problems
	}{
		{"ValidWithChildren", "library.googleapis.com/Publisher", "force", nil},
		{"ValidWithoutChildren", "library.googleapis.com/Book", "other", nil},
		{"InvalidMissingForce", "library.googleapis.com/Publisher", "other", testutils.Problems{{Message: "bool force"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
  import "aep/api/field_info.proto";

				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
					};

					string path = 1;
				}

				message Publisher {
					option (google.api.resource) = {
						type: "library.googleapis.com/Publisher"
						pattern: "publishers/{publisher}"
					};

					string path = 1;
				}

				message DeleteResourceRequest {
					string path = 1 [(aep.api.field_info).resource_reference = "{{.Reference}}"];

					bool {{.BoolField}} = 2;
				}
			`, test)
			msg := f.FindMessage("DeleteResourceRequest")
			problems := forceField.Lint(f)
			if diff := test.problems.SetDescriptor(msg).Diff(problems); diff != "" {
				t.Errorf("Problems did not match: %v", diff)
			}
		})
	}
}
