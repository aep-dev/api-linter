// Copyright 2019 Google LLC
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
	"strings"

	aepapi "buf.build/gen/go/aep/api/protocolbuffers/go/aep/api"
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/locations"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	apb "google.golang.org/genproto/googleapis/api/annotations"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

var resourceVariables = &lint.MessageRule{
	Name:     lint.NewRuleName(4, "resource-variables"),
	OnlyIf:   hasResourceAnnotation,
	RuleType: lint.NewRuleType(lint.MustRule),
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)

		return lintResourceVariables(resource, m, locations.MessageResource(m))
	},
}

// lintResourceVariables lints the resource ID segments of the pattern(s) in the
// give ResourceDescriptor. This is used for the message-level annotation
// aep.api.resource.
func lintResourceVariables(resource *aepapi.ResourceDescriptor, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	return lintResourceVariablesCommon(resource.GetPattern(), desc, loc)
}

// lintResourceVariablesGoogleAPI lints resource variables for Google API ResourceDescriptor (used for file-level resource definitions)
func lintResourceVariablesGoogleAPI(resource *apb.ResourceDescriptor, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	return lintResourceVariablesCommon(resource.GetPattern(), desc, loc)
}

func lintResourceVariablesCommon(patterns []string, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	for _, pattern := range patterns {
		for _, variable := range getVariables(pattern) {
			if strings.ToLower(variable) != variable {
				return []lint.Problem{{
					Message: fmt.Sprintf(
						"Variable names in patterns should use snake case, such as %q.",
						getDesiredPattern(pattern),
					),
					Descriptor: desc,
					Location:   loc,
				}}
			}
		}
	}
	return nil
}
