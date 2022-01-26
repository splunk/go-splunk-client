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
	"io"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

// RequestBuilder defines a function that performs an operation on an http.Request.
type RequestBuilder func(*http.Request) error

// ComposeRequestBuilder creates a new RequestBuilder that performs each RequestBuilder
// provided as an argument, returning the first error encountered, if any.
func ComposeRequestBuilder(builders ...RequestBuilder) RequestBuilder {
	return func(r *http.Request) error {
		for _, builder := range builders {
			if err := builder(r); err != nil {
				return err
			}
		}

		return nil
	}
}

// buildRequest creates a new http.Request and applies the provided RequestBuilder.
func buildRequest(builder RequestBuilder) (*http.Request, error) {
	r := &http.Request{}

	if err := builder(r); err != nil {
		return nil, err
	}

	return r, nil
}

// BuildRequestMethod returns a RequestBuilder that sets the given method.
func BuildRequestMethod(method string) RequestBuilder {
	return func(r *http.Request) error {
		r.Method = method

		return nil
	}
}

// BuildRequestServiceURL returns a RequestBuilder that sets the URL to the ServiceURL
// for a given Service.
func BuildRequestServiceURL(c *Client, service Service) RequestBuilder {
	return func(r *http.Request) error {
		u, err := c.ServiceURL(service)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// BuildRequestBodyValues returns a RequestBuilder that sets the Body to the encoded url.Values for
// a given interface.
func BuildRequestBodyValues(i interface{}) RequestBuilder {
	return func(r *http.Request) error {
		v, err := query.Values(i)
		if err != nil {
			return err
		}

		r.Body = io.NopCloser(strings.NewReader(v.Encode()))

		return nil
	}
}
