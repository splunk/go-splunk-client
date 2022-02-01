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
	"strings"
)

// Strings is a slice of strings that can be explicitly set to empty.
type Strings struct {
	values   []string
	explicit bool
}

// Set sets the values of Strings.
func (s *Strings) Set(values ...string) {
	// create new slice to avoid slice pointer issues
	s.values = make([]string, len(values))
	for i, value := range values {
		s.values[i] = value
	}
	s.explicit = true
}

// Add adds new values to the Strings.
func (s *Strings) Add(values ...string) {
	s.values = append(s.values, values...)
}

// NewStrings returns a new Strings with the given values.
func NewStrings(values ...string) Strings {
	newStrings := Strings{}
	newStrings.Set(values...)

	return newStrings
}

// UnmarshalJSON implements custom JSON unmarshaling of a plain list of strings
// value into an attributes.Strings.
func (s *Strings) UnmarshalJSON(data []byte) error {
	values := new([]string)

	if err := json.Unmarshal(data, values); err != nil {
		return err
	}

	s.values = *values
	s.explicit = true

	return nil
}

// EncodeValues implements custom url.Values encoding of a Strings object. If
// Strings has been explicitly set empty, a single value of an empty string
// will be added to url.Values to permit clearing of previously set REST values.
func (s Strings) EncodeValues(key string, v *url.Values) error {
	// add a single empty value if explicitly
	if s.explicit && len(s.values) == 0 {
		v.Add(key, "")
	}

	// the empty value above will still be the only value if values is empty
	for _, value := range s.values {
		v.Add(key, value)
	}

	return nil
}

// String returns a string representation of Strings.
func (s Strings) String() string {
	return strings.Join(s.values, ", ")
}

// Ok returns true if Strings has at least one value, or if it
// has been explicitly set.
func (s Strings) Ok() bool {
	return (len(s.values) != 0 || s.explicit)
}

// Values returns a slice of strings containing the values set.
func (s Strings) Values() []string {
	// don't bother making a new slice if values is nil
	if s.values == nil {
		return nil
	}

	// make a new slice to prevent slice pointer issues
	values := make([]string, len(s.values))
	for i, value := range s.values {
		values[i] = value
	}

	return values
}

// ValuesOk returns the Strings' values, and a boolean indicating
// if Strings has at least one value, or if it has been explicitly
// set.
func (s Strings) ValuesOk() ([]string, bool) {
	return s.Values(), s.Ok()
}
