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

// Package aep0133 contains rules defined in https://aep.dev/133.
package aep0133

import (
	"strings"

	"github.com/aep-dev/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
)

// AddRules accepts a register function and registers each of
// this AEP's rules to it.
func AddRules(r lint.RuleRegistry) error {
	return r.Register(
		133,
		httpBody,
		httpURIParent,
		httpURIResource,
		httpMethod,
		inputName,
		methodSignature,
		outputName,
		requestIDField,
		requestParentBehavior,
		requestParentField,
		requestParentRequired,
		requestRequiredFields,
		requestResourceBehavior,
		resourceField,
		resourceReferenceType,
		synonyms,
		unknownFields,
	)
}

// get resource message type name from request message
func getResourceMsgNameFromReq(m *desc.MessageDescriptor) string {
	// retrieve the string between the prefix "Create" and suffix "Request" from
	// the name "Create<XXX>Request", and this part will usually be the resource
	// message name(if its naming follows the right principle)
	resourceMsgName := m.GetName()[6 : len(m.GetName())-7]

	// Get the resource field of the request message if it exist, this part will
	// be exactly the resource message name (make a double check here to avoid the
	// issues when request message naming doesn't follow the right principles)
	for _, fieldDesc := range m.GetFields() {
		if msgDesc := fieldDesc.GetMessageType(); msgDesc != nil && strings.Contains(resourceMsgName, msgDesc.GetName()) {
			resourceMsgName = msgDesc.GetName()
		}
	}

	return resourceMsgName
}
