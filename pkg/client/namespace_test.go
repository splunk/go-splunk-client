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
	"reflect"
	"testing"
)

func Test_ParseNamespace(t *testing.T) {
	tests := []struct {
		name          string
		inputID       string
		wantNamespace Namespace
		wantRemnants  []string
		wantError     bool
	}{
		{
			name: "empty (global, no title)",
		},
		{
			name:         "global",
			inputID:      "https://localhost:8089/services/authorization/roles/testrole",
			wantRemnants: []string{"authorization", "roles", "testrole"},
		},
		{
			name:    "user/app namespace",
			inputID: "https://localhost:8089/servicesNS/nobody/search/saved/searches/testsearch",
			wantNamespace: Namespace{
				User: "nobody",
				App:  "search",
			},
			wantRemnants: []string{"saved", "searches", "testsearch"},
		},
		{
			name:         "minimal empty namespace",
			inputID:      "services",
			wantRemnants: []string{},
		},
		{
			name:    "minimal empty namespace, trailing segment",
			inputID: "services/",
			// there is an empty segment, so there is a single remnant with the value of an empty string
			wantRemnants: []string{""},
		},
		{
			name:    "minimal user/app namespace",
			inputID: "servicesNS/nobody/search/",
			wantNamespace: Namespace{
				User: "nobody",
				App:  "search",
			},
			// there is an empty segment, so there is a single remnant with the value of an empty string
			wantRemnants: []string{""},
		},
		{
			name:    "wildcard namespace",
			inputID: "servicesNS/-/-/",
			wantNamespace: Namespace{
				User: "-",
				App:  "-",
			},
			// there is an empty segment, so there is a single remnant with the value of an empty string
			wantRemnants: []string{""},
		},
		{
			name:      "no services/servicesNS",
			inputID:   "https://localhost:8089/whatisthis/saved/searches/testsearch",
			wantError: true,
		},
	}

	for _, test := range tests {
		gotNamespace, gotRemnants, err := parseNamespace(test.inputID)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: ParseNamespace() returned error? %v (%s)", test.name, gotError, err)
		}

		if gotNamespace != test.wantNamespace {
			t.Errorf("%s: ParseNamespace got\n%#v, want\n %#v",
				test.name,
				gotNamespace,
				test.wantNamespace,
			)
		}

		if !reflect.DeepEqual(gotRemnants, test.wantRemnants) {
			t.Errorf("%s: ParseNamespace() got remnants\n%#v, want\n%#v", test.name, gotRemnants, test.wantRemnants)
		}
	}
}
