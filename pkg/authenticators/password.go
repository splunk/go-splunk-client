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
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
	"github.com/splunk/go-sdk/pkg/client"
	"github.com/splunk/go-sdk/pkg/messages"
)

// Password defines password authentication to Splunk.
type Password struct {
	Username string `url:"username"`
	Password string `url:"password"`

	// UseBasicAuth can be set to true if Basic Authentication should always be used,
	// which causes Username/Password to be passed with each authenticated request.
	UseBasicAuth bool `url:"-"`

	// SessionKey holds the SessionKey after initial authentication occurs. Unless
	// UseBasicAuth is set to true, this SessionKey will be used to authenticate requests.
	SessionKey

	// mu is used to enable locking to prevent race conditions when checking for and obtaining
	// a SessionKey.
	mu sync.Mutex

	// Fields below this point have no values, and only define how to interact with
	// the REST API.
	client.GlobalNamespace
	client.Endpoint `endpoint:"auth/login"`
}

// loginResponse represents the response returned from auth/login.
type loginResponse struct {
	Messages messages.Messages
	SessionKey
}

// loginRequest creates an http.Request to perform a password login to Splunk.
func (p *Password) loginRequest(c *client.Client) (*http.Request, error) {
	if p.Username == "" || p.Password == "" {
		return nil, fmt.Errorf("Password authenticator missing Username or Password")
	}

	u, err := c.ServiceURL(p)
	if err != nil {
		return nil, err
	}

	v, err := query.Values(p)
	if err != nil {
		return nil, err
	}

	r := &http.Request{
		Method: http.MethodPost,
		URL:    u,
		Body:   io.NopCloser(strings.NewReader(v.Encode())),
	}

	return r, nil
}

// handleLoginResponse handles a response from a password login request.
func (p *Password) handleLoginResponse(r *http.Response) error {
	if r == nil {
		return fmt.Errorf("handleLoginResponse unable to process nil http.Response")
	}

	lR := loginResponse{}

	d := xml.NewDecoder(r.Body)
	if err := d.Decode(&lR); err != nil {
		return fmt.Errorf("unable to parse login response: %s", err)
	}

	if r.StatusCode != http.StatusOK {
		m, ok := lR.Messages.FirstAndOnly()
		if !ok {
			// this is an unlikely situation, but passing messages back to the calling
			// function in case it does occur
			return fmt.Errorf("login failed with multiple messages: %s - %v", r.Status, lR.Messages.Items)
		}

		return fmt.Errorf("unable to log in: %s", m.Value)
	}

	p.SessionKey = lR.SessionKey

	return nil
}

// authenticate performs the authentication request and handles the response, storing the SessionKey
// if successful.
func (p *Password) authenticate(c *client.Client) error {
	req, err := p.loginRequest(c)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return p.handleLoginResponse(resp)
}

// authenticateOnce calls authenticate only if currently unauthenticated.
func (p *Password) authenticateOnce(c *client.Client) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.authenticated() {
		return p.authenticate(c)
	}

	return nil
}

// AuthenticateRequest adds authentication to an http.Request.
func (p *Password) AuthenticateRequest(c *client.Client, r *http.Request) error {
	if err := p.authenticateOnce(c); err != nil {
		return err
	}

	return p.SessionKey.AuthenticateRequest(c, r)
}
