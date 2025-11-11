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

package aep0004

import (
	"testing"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestDuplicateResource(t *testing.T) {
	f := testutils.ParseProto3Tmpls(t, map[string]string{
		"dep.proto": `
			import "aep/api/resource.proto";
			package xyz;
			message Publisher {
				option (aep.api.resource) = { type: "library.googleapis.com/Publisher" };
			}
			`,
		"test.proto": `
			import "dep.proto";
			import "aep/api/resource.proto";
			package abc;
			message Book {
				option (aep.api.resource) = { type: "library.googleapis.com/Book" };
			}
			message Publisher {
				option (aep.api.resource) = { type: "library.googleapis.com/Publisher" };
			}
			message Foo {
				message Tome {
					option (aep.api.resource) = { type: "library.googleapis.com/Book" };
				}
			}`,
	}, nil)["test.proto"]
	want := testutils.Problems{
		lint.Problem{
			Message:    "Multiple definitions for resource \"library.googleapis.com/Book\": message `abc.Book`, message `abc.Foo.Tome`.",
			Descriptor: f.GetMessageTypes()[0],
		},
		lint.Problem{
			Message:    "Multiple definitions for resource \"library.googleapis.com/Book\": message `abc.Book`, message `abc.Foo.Tome`.",
			Descriptor: f.GetMessageTypes()[2].GetNestedMessageTypes()[0],
		},
		lint.Problem{
			Message:    "Multiple definitions for resource \"library.googleapis.com/Publisher\": message `abc.Publisher`, message `xyz.Publisher`.",
			Descriptor: f.GetMessageTypes()[1],
		},
	}
	if diff := want.Diff(duplicateResource.Lint(f)); diff != "" {
		t.Fatal(diff)
	}
}
