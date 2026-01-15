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

package desc

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// EnumDescriptor methods

// GetName returns the enum name.
func (ed *EnumDescriptor) GetName() string {
	return string(ed.wrapped.Name())
}

// GetFullyQualifiedName returns the fully qualified enum name.
func (ed *EnumDescriptor) GetFullyQualifiedName() string {
	return string(ed.wrapped.FullName())
}

// GetFile returns the file containing this enum.
func (ed *EnumDescriptor) GetFile() *FileDescriptor {
	return ed.file
}

// GetParent returns the parent descriptor.
func (ed *EnumDescriptor) GetParent() Descriptor {
	parent := ed.wrapped.Parent()
	if parent == nil {
		return nil
	}

	switch p := parent.(type) {
	case protoreflect.FileDescriptor:
		return ed.file
	case protoreflect.MessageDescriptor:
		return &MessageDescriptor{wrapped: p, file: ed.file}
	default:
		return nil
	}
}

// GetSourceInfo returns the source code info for this enum.
func (ed *EnumDescriptor) GetSourceInfo() *descriptorpb.SourceCodeInfo_Location {
	path := ed.wrapped.ParentFile().SourceLocations().ByDescriptor(ed.wrapped)
	if path.Path == nil {
		return nil
	}
	return ed.file.srcInfo[pathToString(pathFromSourcePath(path.Path))]
}

// Unwrap returns the underlying protoreflect.Descriptor.
func (ed *EnumDescriptor) Unwrap() protoreflect.Descriptor {
	return ed.wrapped
}

// GetEnumOptions returns the enum options.
func (ed *EnumDescriptor) GetEnumOptions() *descriptorpb.EnumOptions {
	return ed.wrapped.Options().(*descriptorpb.EnumOptions)
}

// GetValues returns all enum values.
func (ed *EnumDescriptor) GetValues() []*EnumValueDescriptor {
	values := ed.wrapped.Values()
	result := make([]*EnumValueDescriptor, values.Len())
	for i := 0; i < values.Len(); i++ {
		result[i] = &EnumValueDescriptor{
			wrapped: values.Get(i),
			file:    ed.file,
		}
	}
	return result
}

// EnumValueDescriptor methods

// GetName returns the enum value name.
func (evd *EnumValueDescriptor) GetName() string {
	return string(evd.wrapped.Name())
}

// GetFullyQualifiedName returns the fully qualified enum value name.
func (evd *EnumValueDescriptor) GetFullyQualifiedName() string {
	return string(evd.wrapped.FullName())
}

// GetFile returns the file containing this enum value.
func (evd *EnumValueDescriptor) GetFile() *FileDescriptor {
	return evd.file
}

// GetParent returns the parent enum.
func (evd *EnumValueDescriptor) GetParent() Descriptor {
	parent := evd.wrapped.Parent()
	if parent == nil {
		return nil
	}

	if ed, ok := parent.(protoreflect.EnumDescriptor); ok {
		return &EnumDescriptor{wrapped: ed, file: evd.file}
	}
	return nil
}

// GetSourceInfo returns the source code info for this enum value.
func (evd *EnumValueDescriptor) GetSourceInfo() *descriptorpb.SourceCodeInfo_Location {
	path := evd.wrapped.ParentFile().SourceLocations().ByDescriptor(evd.wrapped)
	if path.Path == nil {
		return nil
	}
	return evd.file.srcInfo[pathToString(pathFromSourcePath(path.Path))]
}

// Unwrap returns the underlying protoreflect.Descriptor.
func (evd *EnumValueDescriptor) Unwrap() protoreflect.Descriptor {
	return evd.wrapped
}

// GetEnumValueOptions returns the enum value options.
func (evd *EnumValueDescriptor) GetEnumValueOptions() *descriptorpb.EnumValueOptions {
	return evd.wrapped.Options().(*descriptorpb.EnumValueOptions)
}

// GetNumber returns the enum value number.
func (evd *EnumValueDescriptor) GetNumber() int32 {
	return int32(evd.wrapped.Number())
}

// GetEnum returns the parent enum.
func (evd *EnumValueDescriptor) GetEnum() *EnumDescriptor {
	parent := evd.wrapped.Parent()
	if parent == nil {
		return nil
	}

	if ed, ok := parent.(protoreflect.EnumDescriptor); ok {
		return &EnumDescriptor{
			wrapped: ed,
			file:    evd.file,
		}
	}
	return nil
}
