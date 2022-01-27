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

import "net/url"

// Roles is a list of role names.
type Roles []string

// EncodeValues implements custom encoding into url.Values such that an empty
// set of roles is passed as a single value with an empty string. This is
// necessary to coerce Splunk to clear previously set roles.
func (r Roles) EncodeValues(key string, v *url.Values) error {
	EncodeClearableListValues(v, key, r...)

	return nil
}
