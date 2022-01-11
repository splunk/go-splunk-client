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
