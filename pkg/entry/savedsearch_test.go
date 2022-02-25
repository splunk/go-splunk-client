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

package entry

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestSavedSearchActions_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantActions SavedSearchActions
	}{
		{
			"empty",
			`{}`,
			SavedSearchActions{},
		},
		{
			"actions only",
			`{
				"actions":"one,two"
			 }`,
			SavedSearchActions{
				{Name: "one"},
				{Name: "two"},
			},
		},
		{
			"params only",
			`{
				"action.test_action.string_param":"string_value",
				"action.test_action.bool_param":true,
				"action.test_action.int_param":1,
				"action.test_action.float_param":1.234,
				"action.test_action.dotted.decimal.param":"dotted.decimal.value"
			 }`,
			SavedSearchActions{
				{
					Name:     "test_action",
					Disabled: true,
					Parameters: map[string]string{
						"string_param":         "string_value",
						"bool_param":           "true",
						"int_param":            "1",
						"float_param":          "1.234",
						"dotted.decimal.param": "dotted.decimal.value",
					},
				},
			},
		},
		{
			"enabled with params",
			`{
				"actions":"test_action",
				"action.test_action.test_param":"test_value"
			 }`,
			SavedSearchActions{
				{
					Name:       "test_action",
					Parameters: map[string]string{"test_param": "test_value"},
				},
			},
		},
		{
			"action.name bool",
			`{
				"action.test_action":false
			 }`,
			SavedSearchActions{},
		},
	}

	for _, test := range tests {
		gotActions := SavedSearchActions{}
		if err := json.Unmarshal([]byte(test.input), &gotActions); err != nil {
			t.Fatalf("%s unexpected error: %s", test.name, err)
		}

		if !reflect.DeepEqual(gotActions, test.wantActions) {
			t.Errorf("%s json.Unmarshal got\n%#v, want\n%#v", test.name, gotActions, test.wantActions)
		}
	}
}

func TestSavedSearchActions_EncodeValues(t *testing.T) {
	tests := []struct {
		name       string
		input      SavedSearchActions
		wantValues url.Values
	}{
		{
			"empty",
			SavedSearchActions{},
			url.Values{
				"actions": []string{""},
			},
		},
		{
			"enabled, no params",
			SavedSearchActions{
				{Name: "test_action"},
			},
			url.Values{
				"actions": []string{"test_action"},
			},
		},
		{
			"disabled, no params",
			SavedSearchActions{
				{Name: "test_action", Disabled: true},
			},
			url.Values{
				"actions": []string{""},
			},
		},
		{
			"nil",
			nil,
			url.Values{},
		},
		{
			"enabled, params",
			SavedSearchActions{
				{
					Name: "test_action_1",
					Parameters: map[string]string{
						"test_param_name": "test_param_value",
					},
				},
				{
					Name: "test_action_2",
					Parameters: map[string]string{
						"test_param_name_1": "test_param_value_1",
						"test_param_name_2": "test_param_value_2",
					},
				},
			},
			url.Values{
				"actions":                                []string{"test_action_1,test_action_2"},
				"action.test_action_1.test_param_name":   []string{"test_param_value"},
				"action.test_action_2.test_param_name_1": []string{"test_param_value_1"},
				"action.test_action_2.test_param_name_2": []string{"test_param_value_2"},
			},
		},
		{
			"enabled, explicitly empty params",
			SavedSearchActions{
				{
					Name: "test_action",
					Parameters: map[string]string{
						"cleared_param_name": "",
						"set_param_name":     "set_param_value",
					},
				},
			},
			url.Values{
				"actions":                               []string{"test_action"},
				"action.test_action.cleared_param_name": []string{""},
				"action.test_action.set_param_name":     []string{"set_param_value"},
			},
		},
	}

	type actionsStruct struct {
		Actions SavedSearchActions
	}

	for _, test := range tests {
		gotValues, err := query.Values(actionsStruct{test.input})
		if err != nil {
			t.Fatalf("%s unexpected error: %s", test.name, err)
		}

		if !reflect.DeepEqual(gotValues, test.wantValues) {
			t.Errorf("%s query.Values got\n%#v, want\n%#v", test.name, gotValues, test.wantValues)
		}
	}
}
