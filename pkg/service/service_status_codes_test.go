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

func Test_ServiceStatusCodes(t *testing.T) {
	tests := []struct {
		name            string
		input           interface{}
		inputDefaults   StatusCodes
		wantStatusCodes StatusCodes
		wantError       bool
	}{
		{
			name:            "none",
			input:           struct{}{},
			inputDefaults:   StatusCodes{},
			wantStatusCodes: StatusCodes{},
		},
		{
			name:            "defaults",
			input:           struct{}{},
			inputDefaults:   StatusCodes{Created: 201},
			wantStatusCodes: StatusCodes{Created: 201},
		},
		{
			name: "has unexported StatusCodes field",
			input: struct {
				_ StatusCodes `service:"Created=200"`
			}{},
			inputDefaults:   StatusCodes{Created: 201},
			wantStatusCodes: StatusCodes{Created: 200},
		},
		{
			name: "has exported StatusCodes field",
			input: struct {
				Codes StatusCodes `service:"Created=201"`
			}{
				Codes: StatusCodes{Created: 202},
			},
			inputDefaults:   StatusCodes{Created: 201},
			wantStatusCodes: StatusCodes{Created: 202},
		},
	}

	for _, test := range tests {
		gotStatusCodes, err := ServiceStatusCodes(test.input, test.inputDefaults)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: ServiceStatusCodes() returned error? %v (%s)", test.name, gotError, err)
		}

		if !reflect.DeepEqual(gotStatusCodes, test.wantStatusCodes) {
			t.Errorf("%s: ServiceStatusCodes() got\n%#v, want\n%#v", test.name, gotStatusCodes, test.wantStatusCodes)
		}
	}
}
