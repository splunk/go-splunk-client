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
	"fmt"
	"net/url"
)

// Int is an integer that can be explicity set to the zero value.
type Int struct {
	value    int
	explicit bool
}

// NewInt returns an Int with the provided value set explicitly.
func NewInt(v int) Int {
	newInt := Int{}
	newInt.Set(v)

	return newInt
}

// Set explicitly sets the value.
func (i *Int) Set(v int) {
	i.value = v
	i.explicit = true
}

// UnmarshalJSON implements custom JSON unmarshaling of a plain integer
// value into an attributes.Int.
func (i *Int) UnmarshalJSON(data []byte) error {
	value := new(int)

	if err := json.Unmarshal(data, value); err != nil {
		return err
	}

	i.value = *value
	i.explicit = true

	return nil
}

// EncodeValues implements custom url.Values encoding of an attributes.Int
// into url.Values. An Int with a full zero-value (value=0, explicit=false)
// will not add a value to url.Values, but a non-zero value for either value
// or explicit will result in value being added to url.Values.
func (i Int) EncodeValues(key string, v *url.Values) error {
	if i.Ok() {
		v.Add(key, i.String())
	}

	return nil
}

// String returns a string representation of Int.
func (i Int) String() string {
	return fmt.Sprintf("%d", i.value)
}

// Ok returns a boolean indicating if Int has a non-zero value or was
// explicitly set.
func (i Int) Ok() bool {
	return (i.value != 0 || i.explicit)
}

// Value returns an Int's value.
func (i Int) Value() int {
	return i.value
}

// ValueOk returns an Int's value, and a boolean indicating if it has a
// non-zero value or was explicitly set.
func (i Int) ValueOk() (int, bool) {
	return i.Value(), i.Ok()
}
