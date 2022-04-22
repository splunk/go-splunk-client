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

import (
	"testing"
)

func Test_ParseID(t *testing.T) {
	tests := []struct {
		name      string
		inputID   string
		wantID    ID
		wantError bool
	}{
		{
			name:      "empty (global, no title)",
			wantError: true,
		},
		{
			name:    "global",
			inputID: "https://localhost:8089/services/authorization/roles/testrole",
			wantID: ID{
				Title: "testrole",
				url:   "https://localhost:8089/services/authorization/roles/testrole",
			},
		},
		{
			name:    "user/app namespace",
			inputID: "https://localhost:8089/servicesNS/nobody/search/saved/searches/testsearch",
			wantID: ID{
				Namespace: Namespace{
					User: "nobody",
					App:  "search",
				},
				Title: "testsearch",
				url:   "https://localhost:8089/servicesNS/nobody/search/saved/searches/testsearch",
			},
		},
		{
			name:      "minimal empty namespace",
			inputID:   "services",
			wantError: true,
		},
		{
			name:    "minimal empty namespace, trailing segment",
			inputID: "services/",
			wantID: ID{
				// there is an empty segment, so Title is an empty string
				Title: "",
				url:   "services/",
			},
		},
		{
			name:    "minimal user/app namespace",
			inputID: "servicesNS/nobody/search/",
			wantID: ID{
				Namespace: Namespace{
					User: "nobody",
					App:  "search",
				},
				// there is an empty segment, so Title is an empty string
				Title: "",
				url:   "servicesNS/nobody/search/",
			},
		},
		{
			name:    "wildcard namespace",
			inputID: "servicesNS/-/-/testsearch",
			wantID: ID{
				Namespace: Namespace{
					User: "-",
					App:  "-",
				},
				Title: "testsearch",
				url:   "servicesNS/-/-/testsearch",
			},
		},
		{
			name:      "no services/servicesNS",
			inputID:   "https://localhost:8089/whatisthis/saved/searches/testsearch",
			wantError: true,
		},
	}

	for _, test := range tests {
		gotID, err := ParseID(test.inputID)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: ParseNamespace() returned error? %v (%s)", test.name, gotError, err)
		}

		if gotID != test.wantID {
			t.Errorf("%s: ParseID got\n%#v, want\n %#v",
				test.name,
				gotID,
				test.wantID,
			)
		}
	}
}

func TestID_URL(t *testing.T) {
	tests := []struct {
		name      string
		input     ID
		wantURL   string
		wantError bool
	}{
		{
			name:      "empty",
			wantError: true,
		},
		{
			name: "unset url",
			input: ID{
				Title: "testtitle",
				Namespace: Namespace{
					User: "testuser",
					App:  "testapp",
				},
			},
			wantError: true,
		},
		{
			name: "url mismatch",
			input: ID{
				Namespace: Namespace{
					User: "testuser",
					App:  "changedapp",
				},
				Title: "testtitle",
				url:   "https://localhost:8089/servicesNS/testuser/testapp/service/path/testtitle",
			},
			wantError: true,
		},
		{
			name: "url matches",
			input: ID{
				Namespace: Namespace{
					User: "testuser",
					App:  "testapp",
				},
				Title: "testtitle",
				url:   "https://localhost:8089/servicesNS/testuser/testapp/service/path/testtitle",
			},
			wantURL: "https://localhost:8089/servicesNS/testuser/testapp/service/path/testtitle",
		},
	}

	for _, test := range tests {
		gotURL, err := test.input.URL()
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: URL() returned error? %v (%s)", test.name, gotError, err)
		}

		if gotURL != test.wantURL {
			t.Errorf("%s: URL() got\n%s, want\n%s", test.name, gotURL, test.wantURL)
		}
	}
}
