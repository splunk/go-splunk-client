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

// Package client implements a client to the Splunk REST API.
package client

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"

	"github.com/splunk/go-splunk-client/pkg/internal/paths"
	"golang.org/x/net/publicsuffix"
)

// Client defines connectivity and authentication to a Splunk REST API.
type Client struct {
	// URL is the URL to the Splunk REST API. It should include the scheme and port number.
	//
	// Example:
	//   https://localhost:8089
	URL string

	// Authenticator defines which authentication method and credentials to use.
	//
	// Example:
	//   authenticators.Password{Username: "admin", Password: "changeme"}
	Authenticator

	// Set TLSInsecureSkipVerify to true to skip TLS verification.
	TLSInsecureSkipVerify bool

	httpClient *http.Client
	mu         sync.Mutex
}

// urlForPath returns a url.URL for path, relative to Client's URL.
func (c *Client) urlForPath(path ...string) (*url.URL, error) {
	if c.URL == "" {
		return nil, wrapError(ErrorMissingURL, nil, "Client has empty URL")
	}

	combinedPath := paths.Join(path...)

	u := paths.Join(c.URL, combinedPath)

	return url.Parse(u)
}

// ServiceURL returns a url.URL for a Service, relative to the Client's URL.
func (c *Client) ServiceURL(s Service) (*url.URL, error) {
	servicePath, err := servicePath(s)
	if err != nil {
		return nil, err
	}

	return c.urlForPath(servicePath)
}

// EntryURL returns a url.URL for an Entry, relative to the Client's URL.
func (c *Client) EntryURL(e Entry) (*url.URL, error) {
	entryPath, err := entryPath(e)
	if err != nil {
		return nil, err
	}

	return c.urlForPath(entryPath)
}

// httpClientPrep prepares the Client's http.Client.
func (c *Client) httpClientPrep() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpClient == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return wrapError(ErrorHTTPClient, err, "unable to create new cookiejar: %s", err)
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

// do performs a given http.Request via the Client's http.Client.
func (c *Client) do(r *http.Request) (*http.Response, error) {
	if err := c.httpClientPrep(); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, wrapError(ErrorHTTPClient, err, "error encountered performing request: %s", err)
	}

	return resp, nil
}

// RequestAndHandle creates a new http.Request from the given RequestBuilder, performs the
// request, and handles the http.Response with the given ResponseHandler.
func (c *Client) RequestAndHandle(builder RequestBuilder, handler ResponseHandler) error {
	req, err := buildRequest(builder)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handler(resp)
}
