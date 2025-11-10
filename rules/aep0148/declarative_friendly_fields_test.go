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

package aep0148

import (
	"fmt"
	"testing"

	"bitbucket.org/creachadair/stringset"
	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestDeclarativeFriendlyFields(t *testing.T) {
	for _, test := range []struct {
		name    string
		skipped stringset.Set
	}{
		{"Valid", stringset.New()},
		{"Path", stringset.New("path")},
		{"UID", stringset.New("uid")},
		{"DisplayName", stringset.New("display_name")},
		{"CreateTime", stringset.New("create_time")},
		{"UpdateTime", stringset.New("update_time")},
		{"DeleteTime", stringset.New("delete_time")},
		{"AllTimes", stringset.New("create_time", "update_time", "delete_time")},
		{"Randos", stringset.New("uid", "display_name")},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Set up the string with the fields we will include.
			fields := ""
			cursor := 1
			for fieldName, fieldType := range reqFields {
				if !test.skipped.Contains(fieldName) {
					fields += fmt.Sprintf("  %s %s = %d;\n", fieldType, fieldName, cursor)
					cursor++
				}
			}

			// Note: AEP ResourceDescriptor doesn't have a style field, so declarative-friendly
			// detection is not possible. All tests expect nil since resources won't be
			// detected as declarative-friendly.
			f := testutils.ParseProto3Tmpl(t, `
				import "aep/api/resource.proto";
				import "google/protobuf/timestamp.proto";
				message Book {
					option (aep.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
					};
					{{.Fields}}
				}
			`, struct {
				Fields string
			}{Fields: fields})
			m := f.GetMessageTypes()[0]
			got := declarativeFriendlyRequired.Lint(f)
			if diff := testutils.Problems(nil).SetDescriptor(m).Diff(got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDeclarativeFriendlyFieldsSingleton(t *testing.T) {
	// Note: AEP ResourceDescriptor doesn't have a style field, so declarative-friendly
	// detection is not possible. All tests expect nil since resources won't be
	// detected as declarative-friendly.
	for _, test := range []struct {
		name   string
		Fields string
		want   testutils.Problems
	}{
		{
			"NoLongerInvalidNoCreateTime", `string path = 1; string display_name = 2; google.protobuf.Timestamp update_time = 3;`,
			nil,
		},
		{
			"ValidNoDeleteTimeNoUid", `string path = 1; string display_name = 2; ` +
				`google.protobuf.Timestamp create_time = 3; google.protobuf.Timestamp update_time = 4;`,
			nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "aep/api/resource.proto";
				import "google/protobuf/timestamp.proto";
				message Book {
					option (aep.api.resource) = {
						type: "library.googleapis.com/Settings"
						pattern: "publishers/{publisher}/settings"
					};
					{{.Fields}}
				}
			`, test)
			m := f.GetMessageTypes()[0]
			got := declarativeFriendlyRequired.Lint(f)
			if diff := test.want.SetDescriptor(m).Diff(got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
