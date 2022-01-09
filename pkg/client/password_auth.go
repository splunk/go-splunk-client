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
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

// PasswordAuth authenticates to auth/login and stores the resulting sessionKey for
// future authentication.
type PasswordAuth struct {
	Username       string `url:"username"`
	Password       string `url:"password"`
	sessionKeyAuth SessionKeyAuth
}

// requestForLogin creates an http.Response to authenticate to the auth/login endpoint.
func (p PasswordAuth) requestForLogin(c *Client) (*http.Request, error) {
	if p.Username == "" || p.Password == "" {
		return nil, fmt.Errorf("attempted PasswordAuth login with empty Username or Password")
	}

	loginURL, err := c.urlForPath("auth/login")
	if err != nil {
		return nil, fmt.Errorf("unable to determine loginURL: %s", err)
	}

	loginValues, err := query.Values(p)
	if err != nil {
		// don't include obtained err in the returned error in case it has sensitive values
		return nil, fmt.Errorf("unable to create url.Values for PasswordAuth")
	}
	loginBody := io.NopCloser(strings.NewReader(loginValues.Encode()))

	r := &http.Request{
		Method: http.MethodPost,
		URL:    loginURL,
		Body:   loginBody,
	}

	return r, nil
}

// handleLoginResponse checks the http.Response for the correct status code, parses the output,
// and applies the sessionKey or returns an error as needed.
func (p *PasswordAuth) handleLoginResponse(r *http.Response) error {
	authResponse := struct {
		Messages messages
		SessionKeyAuth
	}{}

	if r.Body == nil {
		return fmt.Errorf("handleLoginResponse passed nil Body in http.Response")
	}

	d := xml.NewDecoder(r.Body)
	if err := d.Decode(&authResponse); err != nil {
		return fmt.Errorf("PasswordAuth unable to parse XML response: %s", err)
	}

	if r.StatusCode != http.StatusOK {
		message, ok := authResponse.Messages.firstAndOnly()
		if !ok {
			return fmt.Errorf("unknown failure, status %s", r.Status)
		}

		return fmt.Errorf("unable to log in: %s: %s", message.Code, message.Value)
	}

	p.sessionKeyAuth = authResponse.SessionKeyAuth

	return nil
}

// AuthenticateRequest adds the SessionKey to the http.Request's Header.
func (p *PasswordAuth) AuthenticateRequest(c *Client, r *http.Request) error {
	return p.sessionKeyAuth.AuthenticateRequest(c, r)
}
