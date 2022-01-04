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
	"net/http"
	"net/url"
	"strings"

	"github.com/splunk/go-sdk/pkg/models"
)

type Client struct {
	URL      string
	Username string
	Password string
}

// urlForPath returns a full url.URL for the given path and namespace.
func (c *Client) urlForPath(p string, ns models.Namespace) (*url.URL, error) {
	ctxPath, err := ns.Path()
	if err != nil {
		return nil, err
	}

	urlString := strings.Join([]string{c.URL, ctxPath, p}, "/")

	return url.Parse(urlString)
}

// requestForLogin returns an http.Request that performs authentication.
func (c *Client) requestForLogin() (*http.Request, error) {
	url, err := c.urlForPath("auth/login", models.GlobalNamespace)
	if err != nil {
		return nil, err
	}

	r := &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Header: http.Header{},
	}

	r.SetBasicAuth(c.Username, c.Password)

	return r, nil
}
