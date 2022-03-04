// Copyright 2022 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attributes

import "strings"

// dottedParameterNameParts parses a dotted parameter name and returns the first segment as name,
// and the remaining segments as paramName. If there are not multiple segments, name will be empty.
//
// For example:
//
//   dottedParameterNameParts("name")
//   # ("name", "")
//
//   dottedParameterNameParts("actions.email.to")
//   # ("actions", "email.to")
func dottedParameterNameParts(fullFieldName string) (name string, paramName string) {
	parts := strings.Split(fullFieldName, ".")

	return parts[0], strings.Join(parts[1:], ".")
}

// Parameters is a map of parameter names to string values.
type Parameters map[string]string

// withDottedName returns a new Parameters object containing the nested parameters
// for the given name. The new Parameters name field will have this name prefix removed.
//
// For example:
//
//   Parameters{"action.email": "true", "action.email.to": "whoever@example.com"}.withDottedName("action")
//   # Parameters{"email": "true", "email.to": "whoever@example.com"}
func (p Parameters) withDottedName(name string) Parameters {
	var newParameters Parameters

	for key, value := range p {
		fieldName, fieldParamName := dottedParameterNameParts(key)

		if fieldName == name && fieldParamName != "" {
			if newParameters == nil {
				newParameters = Parameters{}
			}

			newParameters[fieldParamName] = value
		}
	}

	return newParameters
}

// namedParametersWithDottedName returns a NamedParameters with the given name and Status and Parameter values
// as calculated from the input Parameters.
//
//   Parameters{"email":"true","email.to":"whoever@example.com"}.namedParametersWithDottedName("email")
//   # NamedParameters{Name: "email", Status: "true", Parameters{"to": "whoever@example.com"}}
func (p Parameters) namedParametersWithDottedName(name string) NamedParameters {
	return NamedParameters{
		Name:       name,
		Status:     p[name],
		Parameters: p.withDottedName(name),
	}
}
