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

package deepset

import (
	"reflect"
	"testing"
)

type title string

type namespace struct {
	User string
	App  string
}

type id struct {
	Title       title
	URL         string
	Description string
	Namespace   namespace
}

type search struct {
	ID      id
	Content string
}

func Test_Set(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		inputValue interface{}
		want       interface{}
		wantError  bool
	}{
		{
			name:      "nil",
			wantError: true,
		},
		{
			name:  "non-pointer",
			input: id{},
			inputValue: id{
				Title: "any value",
			},
			want:      id{},
			wantError: true,
		},
		{
			name:  "full id",
			input: &id{},
			inputValue: id{
				Title: "any value",
			},
			want: &id{
				Title: "any value",
			},
		},
		{
			name:       "embedded title (unambiguous)",
			input:      &id{},
			inputValue: title("any value"),
			want: &id{
				Title: "any value",
			},
		},
		{
			name:  "embedded string (ambiguous)",
			input: &id{},
			// inputValue is just a string (which is ambiguous), not the custom title type
			inputValue: "any value",
			want:       &id{},
			wantError:  true,
		},
		{
			name:  "embedded namespace",
			input: &id{},
			inputValue: namespace{
				User: "any value",
			},
			want: &id{
				Namespace: namespace{
					User: "any value",
				},
			},
		},
		{
			name:  "set namespace in search (multi-level)",
			input: &search{},
			inputValue: namespace{
				User: "any user",
			},
			want: &search{
				ID: id{
					Namespace: namespace{
						User: "any user",
					},
				},
			},
		},
		{
			name:       "embedded ambiguity",
			input:      &struct{}{},
			inputValue: "plain string",
			want:       &struct{}{},
			wantError:  true,
		},
	}

	for _, test := range tests {
		err := Set(test.input, test.inputValue)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: Set() returned error? %v (%s)", test.name, gotError, err)
		}

		if !reflect.DeepEqual(test.input, test.want) {
			t.Errorf("%s: Set() got\n%#v, want\n%#v", test.name, test.input, test.want)
		}
	}
}
