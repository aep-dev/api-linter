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

// Package protoparse provides proto file parsing using buf's protocompile.
package protoparse

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aep-dev/api-linter/internal/desc"
	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/reporter"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ErrorWithPos represents a parse error with position information.
type ErrorWithPos struct {
	Pos reporter.ErrorWithPos
}

func (e ErrorWithPos) Error() string {
	if e.Pos == nil {
		return "unknown error"
	}
	return e.Pos.Error()
}

// ErrInvalidSource is returned when source parsing fails.
var ErrInvalidSource = fmt.Errorf("parse error")

// Parser is a proto file parser using protocompile.
type Parser struct {
	// ImportPaths are the directories to search for imports.
	ImportPaths []string

	// LookupImport is called to resolve imports.
	LookupImport func(string) (*desc.FileDescriptor, error)

	// IncludeSourceCodeInfo determines whether source code info is included.
	IncludeSourceCodeInfo bool

	// Accessor provides file contents for parsing.
	Accessor FileAccessor

	// ErrorReporter handles errors during parsing.
	ErrorReporter func(ErrorWithPos) error
}

// FileAccessor is an interface for accessing file contents.
type FileAccessor interface {
	Open(filename string) (io.ReadCloser, error)
}

// FileContentsFromMap returns a FileAccessor that reads from a map.
func FileContentsFromMap(files map[string]string) FileAccessor {
	return mapAccessor{files: files}
}

type mapAccessor struct {
	files map[string]string
}

func (m mapAccessor) Open(filename string) (io.ReadCloser, error) {
	content, ok := m.files[filename]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return io.NopCloser(strings.NewReader(content)), nil
}

// ParseFiles parses the given proto files.
func (p *Parser) ParseFiles(filenames ...string) ([]*desc.FileDescriptor, error) {
	resolver := protocompile.SourceResolver{
		ImportPaths: p.ImportPaths,
	}

	if p.Accessor != nil {
		adapter := &accessorAdapter{accessor: p.Accessor, importPaths: p.ImportPaths}
		resolver.Accessor = adapter.Open
	}

	compiler := &protocompile.Compiler{
		Resolver:       &resolver,
		SourceInfoMode: protocompile.SourceInfoStandard,
	}

	if !p.IncludeSourceCodeInfo {
		compiler.SourceInfoMode = protocompile.SourceInfoNone
	}

	if p.ErrorReporter != nil {
		compiler.Reporter = &errorReporter{handler: p.ErrorReporter}
	}

	// Compile the files
	ctx := context.Background()
	compiled, err := compiler.Compile(ctx, filenames...)
	if err != nil {
		return nil, err
	}

	// Convert linker.Files to []protoreflect.FileDescriptor
	fds := make([]protoreflect.FileDescriptor, len(compiled))
	for i, file := range compiled {
		fds[i] = file
	}

	return desc.WrapFiles(fds)
}

type errorReporter struct {
	handler func(ErrorWithPos) error
}

func (er *errorReporter) Error(err reporter.ErrorWithPos) error {
	return er.handler(ErrorWithPos{Pos: err})
}

func (er *errorReporter) Warning(reporter.ErrorWithPos) {
	// Ignore warnings
}

type accessorAdapter struct {
	accessor    FileAccessor
	importPaths []string
}

func (a *accessorAdapter) Open(path string) (io.ReadCloser, error) {
	// Try the path as-is first
	if rc, err := a.accessor.Open(path); err == nil {
		return rc, nil
	}

	// Try with import paths
	for _, importPath := range a.importPaths {
		fullPath := filepath.Join(importPath, path)
		if rc, err := a.accessor.Open(fullPath); err == nil {
			return rc, nil
		}
	}

	// Fall back to file system
	return os.Open(path)
}

// ResolveFilenames resolves file paths relative to import paths.
func ResolveFilenames(importPaths []string, filenames ...string) ([]string, error) {
	resolved := make([]string, len(filenames))
	for i, filename := range filenames {
		found := false
		for _, importPath := range importPaths {
			fullPath := filepath.Join(importPath, filename)
			if _, err := os.Stat(fullPath); err == nil {
				resolved[i] = filename
				found = true
				break
			}
		}
		if !found {
			// Check if file exists as absolute path
			if _, err := os.Stat(filename); err == nil {
				resolved[i] = filename
			} else {
				return nil, fs.ErrNotExist
			}
		}
	}
	return resolved, nil
}
