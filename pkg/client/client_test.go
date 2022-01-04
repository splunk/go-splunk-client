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

	"github.com/splunk/go-sdk/pkg/models"
)

func TestClient_urlForPath(t *testing.T) {
	tests := []struct {
		inputClient    *Client
		inputPath      string
		inputNamespace models.Namespace
		wantURL        string
		wantError      bool
	}{
		{
			&Client{URL: "https://localhost:8089"},
			"auth/login",
			models.GlobalNamespace,
			"https://localhost:8089/services/auth/login",
			false,
		},
	}

	for _, test := range tests {
		gotURL, err := test.inputClient.urlForPath(test.inputPath, test.inputNamespace)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("(%#v).urlForPath(%q, %#v) returned error? %v", test.inputClient, test.inputPath, test.inputNamespace, gotError)
		}

		gotURLString := gotURL.String()
		if gotURLString != test.wantURL {
			t.Errorf("(%#v).urlForPath(%q, %#v) = %s, want %s", test.inputClient, test.inputPath, test.inputNamespace, gotURLString, test.wantURL)
		}
	}
}

func TestClient_requestForLogin(t *testing.T) {
	tests := []struct {
		input        *Client
		wantURL      string
		wantUsername string
		wantPassword string
	}{
		{
			&Client{URL: "https://localhost:8089", Username: "testuser", Password: "testpassword"},
			"https://localhost:8089/services/auth/login",
			"testuser",
			"testpassword",
		},
	}

	for _, test := range tests {
		r, err := test.input.requestForLogin()
		if err != nil {
			t.Fatalf("unexpected error in requestForAuthenticate: %s", err)
		}

		gotURL := r.URL.String()
		if gotURL != test.wantURL {
			t.Errorf("URL = %s, want %s", gotURL, test.wantURL)
		}

		gotUsername, gotPassword, _ := r.BasicAuth()
		if gotUsername != test.wantUsername || gotPassword != test.wantPassword {
			t.Errorf("BasicAuth = (%q, %q), want (%q, %q)", gotUsername, gotPassword, test.wantUsername, test.wantPassword)
		}
	}
}
