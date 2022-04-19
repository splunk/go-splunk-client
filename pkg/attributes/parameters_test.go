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

	"github.com/splunk/go-splunk-client/pkg/internal/checks"
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
				Status: NewExplicit("testname value"),
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

func TestParameters_dottedParameterNameParts(t *testing.T) {
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
					Status: NewExplicit("fieldValueA"),
					Parameters: Parameters{
						"paramA": "paramValueA",
					},
				},
				{
					Name:   "fieldB",
					Status: NewExplicit("fieldValueB"),
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

func TestParameters_UnmarshalJSON(t *testing.T) {
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "empty",
			InputString: `{}`,
			Want:        Parameters(nil),
		},
		{
			Name:        "ignored type (null)",
			InputString: `{"param":null}`,
			Want:        Parameters(nil),
		},
		{
			Name:        "invalid type (list)",
			InputString: `{"param":[]}`,
			Want:        Parameters(nil),
			WantError:   true,
		},
		{
			Name:        "invalid type (dict)",
			InputString: `{"param":{}}`,
			Want:        Parameters(nil),
			WantError:   true,
		},
		{
			Name: "valid",
			InputString: `{
				"stringField": "string value",
				"boolField":   true,
				"intField":    1,
				"floatField":  1.234
			}`,
			Want: Parameters{
				"stringField": "string value",
				"boolField":   "true",
				"intField":    "1",
				"floatField":  "1.234",
			},
		},
	}

	tests.Test(t)
}

// testTypeWithParameters is a type used to test custom unmarshaling of Parameters fields.
type testTypeWithParameters struct {
	Name     string     `json:"name"`
	Args     Parameters `parameters:"args"`
	Dispatch Parameters `parameters:"dispatch"`
}

// UnmarshalJSON implements custom unmarshaling for the test type.
func (valueWithParameters *testTypeWithParameters) UnmarshalJSON(data []byte) error {
	// to permit unmarshaling of the non-Parameter fields as normal, we have to create a new
	// type identical to the type we're actually unmarshaling, as this new type won't also
	// have the UnmarshalJSON override method. without this new type attempting to unmarshal
	// the rest of the type would result in infinite recursion.
	//
	// this can probably be handled directly in UnmarshalJSONForParameters with generics once
	// go 1.18 is released.
	type hasParametersAlias testTypeWithParameters
	var aliasedValueWithParameters hasParametersAlias

	// first unmarshal the data into the aliased type, to get the standard unmarshaling treatment.
	if err := json.Unmarshal(data, &aliasedValueWithParameters); err != nil {
		return err
	}

	// then unmarshal the data into the aliased type using the custom ForParameters method.
	if err := UnmarshalJSONForParameters(data, &aliasedValueWithParameters); err != nil {
		return err
	}

	*valueWithParameters = testTypeWithParameters(aliasedValueWithParameters)

	return nil
}

func TestHasParameters_UnmarshalJSON(t *testing.T) {
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "empty",
			InputString: `{}`,
			Want:        testTypeWithParameters{},
		},
		{
			Name: "valid",
			InputString: `{
				"name": "Test Name",
				"args.argA": "argValueA",
				"dispatch.dispatchA": "dispatchValueA"
			}`,
			Want: testTypeWithParameters{
				Name:     "Test Name",
				Args:     Parameters{"argA": "argValueA"},
				Dispatch: Parameters{"dispatchA": "dispatchValueA"},
			},
		},
	}

	tests.Test(t)
}

func TestParameters_SetURLValues(t *testing.T) {
	type testType struct {
		Name   string     `values:"name,omitzero"`
		Params Parameters `values:"params"`
	}

	tests := checks.QueryValuesTestCases{
		{
			Name:  "empty",
			Input: testType{},
			Want:  map[string][]string{},
		},
		{
			Name: "no param values",
			Input: testType{
				Name: "testName",
			},
			Want: map[string][]string{
				"name": {"testName"},
			},
		},
		{
			Name: "param values",
			Input: testType{
				Name: "testName",
				Params: Parameters{
					"fieldA":        "valueA",
					"fieldA.fieldB": "valueB",
				},
			},
			Want: map[string][]string{
				"name":                 {"testName"},
				"params.fieldA":        {"valueA"},
				"params.fieldA.fieldB": {"valueB"},
			},
		},
	}

	tests.Test(t)
}
