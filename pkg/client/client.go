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
	"sync"

	"golang.org/x/net/publicsuffix"
)

// Client defines how to connect and authenticate to the Splunk REST API.
type Client struct {
	URL                   string
	TLSInsecureSkipVerify bool
	Authenticator         Authenticator
	httpClient            *http.Client
	mu                    sync.Mutex
}

// urlForPath returns a url.URL for the given path components.
func (c *Client) urlForPath(paths ...string) (*url.URL, error) {
	pathsString := urlJoin(paths...)
	urlString := urlJoin(c.URL, pathsString)

	return url.Parse(urlString)
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

func ReadCollection[E collectionEntry](c *Client, e E) ([]E, error) {
	return readEntryCollection(c, e)
}
