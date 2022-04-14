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

package selective

import (
	"reflect"
	"testing"
)

func Test_Encode(t *testing.T) {
	tests := []struct {
		name      string
		inputI    interface{}
		inputTag  string
		wantI     interface{}
		wantError bool
	}{
		{
			name:   "empty",
			inputI: struct{}{},
			wantI:  struct{}{},
		},
		{
			name: "original struct retained, no selective tags present",
			inputI: struct {
				unexportedString string
				unexportedInt    int
				ExportedString   string
				ExportedInt      int
			}{
				unexportedString: "stringValue",
				unexportedInt:    1,
				ExportedString:   "stringValue",
				ExportedInt:      1,
			},
			wantI: struct {
				unexportedString string
				unexportedInt    int
				ExportedString   string
				ExportedInt      int
			}{
				unexportedString: "stringValue",
				unexportedInt:    1,
				ExportedString:   "stringValue",
				ExportedInt:      1,
			},
		},
		{
			name: "exported fields retained, selective tags present",
			inputI: struct {
				unexportedString string
				unexportedInt    int
				ExportedString   string `selective:"matched"`
				ExportedInt      int
			}{
				ExportedString: "stringValue",
				ExportedInt:    1,
			},
			inputTag: "matched",
			wantI: struct {
				ExportedString string `selective:"matched"`
				ExportedInt    int
			}{
				ExportedString: "stringValue",
				ExportedInt:    1,
			},
		},
		{
			name: "embedded struct",
			inputI: struct {
				Embedded struct {
					Value string
				} `selective:"matched"`
			}{
				Embedded: struct {
					Value string
				}{
					Value: "stringValue",
				},
			},
			inputTag: "matched",
			wantI: struct {
				Embedded struct {
					Value string
				} `selective:"matched"`
			}{
				Embedded: struct {
					Value string
				}{
					Value: "stringValue",
				},
			},
		},
		{
			name: "only matched tag",
			inputI: struct {
				MatchedValue   string `selective:"matched"`
				UnmatchedValue string `selective:"unmatched"`
				MatchedStruct  struct {
					AnotherMatchedValue string `selective:"matched"`
				} `selective:"matched"`
				UnmatchedStruct struct {
					// though tagged as "matched", this is under an unmatched tag, so isn't included in want
					AnotherMatchedValue string `selective:"matched"`
				} `selective:"unmatched"`
				UntouchedStruct struct {
					// unexportedValue's presence in the new struct indicates UntouchedStruct was included directly, and not recreated
					unexportedValue string
				} `selctive:"matched"`
			}{
				MatchedValue: "matched",
				UntouchedStruct: struct {
					unexportedValue string
				}{
					unexportedValue: "unexportedValue",
				},
			},
			wantI: struct {
				MatchedValue  string `selective:"matched"`
				MatchedStruct struct {
					AnotherMatchedValue string `selective:"matched"`
				} `selective:"matched"`
				UntouchedStruct struct {
					unexportedValue string
				} `selctive:"matched"`
			}{
				MatchedValue: "matched",
				UntouchedStruct: struct {
					unexportedValue string
				}{
					unexportedValue: "unexportedValue",
				},
			},
			inputTag: "matched",
		},
	}

	for _, test := range tests {
		gotI, err := Encode(test.inputI, test.inputTag)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%q marshalStruct() returned error? %v (%s)", test.name, gotError, err)
		}

		if !reflect.DeepEqual(gotI, test.wantI) {
			t.Errorf("%q marshalStruct() got\n%#v, want\n%#v", test.name, gotI, test.wantI)
		}
	}
}
