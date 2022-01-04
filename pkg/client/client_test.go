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
	"io"
	"net/http"
	"reflect"
	"strings"
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
		input    *Client
		wantURL  string
		wantBody string
	}{
		{
			&Client{URL: "https://localhost:8089", Username: "testuser", Password: "testpassword"},
			"https://localhost:8089/services/auth/login",
			"password=testpassword&username=testuser",
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

		gotBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unexpected error reading Body: %s", err)
		}

		if string(gotBody) != test.wantBody {
			t.Errorf("Body = %q, want %q", gotBody, test.wantBody)
		}
	}
}

func TestClient_handleResponseForLogin(t *testing.T) {
	tests := []struct {
		inputClient   *Client
		inputResponse *http.Response
		wantClient    *Client
		wantError     bool
	}{
		{
			&Client{},
			&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`
				<response>
					<sessionKey>FakeSessionKey</sessionKey>
					<messages>
		  				<msg code=""></msg>
					</messages>
	  			</response>
				`)),
			},
			&Client{sessionKey: "FakeSessionKey"},
			false,
		},
		{
			&Client{},
			&http.Response{
				StatusCode: http.StatusUnauthorized,
				Body: io.NopCloser(strings.NewReader(`
				<response>
  					<messages>
    					<msg type="WARN" code="incorrect_username_or_password">Login failed</msg>
  					</messages>
				</response>
				`)),
			},
			&Client{},
			true,
		}}

	for _, test := range tests {
		gotError := test.inputClient.handleResponseForLogin(test.inputResponse) != nil

		if gotError != test.wantError {
			t.Errorf("Client.handleResponseForLogin(%#v) returned error? %v", test.inputResponse, gotError)
		}

		if !reflect.DeepEqual(test.inputClient, test.wantClient) {
			t.Errorf("Client.handleResponseForLogin(%#v) = %#v, want %#v", test.inputResponse, test.inputClient, test.wantClient)
		}
	}
}
