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

package aep0134

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestAllowMissing(t *testing.T) {
	const singletonPattern = `books/{book}/settings`
	const nonSingletonPattern = `books/{book}`
	// Note: AEP ResourceDescriptor doesn't have a style field, so declarative-friendly
	// detection is not possible. All tests expect nil since resources won't be
	// detected as declarative-friendly.
	for _, test := range []struct {
		name         string
		Pattern      string
		AllowMissing string
		problems     testutils.Problems
	}{
		{"IgnoredNotDF", nonSingletonPattern, "", nil},
		{"ValidIncluded", nonSingletonPattern, "bool allow_missing = 2;", nil},
		{"NoLongerInvalid", nonSingletonPattern, "", nil},
		{"NoLongerInvalidWrongType", nonSingletonPattern, "string allow_missing = 2;", nil},
		{"NoLongerInvalidRepeated", nonSingletonPattern, "repeated bool allow_missing = 2;", nil},
		{"IgnoredSingleton", singletonPattern, "", nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "aep/api/resource.proto";

				service Library {
					rpc UpdateBook(UpdateBookRequest) returns (Book);
				}

				message UpdateBookRequest {
					Book book = 1;
					{{.AllowMissing}}
				}

				message Book {
					option (aep.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "{{.Pattern}}"
					};
				}
			`, test)
			m := f.GetMessageTypes()[0]
			if diff := test.problems.SetDescriptor(m).Diff(allowMissing.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
