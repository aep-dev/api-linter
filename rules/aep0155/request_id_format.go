// Copyright 2023 Google LLC
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

package aep0155

import (
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules/internal/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

var requestIdFormat = &lint.FieldRule{
	Name: lint.NewRuleName(155, "request-id-format"),
	OnlyIf: func(fd *desc.FieldDescriptor) bool {
		return fd.GetName() == "request_id" && 
		fd.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE;
	},
	LintField: func(fd *desc.FieldDescriptor) []lint.Problem {
		if(fd.GetMessageType().GetFullyQualifiedName() != "aep.api.IdempotencyKey") {
			return []lint.Problem{{
				Message:    "The `request_id` field should have type `aep.api.IdempotencyKey`",
				Descriptor: fd,
			}}
	}

		return nil
	},
}
