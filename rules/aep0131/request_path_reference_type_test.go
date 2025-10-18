// Copyright 2021 Google LLC
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

package aep0131

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestRequestNameReferenceType(t *testing.T) {
	for _, test := range []struct {
		testName   string
		Annotation string
		problems   testutils.Problems
	}{
		{"Valid", `[(aep.api.field_info).resource_reference = "library.googleapis.com/Book"]`, nil},
		{"Invalid", `[(aep.api.field_info).resource_reference = ""]`, testutils.Problems{{Message: "should be a direct"}}},
	} {
		t.Run(test.testName, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
  import "aep/api/field_info.proto";
				message GetBookRequest {
					string path = 1 {{.Annotation}};
				}
			`, test)
			field := f.GetMessageTypes()[0].GetFields()[0]
			if diff := test.problems.SetDescriptor(field).Diff(requestPathReferenceType.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
