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

func Test_dottedParameterNameParts(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantName      string
		wantParamName string
	}{
		{
			"empty",
			"",
			"",
			"",
		},
		{
			"name only",
			"testname",
			"testname",
			"",
		},
		{
			"name and param name",
			"testname.testparam",
			"testname",
			"testparam",
		},
		{
			"name and dotted param name",
			"testname.testparamA.testparamB",
			"testname",
			"testparamA.testparamB",
		},
	}

	for _, test := range tests {
		gotName, gotParamName := dottedParameterNameParts(test.input)

		if (gotName != test.wantName) || (gotParamName != test.wantParamName) {
			t.Errorf("%s Test_dottedParameterNameParts() got\n(%q, %q), want\n(%q, %q)", test.name, gotName, gotParamName, test.wantName, test.wantParamName)
		}
	}
}

func TestParameters_withDottedName(t *testing.T) {
	tests := []struct {
		name        string
		inputParams Parameters
		inputName   string
		wantParams  Parameters
	}{
		{
			"nil",
			Parameters(nil),
			"paramName",
			Parameters(nil),
		},
		{
			"empty",
			Parameters{},
			"paramName",
			Parameters(nil),
		},
		{
			"no nested parameters",
			Parameters{
				"paramName": "paramValue",
			},
			"paramName",
			Parameters(nil),
		},
		{
			"nested parameters",
			Parameters{
				"paramName.partA":       "valueA",
				"paramName.partA.partB": "valueB",
			},
			"paramName",
			Parameters{
				"partA":       "valueA",
				"partA.partB": "valueB",
			},
		},
	}

	for _, test := range tests {
		gotParams := test.inputParams.withDottedName(test.inputName)

		if !reflect.DeepEqual(gotParams, test.wantParams) {
			t.Errorf("%s withDottedName got\n%#v, want\n%#v", test.name, gotParams, test.wantParams)
		}
	}
}

func TestParameters_namedParametersWithDottedName(t *testing.T) {
	tests := []struct {
		name        string
		inputParams Parameters
		inputName   string
		want        NamedParameters
	}{
		{
			"empty",
			Parameters{},
			"testname",
			NamedParameters{Name: "testname"},
		},
		{
			"no matches",
			Parameters{
				"unmatched":       "unmatched value",
				"unmatched.field": "unmatched field value",
			},
			"testname",
			NamedParameters{Name: "testname"},
		},
		{
			"matches",
			Parameters{
				"testname":       "testname value",
				"testname.field": "testname field value",
			},
			"testname",
			NamedParameters{
				Name:   "testname",
				Status: "testname value",
				Parameters: Parameters{
					"field": "testname field value",
				},
			},
		},
	}

	for _, test := range tests {
		got := test.inputParams.namedParametersWithDottedName(test.inputName)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s namedParametersWithDottedName got\n%#v, want\n%#v", test.name, got, test.want)
		}
	}
}
