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

// ClearableListValues returns values directly, if it has a length greater than
// zero, otherwise returns a list with a single value of an empty string. This
// functionality provides what the Splunk REST API needs to clear list values,
// as otherwise the lack of any provided values would result in a no-op against
// that field.
func ClearableListValues(values []string) []string {
	if len(values) == 0 {
		return []string{""}
	}

	return values
}
