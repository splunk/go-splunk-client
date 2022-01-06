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

package client

// title represents the title of a Splunk object. It is abstracted into a struct
// to enable immutability from the perspective of external code.
type title struct {
	value string
}

// UnmarshalJSON implements custom unmarshaling of content into a Title.
func (t *title) UnmarshalJSON(data []byte) error {
	t.value = string(data)

	return nil
}
