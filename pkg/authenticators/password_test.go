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

package authenticators

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/splunk/go-sdk/pkg/client"
	"github.com/splunk/go-sdk/pkg/internal/checks"
)

func TestPassword_loginRequest(t *testing.T) {
	tests := []struct {
		name          string
		inputPassword *Password
		inputClient   *client.Client
		wantError     bool
		requestChecks checks.CheckRequestFunc
	}{
		{
			"missing credentials",
			&Password{},
			&client.Client{URL: "https://localhost:8089"},
			true,
			nil,
		},
		{
			"valid",
			&Password{Username: "admin", Password: "changeme"},
			&client.Client{URL: "https://localhost:8089"},
			false,
			checks.ComposeCheckRequestFunc(
				checks.CheckRequestMethod(http.MethodPost),
				checks.CheckRequestURL("https://localhost:8089/services/auth/login"),
				checks.CheckRequestBodyValue("username", "admin"),
				checks.CheckRequestBodyValue("password", "changeme"),
			),
		},
	}

	for _, test := range tests {
		r, err := test.inputPassword.loginRequest(test.inputClient)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s loginRequest returned error? %v", test.name, gotError)
		}

		if test.requestChecks != nil {
			test.requestChecks(r, t)
		}
	}
}

func TestPassword_handleLoginResponse(t *testing.T) {
	tests := []struct {
		name          string
		inputResponse *http.Response
		wantError     bool
		wantPassword  *Password
	}{
		{
			"empty response",
			nil,
			true,
			&Password{},
		},
		{
			"non-200 status",
			&http.Response{
				StatusCode: http.StatusUnauthorized,
				Body: io.NopCloser(strings.NewReader(`
				<?xml version="1.0" encoding="UTF-8"?>
				<response>
				  <messages>
					<msg type="WARN" code="incorrect_username_or_password">Login failed</msg>
				  </messages>
				</response>
			`)),
			},
			true,
			&Password{},
		},
		{
			"valid",
			&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`
					<response>
  						<sessionKey>fake-session-key</sessionKey>
  						<messages>
    						<msg code=""></msg>
  						</messages>
					</response>
				`))},
			false,
			&Password{SessionKey: SessionKey{SessionKey: "fake-session-key"}},
		},
	}

	for _, test := range tests {
		gotPassword := &Password{}
		err := gotPassword.handleLoginResponse(test.inputResponse)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s handleLoginResponse returned error? %v (%s)", test.name, gotError, err)
		}

		if !reflect.DeepEqual(gotPassword, test.wantPassword) {
			t.Errorf("%s handleLoginResponse got\n%#v, want\n%#v", test.name, gotPassword, test.wantPassword)
		}
	}
}
