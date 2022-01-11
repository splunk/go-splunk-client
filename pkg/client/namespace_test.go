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
		inputNamespace Namespace
		wantError      bool
	}{
		{
			Namespace{},
			false,
		},
		{
			GlobalNamespace,
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

func TestNamespace_Path(t *testing.T) {
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
			GlobalNamespace,
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
		gotPath, err := test.inputNamespace.path()
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("(%#v).Path() returned error? %v", test.inputNamespace, gotError)
		}

		if gotPath != test.wantPath {
			t.Errorf("(%#v).Path() = %s, want %s", test.inputNamespace, gotPath, test.wantPath)
		}
	}
}

// func TestNamespace_namespace(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		input         interface{}
// 		wantError     bool
// 		wantNamespace Namespace
// 	}{
// 		{
// 			"non-struct",
// 			"this string is not a struct",
// 			true,
// 			Namespace{},
// 		},
// 		{
// 			"empty struct",
// 			struct{}{},
// 			false,
// 			Namespace{},
// 		},
// 		{
// 			"struct lacking Namespace field",
// 			struct {
// 				Name string
// 			}{
// 				Name: "Unused",
// 			},
// 			false,
// 			Namespace{},
// 		},
// 		{
// 			"struct with non-Namespace{} Namespace field",
// 			struct {
// 				Namespace string
// 			}{
// 				Namespace: "Invalid",
// 			},
// 			true,
// 			Namespace{},
// 		},
// 		{
// 			"struct with Namespace{} Namespace field",
// 			struct {
// 				Namespace Namespace
// 			}{
// 				Namespace: Namespace{
// 					User: "ns_user",
// 					App:  "ns_app",
// 				},
// 			},
// 			false,
// 			Namespace{
// 				User: "ns_user",
// 				App:  "ns_app",
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		gotNamespace, err := namespace(test.input)
// 		gotError := err != nil

// 		if gotError != test.wantError {
// 			t.Errorf("%s namespace() returned error? %v", test.name, gotError)
// 		}

// 		if gotNamespace != test.wantNamespace {
// 			t.Errorf("%s namespace() got\n%#v, want\n%#v", test.name, gotNamespace, test.wantNamespace)
// 		}
// 	}
// }
