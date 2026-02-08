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

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aep-dev/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/proto"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

func (c *cli) runV1(rulesV1 lint.RuleRegistry, configs lint.Configs) ([]lint.Response, error) {
	// Prepare proto import lookup.
	fs, err := loadFileDescriptors(c.ProtoDescPath...)
	if err != nil {
		return nil, err
	}
	lookupImport := func(name string) (*desc.FileDescriptor, error) {
		if f, found := fs[name]; found {
			return f, nil
		}
		return nil, fmt.Errorf("%q is not found", name)
	}
	var errorsWithPos []protoparse.ErrorWithPos
	var lock sync.Mutex
	// Parse proto files into `protoreflect` file descriptors.
	p := protoparse.Parser{
		ImportPaths:           c.ProtoImportPaths,
		IncludeSourceCodeInfo: true,
		LookupImport:          lookupImport,
		ErrorReporter: func(errorWithPos protoparse.ErrorWithPos) error {
			// Protoparse isn't concurrent right now but just to be safe for the future.
			lock.Lock()
			errorsWithPos = append(errorsWithPos, errorWithPos)
			lock.Unlock()
			// Continue parsing. The error returned will be protoparse.ErrInvalidSource.
			return nil
		},
	}
	// Resolve file absolute paths to relative ones.
	protoFiles, err := protoparse.ResolveFilenames(c.ProtoImportPaths, c.ProtoFiles...)
	if err != nil {
		return nil, err
	}
	fd, err := p.ParseFiles(protoFiles...)
	if err != nil {
		if err == protoparse.ErrInvalidSource {
			if len(errorsWithPos) == 0 {
				return nil, errors.New("got protoparse.ErrInvalidSource but no ErrorWithPos errors")
			}
			errStrings := make([]string, len(errorsWithPos))
			for i, errorWithPos := range errorsWithPos {
				errStrings[i] = errorWithPos.Error()
			}
			return nil, errors.New(strings.Join(errStrings, "\n"))
		}
		return nil, err
	}

	// Create a linter to lint the file descriptors.
	l := lint.New(rulesV1, configs, lint.Debug(c.DebugFlag), lint.IgnoreCommentDisables(c.IgnoreCommentDisablesFlag))
	return l.LintProtos(fd...)
}

func loadFileDescriptors(filePaths ...string) (map[string]*desc.FileDescriptor, error) {
	fds := []*dpb.FileDescriptorProto{}
	for _, filePath := range filePaths {
		fs, err := readFileDescriptorSet(filePath)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fs.GetFile()...)
	}
	return desc.CreateFileDescriptors(fds)
}

func readFileDescriptorSet(filePath string) (*dpb.FileDescriptorSet, error) {
	in, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fs := &dpb.FileDescriptorSet{}
	if err := proto.Unmarshal(in, fs); err != nil {
		return nil, err
	}
	return fs, nil
}
