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
	"net/http"
)

// PasswordAuth authenticates to auth/login and stores the resulting sessionKey for
// future authentication.
type PasswordAuth struct {
	Username string
	Password string
	// SessionKeyAuth need not be set, as it is managed by the PasswordAuth instance.
	// It is, however, visible to the caller in case there is a need to get the SessionKey.
	SessionKeyAuth
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

	p.SessionKeyAuth = authResponse.SessionKeyAuth

	return nil
}

// AuthenticateRequest adds the SessionKey to the http.Request's Header.
func (p *PasswordAuth) AuthenticateRequest(c *Client, r *http.Request) error {
	return p.SessionKeyAuth.AuthenticateRequest(c, r)
}
