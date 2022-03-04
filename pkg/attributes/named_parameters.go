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

import "strconv"

// NamedParameters represent a set of Parameters that are associated with an overall Name.
type NamedParameters struct {
	// Name is the overall name of this set of Parameters. It is likely the leftmost segment
	// of a dotted parameter name, such as "actions" for "actions.email".
	Name string

	// Status is the string representation of a NamedParameters' status. This is typically
	// true/false or 0/1, and is the value associated directly with the name segment, such as
	// email=true.
	Status string

	Parameters Parameters
}

// StatusBool returns a boolean value for the Status field, attempting to parse it from boolean-like
// values such as "true", "false", "1", "0". If the Status value can not be parsed, ok will be false.
func (params NamedParameters) StatusBool() (value bool, ok bool) {
	value, err := strconv.ParseBool(params.Status)
	ok = err == nil

	return
}
