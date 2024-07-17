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

package aep0132

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
)

// List methods should have a parent variable if the request has a parent field.
var httpURIParent = &lint.MethodRule{
	Name:     lint.NewRuleName(132, "http-uri-parent"),
	RuleType: lint.NewRuleType(lint.MustRule),
	OnlyIf: func(m *desc.MethodDescriptor) bool {
		return utils.IsListMethod(m) && m.GetInputType().FindFieldByName("parent") != nil
	},
	LintMethod: utils.LintHTTPURIHasParentVariable,
}
