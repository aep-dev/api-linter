// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lint

import (
	"reflect"
	"testing"

	"github.com/aep-dev/api-linter/internal/desc"

)

func TestFileRule(t *testing.T) {
	// Create a file descriptor with nothing in it.
	fd, err := testutils.NewFile(t, "test.proto").Build()
	if err != nil {
		t.Fatalf("Could not build file descriptor: %q", err)
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd) {
		t.Run(test.testName, func(t *testing.T) {
			rule := &FileRule{
				Name: RuleName("test"),
				OnlyIf: func(fd *desc.FileDescriptor) bool {
					return fd.GetName() == "test.proto"
				},
				LintFile: func(fd *desc.FileDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestMessageRule(t *testing.T) {
	// Create a file descriptor with two messages in it.
	fd, err := testutils.NewFile(t, "test.proto").AddMessage(
		testutils.NewMessage(t, "Book"),
	).AddMessage(
		testutils.NewMessage(t, "Author"),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			rule := &MessageRule{
				Name: RuleName("test"),
				OnlyIf: func(m *desc.MessageDescriptor) bool {
					return m.GetName() == "Author"
				},
				LintMessage: func(m *desc.MessageDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

// Establish that nested messages are tested.
func TestMessageRuleNested(t *testing.T) {
	// Create a file descriptor with a message and nested message in it.
	fd, err := testutils.NewFile(t, "test.proto").AddMessage(
		testutils.NewMessage(t, "Book").AddNestedMessage(testutils.NewMessage(t, "Author")),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[0].GetNestedMessageTypes()[0]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the message rule.
			rule := &MessageRule{
				Name: RuleName("test"),
				OnlyIf: func(m *desc.MessageDescriptor) bool {
					return m.GetName() == "Author"
				},
				LintMessage: func(m *desc.MessageDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestFieldRule(t *testing.T) {
	// Create a file descriptor with one message and two fields in that message.
	fd, err := testutils.NewFile(t, "test.proto").AddMessage(
		testutils.NewMessage(t, "Book").AddField(
			newField("title", "string", 1),
		).AddField(
			newField("edition_count", "int32", 2),
		),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build file descriptor.")
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[0].GetFields()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the field rule.
			rule := &FieldRule{
				Name: RuleName("test"),
				OnlyIf: func(f *desc.FieldDescriptor) bool {
					return f.GetName() == "edition_count"
				},
				LintField: func(f *desc.FieldDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestServiceRule(t *testing.T) {
	// Create a file descriptor with a service.
	fd, err := testutils.NewFile(t, "test.proto").AddService(
		builder.NewService("Library"),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build a file descriptor: %q", err)
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetServices()[0]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the service rule.
			rule := &ServiceRule{
				Name: RuleName("test"),
				LintService: func(s *desc.ServiceDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestMethodRule(t *testing.T) {
	// Create a file descriptor with a service.
	book := testutils.NewMessage(t, "Book")
	getBookRequest := testutils.NewMessage(t, "GetBookRequest")
	createBookRequest := testutils.NewMessage(t, "CreateBookRequest")

	fd, err := testutils.NewFile(t, "test.proto").AddMessage(book).AddMessage(getBookRequest).AddMessage(createBookRequest).AddService(
		builder.NewService("Library").AddMethod(
			builder.NewMethod(
				"GetBook",
				builder.RpcTypeMessage(getBookRequest, false),
				builder.RpcTypeMessage(book, false),
			),
		).AddMethod(
			builder.NewMethod(
				"CreateBook",
				builder.RpcTypeMessage(createBookRequest, false),
				builder.RpcTypeMessage(book, false),
			),
		),
	).Build()
	if err != nil {
		t.Fatalf("Failed to build a file descriptor: %q", err)
	}

	// Iterate over the tests and run them.
	for _, test := range makeLintRuleTests(fd.GetServices()[0].GetMethods()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the method rule.
			rule := &MethodRule{
				Name: RuleName("test"),
				OnlyIf: func(m *desc.MethodDescriptor) bool {
					return m.GetName() == "CreateBook"
				},
				LintMethod: func(m *desc.MethodDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestEnumRule(t *testing.T) {
	// Create a file descriptor with top-level enums.
	fd, err := testutils.NewFile(t, "test.proto").AddEnum(
		newEnum("Format").AddValue(newEnumValue("PDF", 0)),
	).AddEnum(
		newEnum("Edition").AddValue(newEnumValue("PUBLISHER_ONLY", 0)),
	).Build()
	if err != nil {
		t.Fatalf("Error building test proto:%s ", err)
	}

	for _, test := range makeLintRuleTests(fd.GetEnumTypes()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum rule.
			rule := &EnumRule{
				Name: RuleName("test"),
				OnlyIf: func(e *desc.EnumDescriptor) bool {
					return e.GetName() == "Edition"
				},
				LintEnum: func(e *desc.EnumDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestEnumValueRule(t *testing.T) {
	// Create a file descriptor with a top-level enum with values.
	fd, err := testutils.NewFile(t, "test.proto").AddEnum(
		newEnum("Format").AddValue(newEnumValue("YAML", 0)).AddValue(newEnumValue("JSON", 1)),
	).Build()
	if err != nil {
		t.Fatalf("Error building test proto:%s ", err)
	}

	for _, test := range makeLintRuleTests(fd.GetEnumTypes()[0].GetValues()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum value rule.
			rule := &EnumValueRule{
				Name: RuleName("test"),
				OnlyIf: func(e *desc.EnumValueDescriptor) bool {
					return e.GetName() == "JSON"
				},
				LintEnumValue: func(e *desc.EnumValueDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestEnumRuleNested(t *testing.T) {
	// Create a file descriptor with top-level enums.
	fd, err := testutils.NewFile(t, "test.proto").AddMessage(
		testutils.NewMessage(t, "Book").AddNestedEnum(
			newEnum("Format").AddValue(newEnumValue("PDF", 0)),
		).AddNestedEnum(
			newEnum("Edition").AddValue(newEnumValue("PUBLISHER_ONLY", 0)),
		),
	).Build()
	if err != nil {
		t.Fatalf("Error building test proto:%s ", err)
	}

	for _, test := range makeLintRuleTests(fd.GetMessageTypes()[0].GetNestedEnumTypes()[1]) {
		t.Run(test.testName, func(t *testing.T) {
			// Create the enum rule.
			rule := &EnumRule{
				Name: RuleName("test"),
				OnlyIf: func(e *desc.EnumDescriptor) bool {
					return e.GetName() == "Edition"
				},
				LintEnum: func(e *desc.EnumDescriptor) []Problem {
					return test.problems
				},
			}

			// Run the rule and assert that we got what we expect.
			test.runRule(rule, fd, t)
		})
	}
}

func TestDescriptorRule(t *testing.T) {
	// Create a file with one of everything in it.
	book := testutils.NewMessage(t, "Book").AddNestedEnum(
		newEnum("Format").AddValue(
			newEnumValue("FORMAT_UNSPECIFIED", 0),
		).AddValue(newEnumValue("PAPERBACK", 1)),
	).AddField(newField("name", "string", 1)).AddNestedMessage(
		testutils.NewMessage(t, "Author"),
	)
	fd, err := testutils.NewFile(t, "library.proto").AddMessage(book).AddService(
		builder.NewService("Library").AddMethod(
			builder.NewMethod(
				"ConjureBook",
				builder.RpcTypeMessage(book, false),
				builder.RpcTypeMessage(book, false),
			),
		),
	).AddEnum(newEnum("State").AddValue(newEnumValue("AVAILABLE", 0))).Build()
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Create a rule that lets us verify that each descriptor was visited.
	visited := make(map[string]desc.Descriptor)
	rule := &DescriptorRule{
		Name: RuleName("test"),
		OnlyIf: func(d desc.Descriptor) bool {
			return d.GetName() != "FORMAT_UNSPECIFIED"
		},
		LintDescriptor: func(d desc.Descriptor) []Problem {
			visited[d.GetName()] = d
			return nil
		},
	}

	// Run the rule.
	rule.Lint(fd)

	// Verify that each descriptor was visited.
	// We do not care what order they were visited in.
	wantDescriptors := []string{
		"Author", "Book", "ConjureBook", "Format", "PAPERBACK",
		"name", "Library", "State", "AVAILABLE",
	}
	if got, want := rule.GetName(), "test"; string(got) != want {
		t.Errorf("Got name %q, wanted %q", got, want)
	}
	if got, want := len(visited), len(wantDescriptors); got != want {
		t.Errorf("Got %d descriptors, wanted %d.", got, want)
	}
	for _, name := range wantDescriptors {
		if _, ok := visited[name]; !ok {
			t.Errorf("Missing descriptor %q.", name)
		}
	}
}

type lintRuleTest struct {
	testName string
	problems []Problem
}

// runRule runs a rule within a test environment.
func (test *lintRuleTest) runRule(rule ProtoRule, fd *desc.FileDescriptor, t *testing.T) {
	// Establish that the metadata methods work.
	if got, want := string(rule.GetName()), string(RuleName("test")); got != want {
		t.Errorf("Got %q for GetName(), expected %q", got, want)
	}

	// Run the rule's lint function on the file descriptor
	// and assert that we got what we expect.
	if got, want := rule.Lint(fd), test.problems; !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v problems; expected %v.", got, want)
	}
}

// makeLintRuleTests generates boilerplate tests that are consistent for
// each type of rule.
func makeLintRuleTests(d desc.Descriptor) []lintRuleTest {
	return []lintRuleTest{
		{"NoProblems", []Problem{}},
		{"OneProblem", []Problem{{
			Message:    "There was a problem.",
			Descriptor: d,
		}}},
		{"TwoProblems", []Problem{
			{
				Message:    "This was the first problem.",
				Descriptor: d,
			},
			{
				Message:    "This was the second problem.",
				Descriptor: d,
			},
		}},
	}
}
