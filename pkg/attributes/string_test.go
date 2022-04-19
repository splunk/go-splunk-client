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

	"github.com/splunk/go-splunk-client/pkg/internal/checks"
)

type testString struct {
	Value Explicit[string] `values:",omitzero"`
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
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "empty",
			InputString: `{}`,
			Want:        testString{},
		},
		{
			Name:        "empty",
			InputString: `{"value":""}`,
			Want:        testString{Value: NewExplicit("")},
		},
		{
			Name:        "non-empty",
			InputString: `{"value":"this string is not empty"}`,
			Want:        testString{Value: NewExplicit("this string is not empty")},
		},
	}

	tests.Test(t)
}

func TestString_SetURLValues(t *testing.T) {
	tests := checks.QueryValuesTestCases{
		{
			Name:  "implicit empty",
			Input: testString{},
			Want:  url.Values{},
		},
		{
			Name:  "explicit empty",
			Input: testString{Value: NewExplicit("")},
			Want:  url.Values{"Value": []string{""}},
		},
		{
			Name:  "non-empty",
			Input: testString{Value: NewExplicit("this string is not empty")},
			Want:  url.Values{"Value": []string{"this string is not empty"}},
		},
	}

	tests.Test(t)
}
