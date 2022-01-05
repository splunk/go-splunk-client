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
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
	"golang.org/x/net/publicsuffix"
)

type Client struct {
	URL                   string
	Username              string
	Password              string
	TLSInsecureSkipVerify bool
	sessionKey            string
	httpClient            *http.Client
	mu                    sync.Mutex
}

// urlForPath returns a full url.URL for the given path and namespace.
func (c *Client) urlForPath(p string, ns Namespace) (*url.URL, error) {
	ctxPath, err := ns.Path()
	if err != nil {
		return nil, err
	}

	urlString := strings.Join([]string{c.URL, ctxPath, p}, "/")

	return url.Parse(urlString)
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	if c.httpClient == nil {
		c.mu.Lock()

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return nil, fmt.Errorf("unable to create new cookiejar: %s", err)
		}

		c.httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.TLSInsecureSkipVerify,
				},
			},
			Jar: jar,
		}

		c.mu.Unlock()
	}

	return c.httpClient.Do(r)
}

// requestForLogin returns an http.Request that performs authentication.
func (c *Client) requestForLogin() (*http.Request, error) {
	url, err := c.urlForPath("auth/login", GlobalNamespace)
	if err != nil {
		return nil, err
	}

	loginValues, err := query.Values(Login{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		return nil, err
	}

	r := &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader(loginValues.Encode())),
		Header: http.Header{},
	}

	return r, nil
}

// handleResponseForLogin handles the http.Response from a login attempt.
func (c *Client) handleResponseForLogin(r *http.Response) error {
	loginResponse := loginResponseElement{}
	d := xml.NewDecoder(r.Body)
	if err := d.Decode(&loginResponse); err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		if len(loginResponse.Messages.MessageElements) != 1 {
			return fmt.Errorf("expected exactly one returned message during login, got:\n%#v", loginResponse)
		}

		return fmt.Errorf("unexpected status code during login: %d, %s", r.StatusCode, loginResponse.Messages.MessageElements[0].Message)
	}

	c.sessionKey = loginResponse.SessionKey

	return nil
}

// Login logs in to the Splunk instance. If successful, it retains sessionKey for subsequent requests.
func (c *Client) Login() error {
	request, err := c.requestForLogin()
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return c.handleResponseForLogin(response)
}
