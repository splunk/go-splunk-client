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

func TestCollection_collectionForInterface(t *testing.T) {
	tests := []struct {
		name           string
		input          collection{}
		wantError      bool
		wantCollection collection
	}{
		{
			"empty interface",
			nil,
			true,
			collection{},
		},
		{
			"missing path",
			struct {
				Namespace
			}{},
			true,
			collection{},
		},
		{
			"missing Namespace",
			struct {
				path string
			}{},
			true,
			collection{},
		},
		{
			"valid",
			struct {
				path string `collection:"authorization/users"`
				Namespace
			}{
				Namespace: Namespace{
					User: "admin",
					App:  "search",
				},
			},
			false,
			collection{
				path: "authorization/users",
				namespace: Namespace{
					User: "admin",
					App:  "search",
				},
			},
		},
	}

	for _, test := range tests {
		gotCollection, err := collectionForInterface(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s returned error? %v", test.name, gotError)
		}

		if gotCollection != test.wantCollection {
			t.Errorf("%s got\n%#v, want\n%#v", test.name, gotCollection, test.wantCollection)
		}
	}
}
