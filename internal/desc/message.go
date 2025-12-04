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

// GetName returns the message name.
func (md *MessageDescriptor) GetName() string {
	return string(md.wrapped.Name())
}

// GetFullyQualifiedName returns the fully qualified message name.
func (md *MessageDescriptor) GetFullyQualifiedName() string {
	return string(md.wrapped.FullName())
}

// GetFile returns the file containing this message.
func (md *MessageDescriptor) GetFile() *FileDescriptor {
	return md.file
}

// GetParent returns the parent descriptor.
func (md *MessageDescriptor) GetParent() Descriptor {
	parent := md.wrapped.Parent()
	if parent == nil {
		return nil
	}

	switch p := parent.(type) {
	case protoreflect.FileDescriptor:
		return md.file
	case protoreflect.MessageDescriptor:
		return &MessageDescriptor{wrapped: p, file: md.file}
	default:
		return nil
	}
}

// GetSourceInfo returns the source code info for this message.
func (md *MessageDescriptor) GetSourceInfo() *descriptorpb.SourceCodeInfo_Location {
	path := md.wrapped.ParentFile().SourceLocations().ByDescriptor(md.wrapped)
	if path.Path == nil {
		return nil
	}
	return md.file.srcInfo[pathToString(pathFromSourcePath(path.Path))]
}

// Unwrap returns the underlying protoreflect.Descriptor.
func (md *MessageDescriptor) Unwrap() protoreflect.Descriptor {
	return md.wrapped
}

// GetMessageOptions returns the message options.
func (md *MessageDescriptor) GetMessageOptions() *descriptorpb.MessageOptions {
	opts := md.wrapped.Options().(*descriptorpb.MessageOptions)
	return opts
}

// GetFields returns all fields in the message.
func (md *MessageDescriptor) GetFields() []*FieldDescriptor {
	fields := md.wrapped.Fields()
	result := make([]*FieldDescriptor, fields.Len())
	for i := 0; i < fields.Len(); i++ {
		result[i] = &FieldDescriptor{
			wrapped: fields.Get(i),
			file:    md.file,
		}
	}
	return result
}

// GetNestedMessageTypes returns all nested messages.
func (md *MessageDescriptor) GetNestedMessageTypes() []*MessageDescriptor {
	msgs := md.wrapped.Messages()
	result := make([]*MessageDescriptor, msgs.Len())
	for i := 0; i < msgs.Len(); i++ {
		result[i] = &MessageDescriptor{
			wrapped: msgs.Get(i),
			file:    md.file,
		}
	}
	return result
}

// GetNestedEnumTypes returns all nested enums.
func (md *MessageDescriptor) GetNestedEnumTypes() []*EnumDescriptor {
	enums := md.wrapped.Enums()
	result := make([]*EnumDescriptor, enums.Len())
	for i := 0; i < enums.Len(); i++ {
		result[i] = &EnumDescriptor{
			wrapped: enums.Get(i),
			file:    md.file,
		}
	}
	return result
}

// IsMapEntry returns true if this message is a synthetic map entry.
func (md *MessageDescriptor) IsMapEntry() bool {
	return md.wrapped.IsMapEntry()
}

// FindFieldByName finds a field by name.
func (md *MessageDescriptor) FindFieldByName(name string) *FieldDescriptor {
	fields := md.wrapped.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if string(field.Name()) == name {
			return &FieldDescriptor{
				wrapped: field,
				file:    md.file,
			}
		}
	}
	return nil
}
