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
	"reflect"
	"testing"
)

func TestStatusCodes_WithDefaults(t *testing.T) {
	tests := []struct {
		name          string
		inputCodes    StatusCodes
		inputDefaults StatusCodes
		wantCodes     StatusCodes
	}{
		{
			name: "empty",
		},
		{
			name:          "specified not replaced",
			inputCodes:    StatusCodes{Created: 200},
			inputDefaults: StatusCodes{Created: 201},
			wantCodes:     StatusCodes{Created: 200},
		},
		{
			name:          "unspecified replaced",
			inputCodes:    StatusCodes{},
			inputDefaults: StatusCodes{Created: 201},
			wantCodes:     StatusCodes{Created: 201},
		},
	}

	for _, test := range tests {
		gotCodes := test.inputCodes.WithDefaults(test.inputDefaults)

		if !reflect.DeepEqual(gotCodes, test.wantCodes) {
			t.Errorf("%s: WithDefaults() got\n%#v, want\n%#v", test.name, gotCodes, test.wantCodes)
		}
	}
}

func TestServiceCodes_withTagDefaults(t *testing.T) {
	tests := []struct {
		name       string
		inputCodes StatusCodes
		inputTag   string
		wantCodes  StatusCodes
		wantError  bool
	}{
		{
			name: "empty",
		},
		{
			name:      "valid action name, no value",
			inputTag:  "Created",
			wantError: true,
		},
		{
			name:      "valid action name, non-int value",
			inputTag:  "Created=200OK",
			wantError: true,
		},
		{
			name:      "invalid action name, int value",
			inputTag:  "created=200",
			wantError: true,
		},
		{
			name:      "valid action name, int value",
			inputTag:  "Created=200",
			wantCodes: StatusCodes{Created: 200},
		},
		{
			name:       "valid action name, int value, existing value",
			inputCodes: StatusCodes{Created: 201},
			inputTag:   "Created=200",
			wantCodes:  StatusCodes{Created: 201},
		},
	}

	for _, test := range tests {
		gotCodes, err := test.inputCodes.withTagDefaults(test.inputTag)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: withTagDefaults() returned error? %v (%s)", test.name, gotError, err)
		}

		if !reflect.DeepEqual(gotCodes, test.wantCodes) {
			t.Errorf("%s: withTagDefaults() got\n%#v, want\n%#v", test.name, gotCodes, test.wantCodes)
		}
	}
}
