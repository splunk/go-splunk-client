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

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Bool is a boolean value that can be explicitly set to false.
type Bool struct {
	value    bool
	explicit bool
}

// Set explicitly sets a Bool's value.
func (b *Bool) Set(v bool) {
	b.value = v
	b.explicit = true
}

// NewBool returns a new Bool with an explicitly set value.
func NewBool(v bool) Bool {
	newBool := Bool{}
	newBool.Set(v)

	return newBool
}

// UnmarshalJSON implements custom JSON unmarshaling of a plain bool
// value into an attributes.Bool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	value := new(bool)

	if err := json.Unmarshal(data, value); err != nil {
		return err
	}

	b.value = *value
	b.explicit = true

	return nil
}

// EncodeValues implements custom url.Values encoding of an attributes.Bool
// into url.Values. A Bool with a full zero-value (Value=false, Explicit=false)
// will not add a value to url.Values, but a non-zero value for either Value
// or Explicit will result in Value being added to url.Values.
func (b Bool) EncodeValues(key string, v *url.Values) error {
	if b.Ok() {
		v.Add(key, b.String())
	}

	return nil
}

// String returns a string representation of Bool.
func (b Bool) String() string {
	return strconv.FormatBool(b.value)
}

// Ok returns a boolean indicating if Bool is true or if it was explicitly
// set false.
func (b Bool) Ok() bool {
	return b.value || b.explicit
}

// Value returns a Bool's value.
func (b Bool) Value() bool {
	return b.value
}

// ValueOk returns a Bool's value and a boolean indicating if it is true or
// was explicitly set false.
func (b Bool) ValueOk() (bool, bool) {
	return b.Value(), b.Ok()
}
