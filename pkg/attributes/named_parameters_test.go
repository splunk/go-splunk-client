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
	"reflect"
	"testing"
)

func TestNamedParameters_StatusBool(t *testing.T) {
	tests := []struct {
		input     string
		wantValue bool
		wantOk    bool
	}{
		{
			"",
			false,
			false,
		},
		{
			"nonsense",
			false,
			false,
		},
		{
			"-1",
			false,
			false,
		},
		{
			"0",
			false,
			true,
		},
		{
			"1",
			true,
			true,
		},
		{
			"f",
			false,
			true,
		},
		{
			"false",
			false,
			true,
		},
		{
			"t",
			true,
			true,
		},
		{
			"true",
			true,
			true,
		},
	}

	for _, test := range tests {
		gotValue, gotOk := NamedParameters{Status: test.input}.StatusBool()

		if (gotValue != test.wantValue) && (gotOk != test.wantOk) {
			t.Errorf("%q StatusBool got\n%#v, want\n%#v", test.input, []bool{gotValue, gotOk}, []bool{test.wantValue, test.wantOk})
		}
	}
}

func TestNamedParametersCollection_EnabledNames(t *testing.T) {
	tests := []struct {
		name  string
		input NamedParametersCollection
		want  []string
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"empty",
			NamedParametersCollection{},
			nil,
		},
		{
			"some enabled",
			NamedParametersCollection{
				{
					Name: "implicitlyDisabledField",
				},
				{
					Name:   "explicitlyDisabledField",
					Status: "false",
				},
				{
					Name:   "explicitlyEnabledBoolField",
					Status: "true",
				},
				{
					Name:   "explicitlyEnabledNumberField",
					Status: "1",
				},
			},
			[]string{
				"explicitlyEnabledBoolField",
				"explicitlyEnabledNumberField",
			},
		},
	}

	for _, test := range tests {
		got := test.input.EnabledNames()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s EnabledNames() got\n%#v, want\n%#v", test.name, got, test.want)
		}
	}
}
