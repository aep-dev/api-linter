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

// ServiceDescriptor methods

// GetName returns the service name.
func (sd *ServiceDescriptor) GetName() string {
	return string(sd.wrapped.Name())
}

// GetFullyQualifiedName returns the fully qualified service name.
func (sd *ServiceDescriptor) GetFullyQualifiedName() string {
	return string(sd.wrapped.FullName())
}

// GetFile returns the file containing this service.
func (sd *ServiceDescriptor) GetFile() *FileDescriptor {
	return sd.file
}

// GetParent returns the parent (always the file for services).
func (sd *ServiceDescriptor) GetParent() Descriptor {
	return sd.file
}

// GetSourceInfo returns the source code info for this service.
func (sd *ServiceDescriptor) GetSourceInfo() *descriptorpb.SourceCodeInfo_Location {
	path := sd.wrapped.ParentFile().SourceLocations().ByDescriptor(sd.wrapped)
	if path.Path == nil {
		return nil
	}
	return sd.file.srcInfo[pathToString(pathFromSourcePath(path.Path))]
}

// Unwrap returns the underlying protoreflect.Descriptor.
func (sd *ServiceDescriptor) Unwrap() protoreflect.Descriptor {
	return sd.wrapped
}

// GetServiceOptions returns the service options.
func (sd *ServiceDescriptor) GetServiceOptions() *descriptorpb.ServiceOptions {
	return sd.wrapped.Options().(*descriptorpb.ServiceOptions)
}

// GetMethods returns all methods in the service.
func (sd *ServiceDescriptor) GetMethods() []*MethodDescriptor {
	methods := sd.wrapped.Methods()
	result := make([]*MethodDescriptor, methods.Len())
	for i := 0; i < methods.Len(); i++ {
		result[i] = &MethodDescriptor{
			wrapped: methods.Get(i),
			file:    sd.file,
		}
	}
	return result
}

// MethodDescriptor methods

// GetName returns the method name.
func (md *MethodDescriptor) GetName() string {
	return string(md.wrapped.Name())
}

// GetFullyQualifiedName returns the fully qualified method name.
func (md *MethodDescriptor) GetFullyQualifiedName() string {
	return string(md.wrapped.FullName())
}

// GetFile returns the file containing this method.
func (md *MethodDescriptor) GetFile() *FileDescriptor {
	return md.file
}

// GetParent returns the parent service.
func (md *MethodDescriptor) GetParent() Descriptor {
	parent := md.wrapped.Parent()
	if parent == nil {
		return nil
	}

	if sd, ok := parent.(protoreflect.ServiceDescriptor); ok {
		return &ServiceDescriptor{wrapped: sd, file: md.file}
	}
	return nil
}

// GetSourceInfo returns the source code info for this method.
func (md *MethodDescriptor) GetSourceInfo() *descriptorpb.SourceCodeInfo_Location {
	path := md.wrapped.ParentFile().SourceLocations().ByDescriptor(md.wrapped)
	if path.Path == nil {
		return nil
	}
	return md.file.srcInfo[pathToString(pathFromSourcePath(path.Path))]
}

// Unwrap returns the underlying protoreflect.Descriptor.
func (md *MethodDescriptor) Unwrap() protoreflect.Descriptor {
	return md.wrapped
}

// GetMethodOptions returns the method options.
func (md *MethodDescriptor) GetMethodOptions() *descriptorpb.MethodOptions {
	return md.wrapped.Options().(*descriptorpb.MethodOptions)
}

// GetInputType returns the input message type.
func (md *MethodDescriptor) GetInputType() *MessageDescriptor {
	input := md.wrapped.Input()
	if input == nil {
		return nil
	}

	return &MessageDescriptor{
		wrapped: input,
		file:    md.file,
	}
}

// GetOutputType returns the output message type.
func (md *MethodDescriptor) GetOutputType() *MessageDescriptor {
	output := md.wrapped.Output()
	if output == nil {
		return nil
	}

	return &MessageDescriptor{
		wrapped: output,
		file:    md.file,
	}
}

// IsClientStreaming returns true if this is a client streaming method.
func (md *MethodDescriptor) IsClientStreaming() bool {
	return md.wrapped.IsStreamingClient()
}

// IsServerStreaming returns true if this is a server streaming method.
func (md *MethodDescriptor) IsServerStreaming() bool {
	return md.wrapped.IsStreamingServer()
}
