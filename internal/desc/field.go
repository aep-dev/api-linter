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

// GetName returns the field name.
func (fd *FieldDescriptor) GetName() string {
	return string(fd.wrapped.Name())
}

// GetFullyQualifiedName returns the fully qualified field name.
func (fd *FieldDescriptor) GetFullyQualifiedName() string {
	return string(fd.wrapped.FullName())
}

// GetFile returns the file containing this field.
func (fd *FieldDescriptor) GetFile() *FileDescriptor {
	return fd.file
}

// GetParent returns the parent message.
func (fd *FieldDescriptor) GetParent() Descriptor {
	parent := fd.wrapped.Parent()
	if parent == nil {
		return nil
	}

	if md, ok := parent.(protoreflect.MessageDescriptor); ok {
		return &MessageDescriptor{wrapped: md, file: fd.file}
	}
	return nil
}

// GetSourceInfo returns the source code info for this field.
func (fd *FieldDescriptor) GetSourceInfo() *descriptorpb.SourceCodeInfo_Location {
	path := fd.wrapped.ParentFile().SourceLocations().ByDescriptor(fd.wrapped)
	if path.Path == nil {
		return nil
	}
	return fd.file.srcInfo[pathToString(pathFromSourcePath(path.Path))]
}

// Unwrap returns the underlying protoreflect.Descriptor.
func (fd *FieldDescriptor) Unwrap() protoreflect.Descriptor {
	return fd.wrapped
}

// GetFieldOptions returns the field options.
func (fd *FieldDescriptor) GetFieldOptions() *descriptorpb.FieldOptions {
	return fd.wrapped.Options().(*descriptorpb.FieldOptions)
}

// GetNumber returns the field number.
func (fd *FieldDescriptor) GetNumber() int32 {
	return int32(fd.wrapped.Number())
}

// GetType returns the field type.
func (fd *FieldDescriptor) GetType() descriptorpb.FieldDescriptorProto_Type {
	return descriptorpb.FieldDescriptorProto_Type(fd.wrapped.Kind())
}

// IsRepeated returns true if the field is repeated.
func (fd *FieldDescriptor) IsRepeated() bool {
	return fd.wrapped.Cardinality() == protoreflect.Repeated
}

// GetMessageType returns the message type if this is a message field.
func (fd *FieldDescriptor) GetMessageType() *MessageDescriptor {
	if fd.wrapped.Kind() != protoreflect.MessageKind && fd.wrapped.Kind() != protoreflect.GroupKind {
		return nil
	}

	md := fd.wrapped.Message()
	if md == nil {
		return nil
	}

	return &MessageDescriptor{
		wrapped: md,
		file:    fd.file,
	}
}

// GetEnumType returns the enum type if this is an enum field.
func (fd *FieldDescriptor) GetEnumType() *EnumDescriptor {
	if fd.wrapped.Kind() != protoreflect.EnumKind {
		return nil
	}

	ed := fd.wrapped.Enum()
	if ed == nil {
		return nil
	}

	return &EnumDescriptor{
		wrapped: ed,
		file:    fd.file,
	}
}

// GetOneOf returns the oneof containing this field, or nil.
func (fd *FieldDescriptor) GetOneOf() *OneofDescriptor {
	oneof := fd.wrapped.ContainingOneof()
	if oneof == nil {
		return nil
	}

	return &OneofDescriptor{
		wrapped: oneof,
		file:    fd.file,
	}
}

// GetMapKeyType returns the key type for map fields.
func (fd *FieldDescriptor) GetMapKeyType() *FieldDescriptor {
	if !fd.wrapped.IsMap() {
		return nil
	}

	msg := fd.wrapped.Message()
	if msg == nil {
		return nil
	}

	keyField := msg.Fields().Get(0) // Map entry key is always first field
	return &FieldDescriptor{
		wrapped: keyField,
		file:    fd.file,
	}
}

// GetMapValueType returns the value type for map fields.
func (fd *FieldDescriptor) GetMapValueType() *FieldDescriptor {
	if !fd.wrapped.IsMap() {
		return nil
	}

	msg := fd.wrapped.Message()
	if msg == nil {
		return nil
	}

	valueField := msg.Fields().Get(1) // Map entry value is always second field
	return &FieldDescriptor{
		wrapped: valueField,
		file:    fd.file,
	}
}

// GetOwner returns the message that owns this field.
func (fd *FieldDescriptor) GetOwner() *MessageDescriptor {
	parent := fd.wrapped.Parent()
	if parent == nil {
		return nil
	}

	if md, ok := parent.(protoreflect.MessageDescriptor); ok {
		return &MessageDescriptor{
			wrapped: md,
			file:    fd.file,
		}
	}
	return nil
}

// IsProto3Optional returns true if this field is a proto3 optional field.
func (fd *FieldDescriptor) IsProto3Optional() bool {
	return fd.wrapped.HasOptionalKeyword()
}
