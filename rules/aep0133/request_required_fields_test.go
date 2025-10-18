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

package aep0133

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
)

func TestRequiredFieldTests(t *testing.T) {
	for _, test := range []struct {
		name                 string
		Fields               string
		problematicFieldName string
		Singular             string
		problems             testutils.Problems
	}{
		{
			"ValidNoExtraFields",
			"",
			"",
			"",
			nil,
		},
		{
			"ValidWithSingularNoExtraFields",
			"",
			"",
			"bookShelf",
			nil,
		},
		{
			"ValidWithSingularAndIdField",
			"string id = 3 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_OPTIONAL];",
			"",
			"bookShelf",
			nil,
		},
		{
			"ValidOptionalValidateOnly",
			"string validate_only = 3 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_OPTIONAL];",
			"validate_only",
			"",
			nil,
		},
		{
			"InvalidRequiredValidateOnly",
			"bool validate_only = 3 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];",
			"validate_only",
			"",
			testutils.Problems{
				{Message: `Create RPCs must only require fields explicitly described in AEPs, not "validate_only"`},
			},
		},
		{
			"InvalidRequiredUnknownField",
			"bool create_iam = 3 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];",
			"create_iam",
			"",
			testutils.Problems{
				{Message: `Create RPCs must only require fields explicitly described in AEPs, not "create_iam"`},
			},
		},
		{
			"InvalidRequiredUnknownMessageField",
			"Foo foo = 3 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];",
			"foo",
			"",
			testutils.Problems{
				{Message: `Create RPCs must only require fields explicitly described in AEPs, not "foo"`},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "aep/api/field_info.proto";
				import "google/api/resource.proto";

				service Library {
					rpc CreateBookShelf(CreateBookShelfRequest) returns (BookShelf) {
						option (google.api.http) = {
							delete: "/v1/{name=publishers/*/bookShelves/*}"
						};
					}
				}

				message BookShelf {
					option (google.api.resource) = {
						type: "library.googleapis.com/BookShelf"
						pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
						singular: "{{.Singular}}"
					};
					string name = 1;
				}

				message Foo {}

				message CreateBookShelfRequest {
					string parent = 1 [
						(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED
					];
					BookShelf book_shelf = 2 [
						(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED
					];
					{{.Fields}}
				}
			`, test)
			var dbr desc.Descriptor = f.FindMessage("CreateBookShelfRequest")
			if test.problematicFieldName != "" {
				dbr = f.FindMessage("CreateBookShelfRequest").FindFieldByName(test.problematicFieldName)
			}
			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
