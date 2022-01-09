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
	"go-sdk/pkg/client/internal/checks"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPasswordAuth_AuthenticateRequest(t *testing.T) {
	tests := AuthenticatorTestCases{
		{
			name:               "empty PasswordAuth",
			inputAuthenticator: &PasswordAuth{},
			wantError:          true,
		},
		{
			name: "PasswordAuth with SessionKey",
			inputAuthenticator: &PasswordAuth{
				sessionKeyAuth: SessionKeyAuth{
					SessionKey: "fake-session-key",
				}},
			wantError:    false,
			requestCheck: checks.CheckRequestHeaderKeyValue("Authorization", "Splunk fake-session-key"),
		},
		{
			name: "use basic auth",
			inputAuthenticator: &PasswordAuth{
				Username:     "admin",
				Password:     "changeme",
				UseBasicAuth: true,
			},
			wantError:    false,
			requestCheck: checks.CheckRequestBasicAuth("admin", "changeme"),
		},
	}

	tests.test(t)
}

func TestPasswordAuth_requestForLogin(t *testing.T) {
	tests := []struct {
		name              string
		inputClient       *Client
		inputPasswordAuth *PasswordAuth
		wantError         bool
		requestCheck      checks.CheckRequestFunc
	}{
		{
			"basic",
			&Client{URL: "https://localhost:8089"},
			&PasswordAuth{Username: "admin", Password: "changeme"},
			false,
			checks.ComposeCheckRequestFunc(
				checks.CheckRequestMethod(http.MethodPost),
				checks.CheckRequestBodyValue("password=changeme&username=admin"),
			),
		},
	}

	for _, test := range tests {
		r, err := test.inputPasswordAuth.requestForLogin(test.inputClient)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s requestForLogin returned error? %v", test.name, gotError)
		}

		test.requestCheck(r, t)
	}
}

func TestPasswordAuth_handleLoginResponse(t *testing.T) {
	tests := []struct {
		name               string
		inputResponse      *http.Response
		wantError          bool
		wantSessionKeyAuth SessionKeyAuth
	}{
		{
			"empty body",
			&http.Response{},
			true,
			SessionKeyAuth{},
		},
		{
			"unauthorized",
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
			SessionKeyAuth{},
		},
		{
			"ok",
			&http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`
					<response>
  						<sessionKey>fake-session-key</sessionKey>
  						<messages>
    						<msg code=""></msg>
  						</messages>
					</response>
				`)),
			},
			false,
			SessionKeyAuth{SessionKey: "fake-session-key"},
		},
	}

	for _, test := range tests {
		p := &PasswordAuth{}
		err := p.handleLoginResponse(test.inputResponse)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s handleLoginResponse returned error? %v", test.name, gotError)
		}

		if p.sessionKeyAuth != test.wantSessionKeyAuth {
			t.Errorf("%s SessionKeyAuth=%#v, want #%v", test.name, p.sessionKeyAuth, test.wantSessionKeyAuth)
		}
	}
}
