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
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/publicsuffix"
)

// Client defines how to connect and authenticate to the Splunk REST API.
type Client struct {
	URL                   string
	Authenticator         Authenticator
	TLSInsecureSkipVerify bool
	httpClient            *http.Client
	mu                    sync.Mutex
}

// urlForPath returns a url.URL for the given Namespace and path components.
func (c *Client) urlForPath(path ...string) (*url.URL, error) {
	// parts will hold the Client URL and all path components, capacity set to accomodate
	parts := make([]string, 0, len(path)+1)

	parts = append(parts, strings.Trim(c.URL, "/"))

	for _, part := range path {
		parts = append(parts, strings.Trim(part, "/"))
	}

	pathURL := strings.Join(parts, "/")

	return url.Parse(pathURL)
}

// httpClientPrep prepares the Client's http.Client.
func (c *Client) httpClientPrep() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpClient == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return fmt.Errorf("unable to create new cookiejar: %s", err)
		}

		c.httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.TLSInsecureSkipVerify,
				},
			},
			Jar: jar,
		}
	}

	return nil
}

// do performs an http.Request, returning its http.Response and error.
func (c *Client) do(r *http.Request) (*http.Response, error) {
	if err := c.httpClientPrep(); err != nil {
		return nil, err
	}

	return c.httpClient.Do(r)
}

func (c *Client) authenticateRequest(r *http.Request) error {
	if c.Authenticator == nil {
		return fmt.Errorf("Client has no Authenticator")
	}

	return c.Authenticator.authenticateRequest(c, r)
}

func (c *Client) Users(ns Namespace) ([]User, error) {
	return nil, nil

}

func ReadEntry[E entry](c *Client, e E) error {
	return readEntry(c, e)
}
