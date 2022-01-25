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

func TestNamespace_validate(t *testing.T) {
	tests := []struct {
		inputNamespace Namespace
		wantError      bool
	}{
		{
			Namespace{},
			false,
		},
		{
			Namespace{},
			false,
		},
		{
			Namespace{User: "-", App: "-"},
			false,
		},
		{
			Namespace{User: "admin", App: "search"},
			false,
		},
		{
			Namespace{User: "admin"},
			true,
		},
		{
			Namespace{App: "search"},
			true,
		},
	}

	for _, test := range tests {
		gotError := test.inputNamespace.validate() != nil

		if gotError != test.wantError {
			t.Errorf("(%#v).validate() returned error? %v", test.inputNamespace, gotError)
		}
	}
}

func TestNamespace_NamespacePath(t *testing.T) {
	tests := []struct {
		inputNamespace Namespace
		wantPath       string
		wantError      bool
	}{
		{
			Namespace{},
			"services",
			false,
		},
		{
			Namespace{},
			"services",
			false,
		},
		{
			Namespace{User: "admin", App: "search"},
			"servicesNS/admin/search",
			false,
		},
		{
			Namespace{User: "admin"},
			"",
			true,
		},
	}

	for _, test := range tests {
		gotPath, err := test.inputNamespace.NamespacePath()
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("(%#v).Path() returned error? %v", test.inputNamespace, gotError)
		}

		if gotPath != test.wantPath {
			t.Errorf("(%#v).Path() = %s, want %s", test.inputNamespace, gotPath, test.wantPath)
		}
	}
}
