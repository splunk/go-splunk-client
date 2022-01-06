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

package internal

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-querystring/query"
)

// TestURLValuesFunc describes functions that perform tests against the calculated url.Values
// representation of an object.
type TestURLValuesFunc func(name string, input interface{}, t *testing.T)

// ComposeTestURLValues composes a new TestURLValuesFunc from any number of other TestURLValuesFunc items.
func ComposeTestURLValues(tests ...TestURLValuesFunc) TestURLValuesFunc {
	return func(name string, input interface{}, t *testing.T) {
		for _, test := range tests {
			test(name, input, t)
		}
	}
}

// TestURLValuesEquality tests against equality of a given url.Values object.
func TestURLValuesEquality(want url.Values) TestURLValuesFunc {
	return func(name string, input interface{}, t *testing.T) {
		got, err := query.Values(input)
		if err != nil {
			t.Errorf("%s: TestURLValuesEquals failed to compute url.Values: %s", name, err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: TestURLValuesEquals got\n%#v, want\n%#v", name, got, want)
		}
	}
}

// TestURLValuesEncoded tests the Encoded value against a given string.
func TestURLValuesEncoded(want string) TestURLValuesFunc {
	return func(name string, input interface{}, t *testing.T) {
		gotValues, err := query.Values(input)
		if err != nil {
			t.Errorf("%s: TestURLValuesEncoded failed to compute url.Values: %s", name, err)
		}

		got := gotValues.Encode()
		if got != want {
			t.Errorf("%s: TestURLValuesEncoded got\n%s, want\n%s", name, got, want)
		}
	}
}
