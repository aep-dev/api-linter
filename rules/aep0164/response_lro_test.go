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

package aep0164

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestResponseLRO(t *testing.T) {
	// Note: AEP ResourceDescriptor doesn't have a style field, so declarative-friendly
	// detection is not possible. All tests expect nil since resources won't be
	// detected as declarative-friendly.
	for _, test := range []struct {
		name         string
		ResponseType string
		problems     testutils.Problems
	}{
		{"ValidNotDF", "Book", nil},
		{"ValidLRO", "google.longrunning.Operation", nil},
		{"NoLongerInvalidDFSync", "Book", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "aep/api/resource.proto";
				import "google/longrunning/operations.proto";
				service Library {
					rpc UndeleteBook(UndeleteBookRequest) returns ({{.ResponseType}});
				}
				message Book {
					option (aep.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
					};
				}
				message UndeleteBookRequest {}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(responseLRO.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
