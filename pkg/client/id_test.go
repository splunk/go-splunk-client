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
			},
		},
		{
			name:      "no services/servicesNS",
			inputID:   "https://localhost:8089/whatisthis/saved/searches/testsearch",
			wantError: true,
		},
	}

	for _, test := range tests {
		gotID, err := parseID(test.inputID)
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
