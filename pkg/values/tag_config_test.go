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

package values

import "testing"

func Test_newOptions(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantOptions tagConfig
		wantError   bool
	}{
		{
			name: "empty",
		},
		{
			name:  "name only",
			input: "fieldname",
			wantOptions: tagConfig{
				Name: "fieldname",
			},
		},
		{
			name:  "empty name, omitempty",
			input: ",omitempty",
			wantOptions: tagConfig{
				Omitempty: true,
			},
		},
		{
			name:      "unknown flag",
			input:     ",invalid",
			wantError: true,
		},
		{
			name:  "field name -",
			input: "-,",
			wantOptions: tagConfig{
				Name: "-",
			},
		},
	}

	for _, test := range tests {
		gotOptions, err := parseTagConfig(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: newOptions() returned error? %v (%s)", test.name, gotError, err)
		}

		if gotOptions != test.wantOptions {
			t.Errorf("%s: newOptions() got\n%#v, want\n%#v", test.name, gotOptions, test.wantOptions)
		}
	}
}
