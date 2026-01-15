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

package aep0133

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
	"github.com/aep-dev/api-linter/internal/desc"

)

type field struct {
	fieldName string
	fieldType *testutils.FieldType
}

func TestUnknownFields(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName      string
		messageName   string
		messageFields []field
		problems      testutils.Problems
		problemDesc   func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			"Parent",
			"CreateBookRequest",
			[]field{{"parent", testutils.FieldTypeString()}},
			testutils.Problems{},
			nil,
		},
		{
			"ValidateOnly",
			"CreateBookRequest",
			[]field{{"parent", testutils.FieldTypeString()}, {"validate_only", testutils.FieldTypeBool()}},
			testutils.Problems{},
			nil,
		},
		{
			"ResourceRelatedField",
			"CreateBookRequest",
			[]field{
				{"book", testutils.FieldTypeMessage(testutils.NewMessage(t, "Book"))},
				{"id", testutils.FieldTypeString()},
			},
			testutils.Problems{},
			nil,
		},
		{
			"ResourceRelatedField",
			"CreateBookStoreRequest",
			[]field{{"id", testutils.FieldTypeString()}},
			testutils.Problems{},
			nil,
		},
		{
			"RequestIdField",
			"CreateBookRequest",
			[]field{{"request_id", testutils.FieldTypeString()}},
			testutils.Problems{},
			nil,
		},
		{
			"Invalid",
			"CreateBookRequest",
			[]field{{"name", testutils.FieldTypeString()}},
			testutils.Problems{{Message: "Create RPCs must only contain fields explicitly described in AEPs, not \"name\"."}},
			func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("name")
			},
		},
		{
			"InvalidResourceRelatedField",
			"CreateBookStoreRequest",
			[]field{{"book_id", testutils.FieldTypeString()}},
			testutils.Problems{{Message: "Create RPCs must only contain fields explicitly described in AEPs, not \"book_id\"."}},
			func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("book_id")
			},
		},
		{
			"Irrelevant",
			"GetBookRequest",
			[]field{{"name", testutils.FieldTypeString()}},
			testutils.Problems{},
			nil,
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			msgBuilder := testutils.NewMessage(t, test.messageName)
			for _, messageField := range test.messageFields {
				msgBuilder.AddField(
					/* field: messageField.fieldName, messageField.fieldType */,
				)
			}
			message, err := msgBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build GetBookRequest message.")
			}

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor = message
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(message)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := unknownFields.Lint(message.GetFile())
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Error(diff)
			}
		})
	}
}
