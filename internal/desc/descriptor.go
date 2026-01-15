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

// Package desc provides descriptor wrappers around Google's protoreflect
// that are API-compatible with jhump/protoreflect/desc.
package desc

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Descriptor is the common interface implemented by all descriptor types.
type Descriptor interface {
	GetName() string
	GetFullyQualifiedName() string
	GetFile() *FileDescriptor
	GetSourceInfo() *descriptorpb.SourceCodeInfo_Location
	GetParent() Descriptor
	Unwrap() protoreflect.Descriptor
}

// FileDescriptor wraps a protoreflect.FileDescriptor.
type FileDescriptor struct {
	wrapped protoreflect.FileDescriptor
	srcInfo map[string]*descriptorpb.SourceCodeInfo_Location
	proto   *descriptorpb.FileDescriptorProto
	deps    []*FileDescriptor
}

// MessageDescriptor wraps a protoreflect.MessageDescriptor.
type MessageDescriptor struct {
	wrapped protoreflect.MessageDescriptor
	file    *FileDescriptor
}

// FieldDescriptor wraps a protoreflect.FieldDescriptor.
type FieldDescriptor struct {
	wrapped protoreflect.FieldDescriptor
	file    *FileDescriptor
}

// MethodDescriptor wraps a protoreflect.MethodDescriptor.
type MethodDescriptor struct {
	wrapped protoreflect.MethodDescriptor
	file    *FileDescriptor
}

// ServiceDescriptor wraps a protoreflect.ServiceDescriptor.
type ServiceDescriptor struct {
	wrapped protoreflect.ServiceDescriptor
	file    *FileDescriptor
}

// EnumDescriptor wraps a protoreflect.EnumDescriptor.
type EnumDescriptor struct {
	wrapped protoreflect.EnumDescriptor
	file    *FileDescriptor
}

// EnumValueDescriptor wraps a protoreflect.EnumValueDescriptor.
type EnumValueDescriptor struct {
	wrapped protoreflect.EnumValueDescriptor
	file    *FileDescriptor
}

// OneofDescriptor wraps a protoreflect.OneofDescriptor.
type OneofDescriptor struct {
	wrapped protoreflect.OneofDescriptor
	file    *FileDescriptor
}

func pathToString(path []int32) string {
	if len(path) == 0 {
		return ""
	}
	parts := make([]string, len(path))
	for i, p := range path {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ",")
}

func pathFromSourcePath(path protoreflect.SourcePath) []int32 {
	result := make([]int32, len(path))
	for i, p := range path {
		result[i] = int32(p)
	}
	return result
}
