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
	"reflect"
	"testing"
)

func TestParameters_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        Parameters(nil),
		},
		{
			name:        "invalid type (list)",
			inputString: `{"param":[]}`,
			want:        Parameters(nil),
			wantError:   true,
		},
		{
			name:        "invalid type (dict)",
			inputString: `{"param":{}}`,
			want:        Parameters(nil),
			wantError:   true,
		},
		{
			name: "valid",
			inputString: `{
				"stringField": "string value",
				"boolField":   true,
				"intField":    1,
				"floatField":  1.234
			}`,
			want: Parameters{
				"stringField": "string value",
				"boolField":   "true",
				"intField":    "1",
				"floatField":  "1.234",
			},
		},
	}

	tests.test(t)
}

func Test_dottedParameterParts(t *testing.T) {
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
	}

	for _, test := range tests {
		gotName, gotParamName := dottedParameterParts(test.input)

		if (gotName != test.wantName) || (gotParamName != test.wantParamName) {
			t.Errorf("%s dottedParameterParts() got\n(%q, %q), want\n(%q, %q)", test.name, gotName, gotParamName, test.wantName, test.wantParamName)
		}
	}
}

func TestParameters_dottedNames(t *testing.T) {
	tests := []struct {
		name  string
		input Parameters
		want  []string
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"empty",
			Parameters{},
			nil,
		},
		{
			"populated",
			Parameters{
				"fieldA":               "fieldAValue",
				"fieldA.paramA":        "paramAValue",
				"fieldA.paramA.paramB": "paramBValue",
				"fieldB":               "fieldBValue",
				// field0 should get sorted to the top of the list
				"field0": "field",
			},
			[]string{
				"field0",
				"fieldA",
				"fieldB",
			},
		},
	}

	for _, test := range tests {
		got := test.input.dottedNames()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s dottedNames() got\n%#v, want\n%#v", test.name, got, test.want)
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

func TestParameters_namedParametersCollection(t *testing.T) {
	tests := []struct {
		name  string
		input Parameters
		want  NamedParametersCollection
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"empty",
			Parameters{},
			nil,
		},
		{
			"populated",
			Parameters{
				"fieldA":        "fieldValueA",
				"fieldA.paramA": "paramValueA",
				"fieldB":        "fieldValueB",
				"fieldB.paramB": "paramValueB",
				// field0 should get sorted to the top of the Collection
				"field0.param0": "paramValue0",
			},
			NamedParametersCollection{
				{
					Name: "field0",
					Parameters: Parameters{
						"param0": "paramValue0",
					},
				},
				{
					Name:   "fieldA",
					Status: "fieldValueA",
					Parameters: Parameters{
						"paramA": "paramValueA",
					},
				},
				{
					Name:   "fieldB",
					Status: "fieldValueB",
					Parameters: Parameters{
						"paramB": "paramValueB",
					},
				},
			},
		},
	}

	for _, test := range tests {
		got := test.input.namedParametersCollection()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s namedParametersCollection() got\n%#v, want\n%#v", test.name, got, test.want)
		}
	}
}

type hasParameters struct {
	StringValue string     `json:"stringField"`
	Params      Parameters `parameters:"params"`
	MoreParams  Parameters `parameters:"moreParams"`
}

func (has *hasParameters) UnmarshalJSON(data []byte) error {
	type hasParametersAlias hasParameters
	var alias hasParametersAlias

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if err := UnmarshalJSONForParameters(data, &alias); err != nil {
		return err
	}

	*has = hasParameters(alias)

	return nil
}

func TestHasParameters_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        hasParameters{},
		},
		{
			name: "valid",
			inputString: `{
				"stringField": "stringValue",
				"params.fieldA": "valueA",
				"moreParams.anotherFieldA": "anotherValueA"
			}`,
			want: hasParameters{
				StringValue: "stringValue",
				Params:      Parameters{"fieldA": "valueA"},
				MoreParams:  Parameters{"anotherFieldA": "anotherValueA"},
			},
		},
	}

	tests.test(t)
}
