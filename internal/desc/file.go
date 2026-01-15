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
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Wrap Files wraps multiple protoreflect.FileDescriptor objects.
func WrapFiles(fds []protoreflect.FileDescriptor) ([]*FileDescriptor, error) {
	fileMap := make(map[string]*FileDescriptor)

	// First pass: wrap all files
	for _, fd := range fds {
		wrapped := wrapFileWithoutDeps(fd)
		fileMap[string(fd.Path())] = wrapped
	}

	// Second pass: resolve dependencies
	for _, fd := range fds {
		wrapped := fileMap[string(fd.Path())]
		deps := make([]*FileDescriptor, fd.Imports().Len())
		for i := 0; i < fd.Imports().Len(); i++ {
			imp := fd.Imports().Get(i)
			if depFile, ok := fileMap[string(imp.Path())]; ok {
				deps[i] = depFile
			} else {
				// Try global registry
				if depFd, err := protoregistry.GlobalFiles.FindFileByPath(string(imp.Path())); err == nil {
					deps[i] = wrapFileWithoutDeps(depFd)
				}
			}
		}
		wrapped.deps = deps
	}

	result := make([]*FileDescriptor, 0, len(fileMap))
	for _, fd := range fds {
		result = append(result, fileMap[string(fd.Path())])
	}

	return result, nil
}

// WrapFile wraps a single protoreflect.FileDescriptor.
func WrapFile(fd protoreflect.FileDescriptor) (*FileDescriptor, error) {
	wrapped := wrapFileWithoutDeps(fd)

	// Resolve dependencies
	deps := make([]*FileDescriptor, fd.Imports().Len())
	for i := 0; i < fd.Imports().Len(); i++ {
		imp := fd.Imports().Get(i)
		if depFd, err := protoregistry.GlobalFiles.FindFileByPath(string(imp.Path())); err == nil {
			deps[i] = wrapFileWithoutDeps(depFd)
		}
	}
	wrapped.deps = deps

	return wrapped, nil
}

func wrapFileWithoutDeps(fd protoreflect.FileDescriptor) *FileDescriptor {
	proto := protodesc.ToFileDescriptorProto(fd)

	// Build source info map
	srcInfo := make(map[string]*descriptorpb.SourceCodeInfo_Location)
	if proto.SourceCodeInfo != nil {
		for _, loc := range proto.SourceCodeInfo.Location {
			srcInfo[pathToString(loc.Path)] = loc
		}
	}

	return &FileDescriptor{
		wrapped: fd,
		srcInfo: srcInfo,
		proto:   proto,
	}
}

// LoadFileDescriptor loads a file by name from the global registry.
func LoadFileDescriptor(name string) (*FileDescriptor, error) {
	fd, err := protoregistry.GlobalFiles.FindFileByPath(name)
	if err != nil {
		return nil, err
	}
	return WrapFile(fd)
}

// CreateFileDescriptors creates file descriptors from FileDescriptorProtos.
func CreateFileDescriptors(protos []*descriptorpb.FileDescriptorProto) (map[string]*FileDescriptor, error) {
	files := &protoregistry.Files{}
	fds := make([]protoreflect.FileDescriptor, 0, len(protos))

	for _, proto := range protos {
		fd, err := protodesc.NewFile(proto, files)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fd)
		files.RegisterFile(fd)
	}

	wrapped, err := WrapFiles(fds)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*FileDescriptor)
	for i, proto := range protos {
		result[proto.GetName()] = wrapped[i]
	}

	return result, nil
}

// GetName returns the file name.
func (fd *FileDescriptor) GetName() string {
	return string(fd.wrapped.Path())
}

// GetFullyQualifiedName returns the package name.
func (fd *FileDescriptor) GetFullyQualifiedName() string {
	return string(fd.wrapped.Package())
}

// GetPackage returns the package name.
func (fd *FileDescriptor) GetPackage() string {
	return string(fd.wrapped.Package())
}

// GetFile returns self.
func (fd *FileDescriptor) GetFile() *FileDescriptor {
	return fd
}

// GetParent returns nil (files have no parent).
func (fd *FileDescriptor) GetParent() Descriptor {
	return nil
}

// GetSourceInfo returns the source code info for the file.
func (fd *FileDescriptor) GetSourceInfo() *descriptorpb.SourceCodeInfo_Location {
	// Files don't have a specific source location
	return nil
}

// Unwrap returns the underlying protoreflect.Descriptor.
func (fd *FileDescriptor) Unwrap() protoreflect.Descriptor {
	return fd.wrapped
}

// UnwrapFile returns the underlying protoreflect.FileDescriptor.
func (fd *FileDescriptor) UnwrapFile() protoreflect.FileDescriptor {
	return fd.wrapped
}

// AsFileDescriptorProto returns the FileDescriptorProto.
func (fd *FileDescriptor) AsFileDescriptorProto() *descriptorpb.FileDescriptorProto {
	return fd.proto
}

// AsProto returns the FileDescriptorProto as a proto.Message.
func (fd *FileDescriptor) AsProto() proto.Message {
	return fd.proto
}

// GetFileOptions returns the file options.
func (fd *FileDescriptor) GetFileOptions() *descriptorpb.FileOptions {
	return fd.proto.Options
}

// GetDependencies returns the file dependencies.
func (fd *FileDescriptor) GetDependencies() []*FileDescriptor {
	return fd.deps
}

// GetMessageTypes returns all top-level messages in the file.
func (fd *FileDescriptor) GetMessageTypes() []*MessageDescriptor {
	msgs := fd.wrapped.Messages()
	result := make([]*MessageDescriptor, msgs.Len())
	for i := 0; i < msgs.Len(); i++ {
		result[i] = &MessageDescriptor{
			wrapped: msgs.Get(i),
			file:    fd,
		}
	}
	return result
}

// GetEnumTypes returns all top-level enums in the file.
func (fd *FileDescriptor) GetEnumTypes() []*EnumDescriptor {
	enums := fd.wrapped.Enums()
	result := make([]*EnumDescriptor, enums.Len())
	for i := 0; i < enums.Len(); i++ {
		result[i] = &EnumDescriptor{
			wrapped: enums.Get(i),
			file:    fd,
		}
	}
	return result
}

// GetServices returns all services in the file.
func (fd *FileDescriptor) GetServices() []*ServiceDescriptor {
	services := fd.wrapped.Services()
	result := make([]*ServiceDescriptor, services.Len())
	for i := 0; i < services.Len(); i++ {
		result[i] = &ServiceDescriptor{
			wrapped: services.Get(i),
			file:    fd,
		}
	}
	return result
}

// IsProto3 returns true if the file uses proto3 syntax.
func (fd *FileDescriptor) IsProto3() bool {
	return fd.proto.GetSyntax() == "proto3" || fd.proto.GetSyntax() == ""
}

// Edition returns the edition of the proto file.
func (fd *FileDescriptor) Edition() descriptorpb.Edition {
	if fd.proto.Edition == nil {
		return descriptorpb.Edition_EDITION_UNKNOWN
	}
	return *fd.proto.Edition
}

// FindMessage finds a message by name.
func (fd *FileDescriptor) FindMessage(msgName string) *MessageDescriptor {
	name := protoreflect.FullName(msgName)
	if !strings.Contains(msgName, ".") && fd.GetPackage() != "" {
		name = protoreflect.FullName(fd.GetPackage() + "." + msgName)
	}

	// Try finding in top-level messages
	msgs := fd.wrapped.Messages()
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		if string(msg.Name()) == msgName || string(msg.FullName()) == string(name) {
			return &MessageDescriptor{wrapped: msg, file: fd}
		}
	}

	// Try finding in nested messages
	for _, m := range fd.GetMessageTypes() {
		if result := findMessageInMessage(m, msgName); result != nil {
			return result
		}
	}

	return nil
}

func findMessageInMessage(m *MessageDescriptor, name string) *MessageDescriptor {
	if m.GetName() == name || string(m.wrapped.FullName()) == name {
		return m
	}
	for _, nested := range m.GetNestedMessageTypes() {
		if result := findMessageInMessage(nested, name); result != nil {
			return result
		}
	}
	return nil
}
