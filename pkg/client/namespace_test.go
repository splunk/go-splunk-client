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

func TestNamespace_validate(t *testing.T) {
	tests := []struct {
		name           string
		inputNamespace Namespace
		wantError      bool
		wantErrorCode  ErrorCode
	}{
		{
			"empty",
			Namespace{},
			false,
			ErrorUndefined,
		},
		{
			"wildcard",
			Namespace{User: "-", App: "-"},
			false,
			ErrorUndefined,
		},
		{
			"user/app",
			Namespace{User: "admin", App: "search"},
			false,
			ErrorUndefined,
		},
		{
			"user only",
			Namespace{User: "admin"},
			true,
			ErrorNamespace,
		},
		{
			"app only",
			Namespace{App: "search"},
			true,
			ErrorNamespace,
		},
	}

	for _, test := range tests {
		err := test.inputNamespace.validate()

		testError(test.name, err, test.wantError, test.wantErrorCode, t)
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

// func TestNamespace_namespaceForID(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		inputId       string
// 		wantNamespace Namespace
// 		wantError     bool
// 	}{
// 		{
// 			"global",
// 			"https://localhost:8089/services/authentication/users/testuser",
// 			Namespace{},
// 			false,
// 		},
// 		{
// 			"namespace",
// 			"https://localhost:8089/servicesNS/admin/search/authentication/users/testuser",
// 			Namespace{
// 				User: "admin",
// 				App:  "search",
// 			},
// 			false,
// 		},
// 		{
// 			"servicesNS in wrong location",
// 			// testuser thrown out, only "roles" remains, which can't be parsed into user/app
// 			"https://localhost:8089/servicesNS/roles/testuser",
// 			Namespace{},
// 			true,
// 		},
// 	}

// 	for _, test := range tests {
// 		gotNamespace, err := ParseID(test.inputId)
// 		gotError := err != nil

// 		if gotError != test.wantError {
// 			t.Errorf("%s namespaceForID returned error? %v (%s)", test.name, gotError, err)
// 		}

// 		if gotNamespace != test.wantNamespace {
// 			t.Errorf("%s namespaceForID got\n%#v, want\n%#v", test.name, gotNamespace, test.wantNamespace)
// 		}
// 	}
// }
