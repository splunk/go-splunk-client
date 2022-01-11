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

package immutable

// Name is a pseudo-immutable means of representing a Title. It aims for
// immutability because an entry's title can't be changed.
// Immutability is accomplished, from the perspective of external code using
// this module, due to this type being defined under internal/.
type Name struct {
	Value string `json:"name"`
}

// Title returns the Name's Value.
func (t Name) Title() string {
	return t.Value
}
