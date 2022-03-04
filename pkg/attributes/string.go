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

// String is a string that can be explicitly set to be empty.
type String struct {
	value    string
	explicit bool
}

// NewString returns a String with the provided value set explicitly.
func NewString(v string) String {
	newString := String{}
	newString.Set(v)

	return newString
}

// Set explicitly sets the value of a String.
func (s *String) Set(v string) {
	s.value = v
	s.explicit = true
}

// Bool returns a boolean value for the String, attempting to parse it from boolean-like values such as "true", "false",
// "1", "0". If the String value can not be parsed, or is not explicitly set, ok will be false.
func (c String) Bool() (value bool, ok bool) {
	if !c.Ok() {
		return
	}

	value, err := strconv.ParseBool(c.value)
	ok = err == nil

	return
}

// UnmarshalJSON implements custom JSON unmarshaling of a plain string
// value into an explicit.String.
func (e *String) UnmarshalJSON(data []byte) error {
	value := new(string)

	if err := json.Unmarshal(data, value); err != nil {
		return err
	}

	e.value = *value
	e.explicit = true

	return nil
}

// EncodeValues implements custom url.Values encoding of an explicit.String
// into url.Values. A String with a full zero-value (Value="", Explicit=false)
// will not add a value to url.Values, but a non-zero value for either Value
// or Explicit will result in Value being added to url.Values.
func (e String) EncodeValues(key string, v *url.Values) error {
	if e.Ok() {
		v.Add(key, e.String())
	}

	return nil
}

// String returns a string representation of a String, which is just its
// value.
func (e String) String() string {
	return e.value
}

// Ok returns a boolean indicating if String has a non-empty value or if
// it was explicitly set.
func (e String) Ok() bool {
	return e.value != "" || e.explicit
}

// Value returns a String's value.
func (e String) Value() string {
	return e.value
}

// ValueOk returns a String's value and a boolean indicating if it has a
// non-empty value or if it was explicitly set.
func (e String) ValueOk() (string, bool) {
	return e.Value(), e.Ok()
}
