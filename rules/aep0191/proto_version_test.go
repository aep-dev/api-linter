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

package aep0191

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestSyntax(t *testing.T) {
	// Set up the two permutations.
	tests := []struct {
		testName string
		edition  descriptorpb.Edition
		problems testutils.Problems
	}{
		{"Valid (proto3)", descriptorpb.Edition_EDITION_PROTO3, testutils.Problems{}},
		{"Valid (2023)", descriptorpb.Edition_EDITION_2023, testutils.Problems{}},
		{"Valid (2024)", descriptorpb.Edition_EDITION_2024, testutils.Problems{}},
		{"Invalid (proto2)", descriptorpb.Edition_EDITION_PROTO2, testutils.Problems{{Suggestion: `edition = "2023";`}}},
	}

	// Run each permutation as an individual test.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Build an appropriate file descriptor.
			f, err := builder.NewFile("library.proto").SetEdition(test.edition).Build()
			if err != nil {
				t.Fatalf("Could not build file descriptor.")
			}
			// Lint the file, and ensure we got the expected problems.
			if diff := test.problems.SetDescriptor(f).Diff(syntax.Lint(f)); diff != "" {
				t.Error(diff)
			}
		})
	}
}
