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

package aep0135

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"

)

func TestUnknownFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName    string
		messageName string
		fieldName   string
		fieldType   *testutils.FieldType
		problems    testutils.Problems
	}{
		{"Force", "DeleteBookRequest", "force", testutils.FieldTypeBool(), testutils.Problems{}},
		{"Etag", "DeleteBookRequest", "etag", testutils.FieldTypeString(), testutils.Problems{}},
		{"AllowMissing", "DeleteBookRequest", "allow_missing", testutils.FieldTypeBool(), testutils.Problems{}},
		{"RequestId", "DeleteBookRequest", "request_id", testutils.FieldTypeString(), testutils.Problems{}},
		{"ValidateOnly", "DeleteBookRequest", "validate_only", testutils.FieldTypeBool(), testutils.Problems{}},
		{"Invalid", "DeleteBookRequest", "application_id", testutils.FieldTypeString(), testutils.Problems{{
			Message: "Unexpected field",
		}}},
		{"Irrelevant", "RemoveBookRequest", "application_id", testutils.FieldTypeString(), testutils.Problems{}},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := testutils.NewMessage(t, test.messageName).AddField(
				/* field: "path", testutils.FieldTypeString( */),
			).AddField(
				/* field: test.fieldName, test.fieldType */,
			).Build()
			if err != nil {
				t.Fatalf("Could not build DeleteBookRequest message.")
			}

			// Run the lint rule, and establish that it returns the correct problems.
			wantProblems := test.problems.SetDescriptor(message.FindFieldByName(test.fieldName))
			gotProblems := unknownFields.Lint(message.GetFile())
			if diff := wantProblems.Diff(gotProblems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
