// Copyright 2024 Google LLC
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

package lint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aep-dev/api-linter/internal/desc"
	"github.com/aep-dev/api-linter/internal/protoparse"
	"google.golang.org/protobuf/types/descriptorpb"
)

// testutils provides builder-style test helpers for creating proto descriptors.
var testutils = &testutilsHelper{}

type testutilsHelper struct{}

// NewFile creates a new file builder.
func (testutilsHelper) NewFile(t *testing.T, name string) *fileBuilder {
	return &fileBuilder{
		t:        t,
		name:     name,
		messages: []string{},
		enums:    []string{},
		services: []string{},
	}
}

// NewMessage creates a new message builder.
func (testutilsHelper) NewMessage(t *testing.T, name string) *messageBuilder {
	return &messageBuilder{
		t:       t,
		name:    name,
		fields:  []string{},
		nested:  []string{},
		comment: "",
	}
}

type fileBuilder struct {
	t              *testing.T
	name           string
	syntaxComment  string
	messages       []string
	enums          []string
	services       []string
}

func (fb *fileBuilder) SetSyntaxComments(c Comments) *fileBuilder {
	fb.syntaxComment = c.LeadingComment
	return fb
}

func (fb *fileBuilder) AddMessage(mb *messageBuilder) *fileBuilder {
	fb.messages = append(fb.messages, mb.build())
	return fb
}

func (fb *fileBuilder) AddEnum(eb *enumBuilder) *fileBuilder {
	fb.enums = append(fb.enums, eb.build())
	return fb
}

func (fb *fileBuilder) AddService(sb *serviceBuilder) *fileBuilder {
	fb.services = append(fb.services, sb.build())
	return fb
}

func (fb *fileBuilder) Build() (*desc.FileDescriptor, error) {
	var parts []string

	// Add syntax with optional comment
	if fb.syntaxComment != "" {
		// Handle multi-line comments
		for _, line := range strings.Split(fb.syntaxComment, "\n") {
			parts = append(parts, "// "+line)
		}
	}
	parts = append(parts, "syntax = \"proto3\";")
	parts = append(parts, "")

	// Add messages
	for _, msg := range fb.messages {
		parts = append(parts, msg)
		parts = append(parts, "")
	}

	// Add enums
	for _, enum := range fb.enums {
		parts = append(parts, enum)
		parts = append(parts, "")
	}

	// Add services
	for _, svc := range fb.services {
		parts = append(parts, svc)
		parts = append(parts, "")
	}

	content := strings.Join(parts, "\n")

	parser := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			fb.name: content,
		}),
		IncludeSourceCodeInfo: true,
		LookupImport:          desc.LoadFileDescriptor,
	}

	fds, err := parser.ParseFiles(fb.name)
	if err != nil {
		return nil, err
	}

	if len(fds) == 0 {
		return nil, fmt.Errorf("no file descriptors returned")
	}

	return fds[0], nil
}

type messageBuilder struct {
	t       *testing.T
	name    string
	comment string
	fields  []string
	nested  []string
	options *descriptorpb.MessageOptions
}

func (mb *messageBuilder) SetComments(c Comments) *messageBuilder {
	mb.comment = c.LeadingComment
	return mb
}

func (mb *messageBuilder) SetOptions(opts *descriptorpb.MessageOptions) *messageBuilder {
	mb.options = opts
	return mb
}

func (mb *messageBuilder) AddField(fb *fieldBuilder) *messageBuilder {
	if fb != nil {
		mb.fields = append(mb.fields, fb.build())
	}
	return mb
}

func (mb *messageBuilder) AddNestedMessage(nmb *messageBuilder) *messageBuilder {
	mb.nested = append(mb.nested, nmb.build())
	return mb
}

func (mb *messageBuilder) AddNestedEnum(eb *enumBuilder) *messageBuilder {
	mb.nested = append(mb.nested, eb.build())
	return mb
}

func (mb *messageBuilder) build() string {
	var parts []string

	if mb.comment != "" {
		for _, line := range strings.Split(mb.comment, "\n") {
			parts = append(parts, "// "+line)
		}
	}

	parts = append(parts, fmt.Sprintf("message %s {", mb.name))

	// Add deprecated option if needed
	if mb.options != nil {
		if mb.options.Deprecated != nil && *mb.options.Deprecated {
			parts = append(parts, "  option deprecated = true;")
		}
	}

	// Add nested definitions
	for _, nested := range mb.nested {
		for _, line := range strings.Split(nested, "\n") {
			parts = append(parts, "  "+line)
		}
	}

	// Add fields
	for i, field := range mb.fields {
		parts = append(parts, fmt.Sprintf("  %s", field))
		_ = i
	}

	parts = append(parts, "}")

	return strings.Join(parts, "\n")
}

type fieldBuilder struct {
	t         *testing.T
	name      string
	fieldType string
	number    int
	options   *descriptorpb.FieldOptions
}

func (fb *fieldBuilder) SetOptions(opts *descriptorpb.FieldOptions) *fieldBuilder {
	fb.options = opts
	return fb
}

func (fb *fieldBuilder) build() string {
	optStr := ""
	if fb.options != nil {
		if fb.options.Deprecated != nil && *fb.options.Deprecated {
			optStr = " [deprecated = true]"
		}
	}
	return fmt.Sprintf("%s %s = %d%s;", fb.fieldType, fb.name, fb.number, optStr)
}

type enumBuilder struct {
	t      *testing.T
	name   string
	values []string
}

func (eb *enumBuilder) AddValue(evb *enumValueBuilder) *enumBuilder {
	if evb != nil {
		eb.values = append(eb.values, evb.build())
	}
	return eb
}

func (eb *enumBuilder) build() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("enum %s {", eb.name))

	for i, value := range eb.values {
		parts = append(parts, fmt.Sprintf("  %s", value))
		_ = i
	}

	parts = append(parts, "}")

	return strings.Join(parts, "\n")
}

type enumValueBuilder struct {
	t      *testing.T
	name   string
	number int
}

func (evb *enumValueBuilder) build() string {
	return fmt.Sprintf("%s = %d;", evb.name, evb.number)
}

type serviceBuilder struct {
	t       *testing.T
	name    string
	methods []string
}

func (sb *serviceBuilder) AddMethod(mb *methodBuilder) *serviceBuilder {
	if mb != nil {
		sb.methods = append(sb.methods, mb.build())
	}
	return sb
}

func (sb *serviceBuilder) build() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("service %s {", sb.name))

	for _, method := range sb.methods {
		parts = append(parts, fmt.Sprintf("  %s", method))
	}

	parts = append(parts, "}")

	return strings.Join(parts, "\n")
}

type methodBuilder struct {
	t                *testing.T
	name             string
	inputType        string
	outputType       string
	clientStreaming  bool
	serverStreaming  bool
}

func (mb *methodBuilder) build() string {
	inputPrefix := ""
	outputPrefix := ""
	if mb.clientStreaming {
		inputPrefix = "stream "
	}
	if mb.serverStreaming {
		outputPrefix = "stream "
	}
	return fmt.Sprintf("rpc %s(%s%s) returns (%s%s);", mb.name, inputPrefix, mb.inputType, outputPrefix, mb.outputType)
}

// Comments represents proto comments.
type Comments struct {
	LeadingComment string
}

// builder namespace for compatibility
var builder = &builderHelper{}

type builderHelper struct{}

func (builderHelper) NewService(name string) *serviceBuilder {
	return &serviceBuilder{
		name:    name,
		methods: []string{},
	}
}

func (builderHelper) NewMethod(name string, inputType *rpcType, outputType *rpcType) *methodBuilder {
	return &methodBuilder{
		name:            name,
		inputType:       inputType.typeName,
		outputType:      outputType.typeName,
		clientStreaming: inputType.streaming,
		serverStreaming: outputType.streaming,
	}
}

type rpcType struct {
	typeName  string
	streaming bool
}

func (builderHelper) RpcTypeMessage(mb *messageBuilder, streaming bool) *rpcType {
	return &rpcType{
		typeName:  mb.name,
		streaming: streaming,
	}
}

func (builderHelper) Comments(leading string) Comments {
	return Comments{LeadingComment: leading}
}

// newField creates a field builder with a given name, type, and number.
func newField(name, fieldType string, number int) *fieldBuilder {
	return &fieldBuilder{
		name:      name,
		fieldType: fieldType,
		number:    number,
	}
}

// newEnum creates an enum builder with a given name.
func newEnum(name string) *enumBuilder {
	return &enumBuilder{
		name:   name,
		values: []string{},
	}
}

// newEnumValue creates an enum value builder.
func newEnumValue(name string, number int) *enumValueBuilder {
	return &enumValueBuilder{
		name:   name,
		number: number,
	}
}
