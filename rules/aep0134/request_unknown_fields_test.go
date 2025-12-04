// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/aep-dev/api-linter/internal/desc"

	fpb "google.golang.org/genproto/protobuf/field_mask"
)

func TestUnknownFields(t *testing.T) {
	// Get the correct message type for google.protobuf.FieldMask.
	fieldMask, err := desc.LoadMessageDescriptorForMessage(&fpb.FieldMask{})
	if err != nil {
		t.Fatalf("Unable to load the field mask message.")
	}

	// Set up the testing permutations.
	tests := []struct {
		testName    string
		messageName string
		fieldName   string
		fieldType   *testutils.FieldType
		problems    testutils.Problems
	}{
		// Use BigBook instead of Book to test correct casing logic
		{
			"UpdateMask", "UpdateBigBookRequest", "update_mask",
			testutils.FieldTypeImportedMessage(fieldMask),
			testutils.Problems{},
		},
		{
			"ValidateOnly", "UpdateBigBookRequest", "validate_only",
			testutils.FieldTypeBool(),
			testutils.Problems{},
		},
		{
			"PathOnly", "UpdateBigBookRequest", "path",
			testutils.FieldTypeString(),
			testutils.Problems{},
		},
		{
			"Invalid", "UpdateBigBookRequest", "application_id",
			testutils.FieldTypeString(),
			testutils.Problems{{Message: "Unexpected field"}},
		},
		{
			"InvalidCasing", "UpdateBigBookRequest", "bigbook",
			testutils.FieldTypeString(),
			testutils.Problems{{Message: "Unexpected field"}},
		},
		{
			"Irrelevant", "AcquireBigBookRequest", "application_id",
			testutils.FieldTypeString(),
			testutils.Problems{},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			message, err := testutils.NewMessage(t, test.messageName).AddField(
				/* field: "big_book", testutils.FieldTypeMessage(testutils.NewMessage(t, "BigBook" */)),
			).AddField(
				/* field: test.fieldName, test.fieldType */,
			).Build()
			if err != nil {
				t.Fatalf("Could not build UpdateBookRequest message.")
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
