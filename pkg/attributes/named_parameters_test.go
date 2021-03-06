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

func TestNamedParameters_SetURLValues(t *testing.T) {
	type testType struct {
		Description string          `values:"description,omitzero"`
		Action      NamedParameters `values:"actions"`
	}

	tests := checks.QueryValuesTestCases{
		{
			Name:      "empty",
			Input:     testType{},
			WantError: true,
		},
		{
			Name: "description only",
			Input: testType{
				Description: "testDescription",
			},
			WantError: true,
		},
		{
			Name: "no status",
			Input: testType{
				Description: "testDescription",
				Action: NamedParameters{
					Name: "email",
					Parameters: Parameters{
						"to":      "whocares@example.com",
						"subject": "10 tricks your Splunk admin doesn't want you to know!",
					},
				},
			},
			Want: map[string][]string{
				"description":           {"testDescription"},
				"actions.email.to":      {"whocares@example.com"},
				"actions.email.subject": {"10 tricks your Splunk admin doesn't want you to know!"},
			},
		},
		{
			Name: "email action",
			Input: testType{
				Description: "testDescription",
				Action: NamedParameters{
					Name:   "email",
					Status: NewExplicit("true"),
					Parameters: Parameters{
						"to":      "whocares@example.com",
						"subject": "10 tricks your Splunk admin doesn't want you to know!",
					},
				},
			},
			Want: map[string][]string{
				"description":           {"testDescription"},
				"actions.email":         {"true"},
				"actions.email.to":      {"whocares@example.com"},
				"actions.email.subject": {"10 tricks your Splunk admin doesn't want you to know!"},
			},
		},
	}

	tests.Test(t)
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
					Status: NewExplicit("false"),
				},
				{
					Name:   "explicitlyEnabledBoolField",
					Status: NewExplicit("true"),
				},
				{
					Name:   "explicitlyEnabledNumberField",
					Status: NewExplicit("1"),
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

// testTypeWithNamedParametersCollection is a type used to test custom unmarshaling of NamedParametersCollection fields..
type testTypeWithNamedParametersCollection struct {
	Name    string                    `json:"name"`
	Options NamedParametersCollection `named_parameters_collection:"options"`
}

// UnmarshalJSON implements custom unmarshaling for the test type.
func (valueWithCollections *testTypeWithNamedParametersCollection) UnmarshalJSON(data []byte) error {
	// to permit unmarshaling of the non-Parameter fields as normal, we have to create a new
	// type identical to the type we're actually unmarshaling, as this new type won't also
	// have the UnmarshalJSON override method. without this new type attempting to unmarshal
	// the rest of the type would result in infinite recursion.
	//
	// this can probably be handled directly in UnmarshalJSONForNamedParametersCollection with
	// generics once go 1.18 is released.
	type testTypeWithNamedParametersCollectionAlias testTypeWithNamedParametersCollection
	var aliasedValueWithCollections testTypeWithNamedParametersCollectionAlias

	if err := json.Unmarshal(data, &aliasedValueWithCollections); err != nil {
		return err
	}

	if err := UnmarshalJSONForNamedParametersCollections(data, &aliasedValueWithCollections); err != nil {
		return err
	}

	*valueWithCollections = testTypeWithNamedParametersCollection(aliasedValueWithCollections)

	return nil
}

func TestNamedParametersCollection_UnmarshalJSON(t *testing.T) {
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "empty",
			InputString: `{}`,
			Want:        testTypeWithNamedParametersCollection{},
			WantError:   false,
		},
		{
			Name: "working",
			InputString: `{
				"name":"working",
				"options.disabledOption.description":"this option is not enabled",
				"options.enabledOption":"true",
				"options.enabledOption.description":"this option is enabled"
			}`,
			Want: testTypeWithNamedParametersCollection{
				Name: "working",
				Options: NamedParametersCollection{
					{
						Name: "disabledOption",
						Parameters: Parameters{
							"description": "this option is not enabled",
						},
					},
					{
						Name:   "enabledOption",
						Status: NewExplicit("true"),
						Parameters: Parameters{
							"description": "this option is enabled",
						},
					},
				},
			},
			WantError: false,
		},
	}

	tests.Test(t)
}

func TestNamedParametersCollection_SetURLValues(t *testing.T) {
	type testType struct {
		Description string                    `values:"description,omitzero"`
		Actions     NamedParametersCollection `named_parameters_collection:"actions" values:"actions,omitzero"`
	}

	tests := checks.QueryValuesTestCases{
		{
			Name:  "empty",
			Input: testType{},
			Want:  map[string][]string{},
		},
		{
			Name: "description only",
			Input: testType{
				Description: "testDescription",
			},
			Want: map[string][]string{
				"description": {"testDescription"},
			},
		},
		{
			Name: "with status",
			Input: testType{
				Description: "testDescription",
				Actions: NamedParametersCollection{
					{
						Name:   "email",
						Status: NewExplicit("false"),
					},
				},
			},
			Want: map[string][]string{
				"description":   {"testDescription"},
				"actions.email": {"false"},
			},
		},
		{
			Name: "with fields",
			Input: testType{
				Description: "testDescription",
				Actions: NamedParametersCollection{
					{
						Name:   "email",
						Status: NewExplicit("true"),
						Parameters: Parameters{
							"to": "whocares@example.com",
						},
					},
				},
			},
			Want: map[string][]string{
				"description":      {"testDescription"},
				"actions.email":    {"true"},
				"actions.email.to": {"whocares@example.com"},
			},
		},
		{
			Name: "with empty fields",
			Input: testType{
				Description: "testDescription",
				Actions: NamedParametersCollection{
					{
						Name: "email",
						Parameters: map[string]string{
							"to": "whocares@example.com",
							// this allows clearing the subject field
							"subject": "",
						},
					},
				},
			},
			Want: map[string][]string{
				"description":           {"testDescription"},
				"actions.email.to":      {"whocares@example.com"},
				"actions.email.subject": {""},
			},
		},
	}

	tests.Test(t)
}
