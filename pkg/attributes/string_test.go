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
	"net/url"
	"testing"
)

type testString struct {
	Value Explicit[string] `values:",omitempty"`
}

func TestString_Bool(t *testing.T) {
	tests := []struct {
		input     Explicit[string]
		wantValue bool
		wantOk    bool
	}{
		{
			Explicit[string]{},
			false,
			false,
		},
		{
			NewExplicit(""),
			false,
			true,
		},
		{
			NewExplicit("nonsense"),
			false,
			false,
		},
		{
			NewExplicit("-1"),
			false,
			false,
		},
		{
			NewExplicit("0"),
			false,
			true,
		},
		{
			NewExplicit("1"),
			true,
			true,
		},
		{
			NewExplicit("f"),
			false,
			true,
		},
		{
			NewExplicit("false"),
			false,
			true,
		},
		{
			NewExplicit("t"),
			true,
			true,
		},
		{
			NewExplicit("true"),
			true,
			true,
		},
	}

	for _, test := range tests {
		gotValue, gotOk := test.input.Bool()

		if (gotValue != test.wantValue) && (gotOk != test.wantOk) {
			t.Errorf("%q StatusBool got\n%#v, want\n%#v", test.input.Value(), []bool{gotValue, gotOk}, []bool{test.wantValue, test.wantOk})
		}
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        testString{},
		},
		{
			name:        "empty",
			inputString: `{"value":""}`,
			want:        testString{Value: NewExplicit("")},
		},
		{
			name:        "non-empty",
			inputString: `{"value":"this string is not empty"}`,
			want:        testString{Value: NewExplicit("this string is not empty")},
		},
	}

	tests.test(t)
}

func TestString_SetURLValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			name:  "implicit empty",
			input: testString{},
			want:  url.Values{},
		},
		{
			name:  "explicit empty",
			input: testString{Value: NewExplicit("")},
			want:  url.Values{"Value": []string{""}},
		},
		{
			name:  "non-empty",
			input: testString{Value: NewExplicit("this string is not empty")},
			want:  url.Values{"Value": []string{"this string is not empty"}},
		},
	}

	tests.test(t)
}
