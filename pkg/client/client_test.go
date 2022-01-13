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

func TestClient_urlForPath(t *testing.T) {
	tests := []struct {
		name        string
		inputClient *Client
		inputPath   []string
		wantURL     string
	}{
		{
			"excessive slashes removed",
			&Client{URL: "https://localhost:8089/"},
			[]string{"/path1/", "/path2/"},
			"https://localhost:8089/path1/path2",
		},
	}

	for _, test := range tests {
		gotURL, err := test.inputClient.urlForPath(test.inputPath...)
		if err != nil {
			t.Fatalf("%s urlForPath unexpected error parsing URL: %s", test.name, err)
		}

		gotURLString := gotURL.String()
		if gotURLString != test.wantURL {
			t.Errorf("%s urlForPath got\n%s, want %s", test.name, gotURLString, test.wantURL)
		}
	}
}
