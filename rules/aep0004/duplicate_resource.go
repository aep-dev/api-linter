// Copyright 2021 Google LLC
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

package aep0004

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

type resourceDef struct {
	desc *desc.MessageDescriptor
}

func (d *resourceDef) String() string {
	return fmt.Sprintf("message `%s`", d.desc.GetFullyQualifiedName())
}

func (d *resourceDef) location() *dpb.SourceCodeInfo_Location {
	return locations.MessageResource(d.desc)
}

func resourceDefsInFile(f *desc.FileDescriptor, defs map[string][]resourceDef) map[string][]resourceDef {
	for _, m := range f.GetMessageTypes() {
		resourceDefsInMsg(m, defs)
	}
	return defs
}

func resourceDefsInMsg(m *desc.MessageDescriptor, defs map[string][]resourceDef) {
	if t := utils.GetResource(m).GetType(); t != "" {
		defs[t] = append(defs[t], resourceDef{m})
	}
	for _, m := range m.GetNestedMessageTypes() {
		resourceDefsInMsg(m, defs)
	}
}

func allDeps(f *desc.FileDescriptor, deps map[string]*desc.FileDescriptor) map[string]*desc.FileDescriptor {
	for _, f := range f.GetDependencies() {
		name := f.GetName()
		if _, ok := deps[name]; !ok {
			deps[name] = f
			allDeps(f, deps)
		}
	}
	return deps
}

var duplicateResource = &lint.FileRule{
	Name:     lint.NewRuleName(4, "duplicate-resource"),
	RuleType: lint.NewRuleType(lint.MustRule),
	LintFile: func(f *desc.FileDescriptor) []lint.Problem {
		defsInFile := resourceDefsInFile(f, map[string][]resourceDef{})
		if len(defsInFile) == 0 {
			return nil
		}

		defsInDeps := map[string][]resourceDef{}
		for _, f := range allDeps(f, map[string]*desc.FileDescriptor{}) {
			resourceDefsInFile(f, defsInDeps)
		}

		var resourceTypes []string
		for t := range defsInFile {
			resourceTypes = append(resourceTypes, t)
		}
		sort.Strings(resourceTypes)

		ps := []lint.Problem{}
		for _, t := range resourceTypes {
			ds := defsInFile[t]
			locs := []string{}
			// Collect all definitions from current file
			for _, d := range ds {
				locs = append(locs, d.String())
			}
			// Collect all definitions from dependencies
			for _, d := range defsInDeps[t] {
				locs = append(locs, d.String())
			}
			// Only report if there are duplicates (2 or more total definitions)
			if len(locs) <= 1 {
				continue
			}
			sort.Strings(locs)
			msg := fmt.Sprintf("Multiple definitions for resource %q: %s.", t, strings.Join(locs, ", "))
			for _, d := range ds {
				ps = append(ps, lint.Problem{
					Message:    msg,
					Descriptor: d.desc,
					Location:   d.location(),
				})
			}
		}
		return ps
	},
}
