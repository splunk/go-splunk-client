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

package values

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

type testMapCustomKey map[string]string

func (m testMapCustomKey) GetURLKey(key1 string, key2 string) (string, error) {
	return fmt.Sprintf("%s[%s]", key1, key2), nil
}

type testCustomEncoder struct {
	key   string
	value string
}

func (e testCustomEncoder) SetURLValues(key string, values *url.Values) error {
	values.Add(e.key, e.value)

	return nil
}

type testCustomKeySlice []string

func (s testCustomKeySlice) GetURLKey(parentKey string, index string) (string, error) {
	return fmt.Sprintf("%s.%s", parentKey, index), nil
}

type testCustomValuesSlice []string

func (s testCustomValuesSlice) SetURLValues(key string, values *url.Values) error {
	for _, value := range s {
		values.Add("testCustomValuesSlice", value)
	}

	return nil
}

func Test_Encode(t *testing.T) {
	type StructSliceField []struct {
		StringField string
		IntField    int
	}

	type TestStructField struct {
		StringField      string
		IntField         int
		StringSliceField []string
		StructSliceField StructSliceField
		MapField         map[string]string
		MapFieldCustom   testMapCustomKey
	}

	tests := []struct {
		name       string
		input      interface{}
		wantValues url.Values
		wantError  bool
	}{
		{
			name:      "nil",
			wantError: true,
		},
		{
			name:  "empty struct",
			input: struct{}{},
			// empty, but not nil, values
			wantValues: url.Values{},
		},
		{
			name:  "struct with zero values",
			input: TestStructField{},
			wantValues: url.Values{
				"StringField":                  []string{""},
				"IntField":                     []string{"0"},
				"StringSliceField":             []string{""},
				"StructSliceField.StringField": []string{""},
				"StructSliceField.IntField":    []string{"0"},
			},
		},
		{
			name: "struct with zero values, omitempty",
			input: struct {
				StringField string   `values:",omitempty"`
				IntField    int      `values:",omitempty"`
				SliceField  []string `values:",omitempty"`
			}{},
			wantValues: url.Values{},
		},
		{
			name: "struct with omitted field",
			input: struct {
				StringValue string `values:"-"`
			}{
				StringValue: "testString",
			},
			wantValues: url.Values{},
		},
		{
			name: "struct with field named -",
			input: struct {
				StringValue string `values:"-,"`
			}{
				StringValue: "testString",
			},
			wantValues: url.Values{
				"-": []string{"testString"},
			},
		},
		{
			name: "struct with non-zero values",
			input: TestStructField{
				StringField:      "testString",
				IntField:         1,
				StringSliceField: []string{"one", "two"},
				StructSliceField: StructSliceField{
					{
						StringField: "testString",
						IntField:    1,
					},
				},
				MapField: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				MapFieldCustom: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			wantValues: url.Values{
				"StringField":                  []string{"testString"},
				"IntField":                     []string{"1"},
				"StringSliceField":             []string{"one", "two"},
				"StructSliceField.StringField": []string{"testString"},
				"StructSliceField.IntField":    []string{"1"},
				"MapField.key1":                []string{"value1"},
				"MapField.key2":                []string{"value2"},
				"MapFieldCustom[key1]":         []string{"value1"},
				"MapFieldCustom[key2]":         []string{"value2"},
			},
		},
		{
			name: "nested struct",
			input: struct {
				StructField TestStructField
			}{},
			wantValues: url.Values{
				"StructField.StringField":                  []string{""},
				"StructField.IntField":                     []string{"0"},
				"StructField.StringSliceField":             []string{""},
				"StructField.StructSliceField.StringField": []string{""},
				"StructField.StructSliceField.IntField":    []string{"0"},
			},
		},
		{
			name: "nested struct, anonymous",
			input: struct {
				TestStructField
			}{},
			wantValues: url.Values{
				"StringField":                  []string{""},
				"IntField":                     []string{"0"},
				"StringSliceField":             []string{""},
				"StructSliceField.StringField": []string{""},
				"StructSliceField.IntField":    []string{"0"},
			},
		},
		{
			name: "nested struct, anonymized",
			input: struct {
				StructField TestStructField `values:",anonymize"`
			}{},
			wantValues: url.Values{
				"StringField":                  []string{""},
				"IntField":                     []string{"0"},
				"StringSliceField":             []string{""},
				"StructSliceField.StringField": []string{""},
				"StructSliceField.IntField":    []string{"0"},
			},
		},
		{
			name: "custom encoder",
			input: testCustomEncoder{
				key:   "customKey",
				value: "customValue",
			},
			wantValues: url.Values{
				"customKey": []string{"customValue"},
			},
		},
		{
			name: "map (not in struct)",
			input: map[string]string{
				"key1": "value1",
			},
			wantValues: url.Values{
				"key1": []string{"value1"},
			},
		},
		{
			name: "slice (not in struct)",
			input: []string{
				"value1",
				"value2",
			},
			wantError: true,
		},
		{
			// testCustomKeySlice implements custom key by implementing the GetURLKey interface
			name: "custom key slice",
			input: struct {
				Values testCustomKeySlice
			}{
				Values: testCustomKeySlice{
					"value1",
					"value2",
				},
			},
			wantValues: url.Values{
				"Values.0": []string{"value1"},
				"Values.1": []string{"value2"},
			},
		},
		{
			// testCustomValuesSlice implements custom key by implementing the EncodeValues interface
			name: "custom values slice",
			input: testCustomValuesSlice{
				"value1",
				"value2",
			},
			wantValues: url.Values{
				"testCustomValuesSlice": []string{
					"value1",
					"value2",
				},
			},
		},
	}

	for _, test := range tests {
		gotValues, err := Encode(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: Encode() returned error? %v (%s)", test.name, gotError, err)
		}

		if !reflect.DeepEqual(gotValues, test.wantValues) {
			t.Errorf("%s: Encode() got\n%#v, want\n%#v", test.name, gotValues, test.wantValues)
		}
	}
}
