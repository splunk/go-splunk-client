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

package service

import (
	"fmt"
	"testing"
)

type testServiceFieldPathGetter string

func (t testServiceFieldPathGetter) GetServicePath(path string) (string, error) {
	return path, nil
}

func (t testServiceFieldPathGetter) GetEntryPath(path string) (string, error) {
	servicePath, err := t.GetServicePath(path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", servicePath, t), nil
}

type testObjectHasServicePathGetter struct {
	Name testServiceFieldPathGetter `service:"test/object/field"`
}

type testObjectIsServicePathGetter struct {
	Name string
}

func (t testObjectIsServicePathGetter) GetServicePath(string) (string, error) {
	return "test/object", nil
}

func (t testObjectIsServicePathGetter) GetEntryPath(path string) (string, error) {
	servicePath, err := t.GetServicePath(path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", servicePath, t.Name), nil
}

func Test_ServicePath(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		want      string
		wantError bool
	}{
		{
			name:      "nil",
			input:     nil,
			wantError: true,
		},
		{
			name:      "empty struct",
			input:     struct{}{},
			wantError: true,
		},
		{
			name:  "testObjectHasServicePathGetter tag passed to GetServicePath",
			input: testObjectHasServicePathGetter{},
			want:  "test/object/field",
		},
		{
			name:  "testObjectIsServicePathGetter",
			input: testObjectIsServicePathGetter{},
			want:  "test/object",
		},
	}

	for _, test := range tests {
		got, err := ServicePath(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%q ServicePath returned error? %v (%s)", test.name, gotError, err)
		}

		if got != test.want {
			t.Errorf("%q ServicePath got\n%s, want\n%s", test.name, got, test.want)
		}
	}
}

func Test_EntryPath(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		want      string
		wantError bool
	}{
		{
			name:      "nil",
			input:     nil,
			wantError: true,
		},
		{
			name:      "empty struct",
			input:     struct{}{},
			wantError: true,
		},
		{
			name:  "testObjectHasServicePathGetter tag passed to GetServicePath",
			input: testObjectHasServicePathGetter{Name: "objectId"},
			want:  "test/object/field/objectId",
		},
		{
			name:  "testObjectIsServicePathGetter",
			input: testObjectIsServicePathGetter{Name: "objectName"},
			want:  "test/object/objectName",
		},
	}

	for _, test := range tests {
		got, err := EntryPath(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%q ServicePath returned error? %v (%s)", test.name, gotError, err)
		}

		if got != test.want {
			t.Errorf("%q ServicePath got\n%s, want\n%s", test.name, got, test.want)
		}
	}
}
