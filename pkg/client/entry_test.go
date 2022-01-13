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
	"fmt"
	"go-sdk/pkg/client/internal/checks"
	"go-sdk/pkg/client/internal/immutable"
	"reflect"
	"testing"
)

func TestEntryCollection_firstAndOnly(t *testing.T) {
	tests := []struct {
		name         string
		inputEntries interface{}
		wantError    bool
		wantEntry    entry
	}{
		{
			"empty",
			Users{},
			true,
			nil,
		},
		{
			"exactly one",
			[]User{
				{Name: immutable.Name{Value: "user 1"}},
			},
			false,
			User{Name: immutable.Name{Value: "user 1"}},
		},
		{
			"too many",
			[]User{
				{Name: immutable.Name{Value: "user 1"}},
				{Name: immutable.Name{Value: "user 2"}},
			},
			true,
			nil,
		},
		{
			"not a slice",
			"this string is not a slice",
			true,
			nil,
		},
		{
			"slice of non-entry objects",
			[]string{},
			true,
			nil,
		},
	}

	for _, test := range tests {
		gotEntry, err := firstAndOnlyEntry(test.inputEntries)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s firstAndOnlyEntry returned error? %v - %s", test.name, gotError, err)
		}

		if !reflect.DeepEqual(gotEntry, test.wantEntry) {
			t.Errorf("%s firstAndOnlyEntry got\n%#v, want\n%#v", test.name, gotEntry, test.wantEntry)
		}
	}
}

func TestEntry_entriesAsType(t *testing.T) {
	tests := []struct {
		name         string
		inputEntries []entry
		inputType    User
		wantError    bool
		wantEntries  []User
	}{
		{
			"valid",
			[]entry{
				User{},
			},
			User{},
			false,
			[]User{
				{},
			},
		},
	}

	for _, test := range tests {
		gotEntries, err := entriesAsType(test.inputEntries, test.inputType)
		gotError := err != nil

		if gotError != test.wantError {
			fmt.Errorf("%s entriesAsType returned error? %v", test.name, gotError)
		}

		if !reflect.DeepEqual(gotEntries, test.wantEntries) {
			t.Errorf("%s entriesAsType got\n%#v, want\n%#v", test.name, gotEntries, test.wantEntries)
		}
	}
}

func TestEntry_entryReadRequest(t *testing.T) {
	tests := []struct {
		name         string
		inputEntry   entry
		wantError    bool
		checkRequest checks.CheckRequestFunc
	}{
		{
			"missing title",
			User{},
			true,
			checks.ComposeCheckRequestFunc(),
		},
		{
			"valid",
			User{Name: immutable.Name{Value: "admin"}},
			false,
			checks.ComposeCheckRequestFunc(
				checks.CheckRequestURL("https://localhost:8089/services/authentication/users/admin"),
			),
		},
	}

	c := &Client{
		URL:           "https://localhost:8089",
		Authenticator: &SessionKeyAuth{SessionKey: "fake-session-key"},
	}

	for _, test := range tests {
		r, err := entryReadRequest(c, test.inputEntry)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s entryCollectionReadRequest returned error? %v (%s)", test.name, gotError, err)
		}

		test.checkRequest(r, t)
	}
}
