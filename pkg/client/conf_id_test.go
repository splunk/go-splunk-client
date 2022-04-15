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

func Test_parseConfID(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantConfID ConfID
		wantError  bool
	}{
		{
			name:      "empty",
			wantError: true,
		},
		{
			name:      "malformed, missing file",
			input:     "/services/conf-/general",
			wantError: true,
		},
		{
			name:      "malformed, bad prefix",
			input:     "/services/confs-server/general",
			wantError: true,
		},
		{
			name:  "valid",
			input: "/services/conf-server/general",
			wantConfID: ConfID{
				File:   "server",
				Stanza: "general",
			},
			wantError: false,
		},
	}

	for _, test := range tests {
		gotConfID, err := parseConfID(test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s: parseConfID() returned error? %v (%s)", test.name, gotError, err)
		}

		if gotConfID != test.wantConfID {
			t.Errorf("%s: parseConfID() got\n%#v, want\n%#v", test.name, gotConfID, test.wantConfID)
		}
	}
}

func TestConfID_ServicePaths(t *testing.T) {
	tests := []struct {
		name                 string
		input                ConfID
		inputTag             string
		wantServicePath      string
		wantServicePathError bool
		wantEntryPath        string
		wantEntryPathError   bool
	}{
		{
			name:                 "empty",
			wantServicePathError: true,
			wantEntryPathError:   true,
		},
		{
			name: "valid",
			input: ConfID{
				File:   "server",
				Stanza: "general",
			},
			inputTag:        "configs",
			wantServicePath: "services/configs/conf-server",
			wantEntryPath:   "services/configs/conf-server/general",
		},
	}

	for _, test := range tests {
		gotServicePath, err := test.input.GetServicePath(test.inputTag)
		gotServicePathError := err != nil

		if gotServicePathError != test.wantServicePathError {
			t.Errorf("%s: GetServicePath() returned error? %v (%s)", test.name, gotServicePathError, err)
		}

		if gotServicePath != test.wantServicePath {
			t.Errorf("%s: GetServicePath() got\n%s, want\n%s", test.name, gotServicePath, test.wantServicePath)
		}

		gotEntryPath, err := test.input.GetEntryPath(test.inputTag)
		gotEntryPathError := err != nil

		if gotEntryPathError != test.wantEntryPathError {
			t.Errorf("%s: GetEntryPath() returned error? %v (%s)", test.name, gotEntryPathError, err)
		}

		if gotEntryPath != test.wantEntryPath {
			t.Errorf("%s: GetEntryPath() got\n%s, want\n%s", test.name, gotEntryPath, test.wantEntryPath)
		}
	}
}
