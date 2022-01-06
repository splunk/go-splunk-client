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

import "net/url"

// Capabilities is a slice of capability names to be applied to a role.
type Capabilities []string

// EncodeValues implements custom encoding into url.Values such that an empty
// set of capabilities is passed as a single value with an empty string. This is
// necessary to coerce Splunk to clear previously set capabilities.
func (c Capabilities) EncodeValues(key string, v *url.Values) error {
	if len(c) == 0 {
		v.Add(key, "")
	} else {
		for _, capability := range c {
			v.Add(key, capability)
		}
	}

	return nil
}
