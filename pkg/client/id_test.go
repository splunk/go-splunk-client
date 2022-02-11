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

func TestID_GetIDParts(t *testing.T) {
	tests := []struct {
		name         string
		inputID      ID
		wantUser     string
		wantApp      string
		wantEndpoint string
		wantTitle    string
		wantError    bool
	}{
		{
			"global",
			"https://localhost:8089/services/authorization/roles/testrole",
			"",
			"",
			"authorization/roles",
			"testrole",
			false,
		},
		{
			"user/app namespace",
			"https://localhost:8089/servicesNS/nobody/search/saved/searches/testsearch",
			"nobody",
			"search",
			"saved/searches",
			"testsearch",
			false,
		},
		{
			"no services/servicesNS",
			"https://localhost:8089/whatisthis/saved/searches/testsearch",
			"",
			"",
			"",
			"",
			true,
		},
	}

	for _, test := range tests {
		gotUser, gotApp, gotEndpoint, gotTitle, err := test.inputID.GetIDParts()
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s GetIDParts returned error? %v (%s)", test.name, gotError, err)
		}

		if !((gotUser == test.wantUser) &&
			(gotApp == test.wantApp) &&
			(gotEndpoint == test.wantEndpoint) &&
			(gotTitle == test.wantTitle)) {
			t.Errorf("%s GetParts got\n%#v, want\n %#v",
				test.name,
				[]string{gotUser, gotApp, gotEndpoint, gotTitle},
				[]string{test.wantUser, test.wantApp, test.wantEndpoint, test.wantTitle},
			)
		}
	}
}
