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

package client

import "testing"

func TestService_servicePath(t *testing.T) {
	tests := []struct {
		name            string
		input           interface{}
		wantError       bool
		wantServicePath string
	}{
		{
			"non-struct",
			"this string is not a struct",
			true,
			"",
		},
		{
			"struct has no field named servicePath",
			struct {
				path service `service:"test/endpoint"`
			}{},
			true,
			"",
		},
		{
			"servicePath field isn't servicePath type",
			struct {
				servicePath string `service:"test/endpoint"`
			}{},
			true,
			"",
		},
		{
			"servicePath field missing service tag",
			struct {
				service
			}{},
			true,
			"",
		},
		{
			"everything in place",
			struct {
				service `service:"test/endpoint"`
			}{},
			false,
			"test/endpoint",
		},
	}

	s := service{}
	for _, test := range tests {
		gotServicePath, err := s.servicePath(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s returned error? %v", test.name, gotError)
		}

		if gotServicePath != test.wantServicePath {
			t.Errorf("%s got %s, want %s", test.name, gotServicePath, test.wantServicePath)
		}
	}
}
