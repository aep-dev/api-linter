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

package aep0121

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestNoMutableCycles(t *testing.T) {

	for _, test := range []struct {
		name                                                                             string
		BookExtensions, PublisherExtensions, LibraryExtensions, OtherPublisherExtensions string
		problems                                                                         testutils.Problems
	}{
		{
			"ValidNoCycle",
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Library"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Library"]`,
			"",
			"",
			nil,
		},
		{
			"InvalidCycle",
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Book"]`,
			"",
			"",
			testutils.Problems{{
				Message: "cycle",
			}},
		},
		{
			"InvalidSelfReferenceCycle",
			"",
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"]`,
			"",
			"",
			testutils.Problems{{
				Message: "cycle",
			}},
		},
		{
			"InvalidDeepCycle",
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Library"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Book"]`,
			"",
			testutils.Problems{{
				Message: "cycle",
			}},
		},
		{
			"InvalidDeepAndShallowCycles",
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Library"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Book"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Book"]`,
			testutils.Problems{
				{
					Message: "cycle",
				},
				{
					Message: "cycle",
				},
			},
		},
		{
			"ValidOutputOnlyCyclicReference",
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"]`,
			`[
				(aep.api.field_info).resource_reference = "library.googleapis.com/Book",
				(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_OUTPUT_ONLY
			]`,
			"",
			"",
			nil,
		},
		{
			"ValidOutputOnlyDeepCyclicReference",
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Publisher"]`,
			`[(aep.api.field_info).resource_reference = "library.googleapis.com/Library"]`,
			`[
				(aep.api.field_info).resource_reference = "library.googleapis.com/Book",
				(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_OUTPUT_ONLY
			]`,
			"",
			nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
			import "google/api/resource.proto";
			import "aep/api/field_info.proto";
			message Book {
				option (aep.api.resource) = {
					type: "library.googleapis.com/Book"
					pattern: "publishers/{publisher}/books/{book}"
				};
				string name = 1;

				string resource = 2 {{.BookExtensions}};
			}

			message Publisher {
				option (aep.api.resource) = {
					type: "library.googleapis.com/Publisher"
					pattern: "publishers/{publisher}"
				};
				string name = 1;

				string resource = 2 {{.PublisherExtensions}};

				string other_resource = 3 {{.OtherPublisherExtensions}};
			}

			message Library {
				option (aep.api.resource) = {
					type: "library.googleapis.com/Library"
					pattern: "libraries/{library}"
				};
				string name = 1;

				string resource = 3 {{.LibraryExtensions}};
			}
			`, test)

			msg := f.FindMessage("Publisher")
			want := test.problems
			if len(want) >= 1 {
				want[0].Descriptor = msg.FindFieldByName("resource")
			}
			if len(want) == 2 {
				want[1].Descriptor = msg.FindFieldByName("other_resource")
			}
			// If this rule was run on the entire test file, there would be two
			// findings, one for each resource in the cycle. To simplify that,
			// we just lint one of the offending messages.
			if diff := want.Diff(noMutableCycles.LintMessage(msg)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
