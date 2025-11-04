// Copyright 2022 Google LLC
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
)

func TestResourceReferenceType(t *testing.T) {
	// Set up testing permutations.
	tests := []struct {
		testName   string
		Annotation string
		problems   testutils.Problems
	}{
		{"ValidMatch_resource_reference", `resource_reference = "library.googleapis.com/Book"`, nil},
		{"InvalidMismatch_resource_reference", `resource_reference = "library.googleapis.com/Shelf"`, testutils.Problems{{Message: "`resource_reference_child_type`"}}},
		{"ValidMatch_resource_reference_child_type", `resource_reference_child_type = "library.googleapis.com/Book"`, nil},
		{"InvalidMismatch_resource_reference_child_type", `resource_reference_child_type = "library.googleapis.com/Shelf"`, testutils.Problems{{Message: "`resource_reference_child_type`"}}},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
  import "aep/api/field_info.proto";
				import "google/longrunning/operations.proto";
				service Library {
					rpc CreateBook(CreateBookRequest) returns (Book) {}
				}
				message CreateBookRequest {
					string parent = 1 [(aep.api.field_info).{{ .Annotation }}];
				}
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "shelves/{shelf}/books/{book}"
					};
					string name = 1;
				}
			`, test)
			field := file.GetServices()[0].GetMethods()[0].GetInputType().FindFieldByName("parent")
			problems := resourceReferenceType.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestResourceReferenceTypeLRO(t *testing.T) {
	// Set up testing permutations.
	tests := []struct {
		testName     string
		Annotation   string
		ResponseType string
		problems     testutils.Problems
	}{
		{"ValidMatch_resource_reference", `resource_reference = "library.googleapis.com/Book"`, "Book", nil},
		{"InvalidMismatch_resource_reference", `resource_reference = "library.googleapis.com/Shelf"`, "Book", testutils.Problems{{Message: "`resource_reference_child_type`"}}},
		{"ValidMatch_resource_reference_child_type", `resource_reference_child_type = "library.googleapis.com/Book"`, "Book", nil},
		{"InvalidMismatch_resource_reference_child_type", `resource_reference_child_type = "library.googleapis.com/Shelf"`, "Book", testutils.Problems{{Message: "`resource_reference_child_type`"}}},
		{"SkipUnresolvableResponse", `resource_reference = "library.googleapis.com/Shelf"`, "Foo", nil},
	}

	// Run each test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			file := testutils.ParseProto3Tmpl(t, `
				import "google/api/resource.proto";
  import "aep/api/field_info.proto";
				import "google/longrunning/operations.proto";
				service Library {
					rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
						option (google.longrunning.operation_info) = {
							response_type: "{{ .ResponseType }}"
							metadata_type: "Book"
						};
					}
				}
				message CreateBookRequest {
					string parent = 1 [(aep.api.field_info).{{ .Annotation }}];
				}
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "shelves/{shelf}/books/{book}"
					};
					string name = 1;
				}
			`, test)
			field := file.GetServices()[0].GetMethods()[0].GetInputType().FindFieldByName("parent")
			problems := resourceReferenceType.Lint(file)
			if diff := test.problems.SetDescriptor(field).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
