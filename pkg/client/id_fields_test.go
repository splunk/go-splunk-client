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

func TestIDFields_ParseIDFields(t *testing.T) {
	tests := []struct {
		name       string
		inputID    string
		wantFields IDFields
		wantError  bool
	}{
		{
			"empty (global, no title)",
			"",
			IDFields{},
			false,
		},
		{
			"global",
			"https://localhost:8089/services/authorization/roles/testrole",
			IDFields{
				baseURL:  "https://localhost:8089",
				endpoint: "authorization/roles",
				Title:    "testrole",
			},
			false,
		},
		{
			"user/app namespace",
			"https://localhost:8089/servicesNS/nobody/search/saved/searches/testsearch",
			IDFields{
				baseURL:  "https://localhost:8089",
				User:     "nobody",
				App:      "search",
				endpoint: "saved/searches",
				Title:    "testsearch",
			},
			false,
		},
		{
			"minimal empty namespace",
			"services/",
			IDFields{},
			false,
		},
		{
			"minimal user/app namespace",
			"servicesNS/nobody/search/",
			IDFields{
				User: "nobody",
				App:  "search",
			},
			false,
		},
		{
			"wildcard namespace",
			"servicesNS/-/-/",
			IDFields{
				User: "-",
				App:  "-",
			},
			false,
		},
		{
			"no services/servicesNS",
			"https://localhost:8089/whatisthis/saved/searches/testsearch",
			IDFields{},
			true,
		},
		{
			"servicesNS without user/app",
			// at first glance this seems valid, but the last segment will always be the title
			"https://localhost:8089/servicesNS/nobody/search",
			IDFields{},
			true,
		},
	}

	for _, test := range tests {
		gotFields, err := ParseIDFields(test.inputID)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s Fields returned error? %v (%s)", test.name, gotError, err)
		}

		if gotFields != test.wantFields {
			t.Errorf("%s Fields got\n%#v, want\n %#v",
				test.name,
				gotFields,
				test.wantFields,
			)
		}
	}
}
