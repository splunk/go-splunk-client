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

func TestNamedParametersCollection_Enabled(t *testing.T) {
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

type testNamedParametersCollection struct {
	Name    string                    `json:"name"`
	Options NamedParametersCollection `npc:"options"`
}

func (collection *testNamedParametersCollection) UnmarshalJSON(data []byte) error {
	type testNamedParametersCollectionAlias testNamedParametersCollection
	var alias testNamedParametersCollectionAlias

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if err := UnmarshalJSONForNamedParametersCollection(data, &alias); err != nil {
		return err
	}

	*collection = testNamedParametersCollection(alias)

	return nil
}

func TestNamedParametersCollection_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        testNamedParametersCollection{},
			wantError:   false,
		},
		{
			name: "working",
			inputString: `{
				"name":"working",
				"options.disabledOption.description":"this option is not enabled",
				"options.enabledOption":"true",
				"options.enabledOption.description":"this option is enabled"
			}`,
			want: testNamedParametersCollection{
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
						Status: "true",
						Parameters: Parameters{
							"description": "this option is enabled",
						},
					},
				},
			},
			wantError: false,
		},
	}

	tests.test(t)
}
