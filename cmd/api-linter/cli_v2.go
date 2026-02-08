// Copyright 2026 The AEP Authors
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

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"github.com/bufbuild/protocompile/reporter"
	lint_v2 "github.com/aep-dev/api-linter/lint/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

func (c *cli) runV2(rulesV2 lint_v2.RuleRegistry, configs lint_v2.Configs) ([]lint_v2.Response, error) {
	// Create resolver for descriptor sets.
	descResolver, err := loadFileDescriptorsAsResolver(c.ProtoDescPath...)
	if err != nil {
		return nil, err
	}

	// Create resolver for source files.
	imports := resolveImports(c.ProtoImportPaths)
	sourceResolver := &protocompile.SourceResolver{
		ImportPaths: imports,
	}

	resolvers := []protocompile.Resolver{sourceResolver}
	if descResolver != nil {
		resolvers = append(resolvers, descResolver)
	}

	var collectedErrors []error
	rep := reporter.NewReporter(func(err reporter.ErrorWithPos) error {
		collectedErrors = append(collectedErrors, err)
		return nil
	}, nil)

	compiler := protocompile.Compiler{
		Resolver:       protocompile.WithStandardImports(protocompile.CompositeResolver(resolvers)),
		SourceInfoMode: protocompile.SourceInfoExtraOptionLocations,
		Reporter:       rep,
	}

	var compiledFiles linker.Files
	for _, protoFile := range c.ProtoFiles {
		f, err := compiler.Compile(context.Background(), protoFile)
		if len(collectedErrors) > 0 {
			errStrings := make([]string, len(collectedErrors))
			for i, e := range collectedErrors {
				errStrings[i] = e.Error()
			}
			return nil, errors.New(strings.Join(errStrings, "\n"))
		}
		if err != nil {
			return nil, err
		}
		compiledFiles = append(compiledFiles, f...)
	}

	var fileDescriptors []protoreflect.FileDescriptor
	for _, f := range compiledFiles {
		fileDescriptors = append(fileDescriptors, f)
	}

	l := lint_v2.New(rulesV2, configs, lint_v2.Debug(c.DebugFlag), lint_v2.IgnoreCommentDisables(c.IgnoreCommentDisablesFlag))
	return l.LintProtos(fileDescriptors...)
}

func resolveImports(imports []string) []string {
	if len(imports) == 0 {
		return []string{"."}
	}
	return imports
}

type v2Resolver struct {
	files *protoregistry.Files
}

func (r *v2Resolver) FindFileByPath(path string) (protocompile.SearchResult, error) {
	fd, err := r.files.FindFileByPath(path)
	if err != nil {
		return protocompile.SearchResult{}, err
	}
	return protocompile.SearchResult{Desc: fd}, nil
}

func loadFileDescriptorsAsResolver(filePaths ...string) (protocompile.Resolver, error) {
	if len(filePaths) == 0 {
		return nil, nil
	}

	fdsSet := make(map[string]*dpb.FileDescriptorProto)
	for _, filePath := range filePaths {
		in, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		fs := &dpb.FileDescriptorSet{}
		if err := proto.Unmarshal(in, fs); err != nil {
			return nil, err
		}
		for _, fd := range fs.GetFile() {
			if _, exists := fdsSet[fd.GetName()]; !exists {
				fdsSet[fd.GetName()] = fd
			}
		}
	}

	fds := &dpb.FileDescriptorSet{}
	for _, fd := range fdsSet {
		fds.File = append(fds.File, fd)
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return nil, fmt.Errorf("failed to create protoregistry.Files: %w", err)
	}
	return &v2Resolver{files: files}, nil
}
