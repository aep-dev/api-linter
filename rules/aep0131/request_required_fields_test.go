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

package aep0131

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
		problems             testutils.Problems
	}{
		{
			"ValidNoExtraFields",
			"",
			"",
			nil,
		},
		{
			"ValidOptionalReadMask",
			"google.protobuf.FieldMask read_mask = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_OPTIONAL];",
			"read_mask",
			nil,
		},
		{
			"InvalidRequiredReadMask",
			"google.protobuf.FieldMask read_mask = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];",
			"read_mask",
			testutils.Problems{
				{Message: `Get RPCs must only require fields explicitly described in AEPs, not "read_mask"`},
			},
		},
		{
			"InvalidRequiredUnknownField",
			"bool create_iam = 3 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];",
			"create_iam",
			testutils.Problems{
				{Message: `Get RPCs must only require fields explicitly described in AEPs, not "create_iam"`},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "aep/api/field_info.proto";
				import "aep/api/resource.proto";
				import "google/protobuf/field_mask.proto";

				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							get: "/v1/{path=publishers/*/books/*}"
						};
					}
				}

				message GetBookRequest {
					// The path of the book to retrieve.
					// Format: publishers/{publisher}/books/{book}
					string path = 1 [
					    (aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED,
						(aep.api.field_info).resource_reference = "library.googleapis.com/Book"
					];
					{{.Fields}}
				}

				message Book {
					option (aep.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "publishers/{publisher}/books/{book}"
					};
					string path = 1;
				}
			`, test)
			var dbr desc.Descriptor = f.FindMessage("GetBookRequest")
			if test.problematicFieldName != "" {
				dbr = f.FindMessage("GetBookRequest").FindFieldByName(test.problematicFieldName)
			}
			if diff := test.problems.SetDescriptor(dbr).Diff(requestRequiredFields.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
