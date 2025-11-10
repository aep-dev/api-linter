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

var resourcePattern = &lint.MessageRule{
	Name:     lint.NewRuleName(4, "resource-pattern"),
	OnlyIf:   hasResourceAnnotation,
	RuleType: lint.NewRuleType(lint.MustRule),
	LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
		resource := utils.GetResource(m)
		return lintResourcePattern(resource, m, locations.MessageResource(m))
	},
}

func lintResourcePattern(resource *aepapi.ResourceDescriptor, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	return lintResourcePatternCommon(resource.GetPattern(), desc, loc)
}

// lintResourcePatternGoogleAPI lints resource patterns for Google API ResourceDescriptor (used for file-level resource definitions)
func lintResourcePatternGoogleAPI(resource *apb.ResourceDescriptor, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	return lintResourcePatternCommon(resource.GetPattern(), desc, loc)
}

func lintResourcePatternCommon(patterns []string, desc desc.Descriptor, loc *dpb.SourceCodeInfo_Location) []lint.Problem {
	// Are any patterns declared at all? If not, complain.
	if len(patterns) == 0 {
		return []lint.Problem{{
			Message:    "Resources should declare resource name pattern(s).",
			Descriptor: desc,
			Location:   loc,
		}}
	}

	// Ensure that the constant segments of the pattern uses camel case,
	// not snake case, and there are no spaces.
	for _, pattern := range patterns {
		plainPattern := getPlainPattern(pattern)

		if strings.Contains(plainPattern, "_") {
			return []lint.Problem{{
				Message: fmt.Sprintf(
					"Resource patterns should use camel case (apart from the variable names), such as %q.",
					getDesiredPattern(pattern),
				),
				Descriptor: desc,
				Location:   loc,
			}}
		}
		if strings.Contains(plainPattern, " ") {
			return []lint.Problem{{
				Message:    "Resource patterns should not have spaces",
				Descriptor: desc,
				Location:   loc,
			}}
		}
	}
	return nil
}
