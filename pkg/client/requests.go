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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// RequestBuilder defines a function that performs an operation on an http.Request.
type RequestBuilder func(*http.Request) error

// ComposeRequestBuilder creates a new buildRequestFunc that performs each buildRequestFunc
// provided as an argument.
func ComposeRequestBuilder(mods ...RequestBuilder) RequestBuilder {
	return func(r *http.Request) error {
		for _, mod := range mods {
			if err := mod(r); err != nil {
				return err
			}
		}

		return nil
	}
}

// buildRequest creates a new http.Request and applies each requestBuilder defined.
// If any of requestBuilder returns an error, buildRequest will return a nil http.Request
// and the error.
func buildRequest(mods ...RequestBuilder) (*http.Request, error) {
	r := &http.Request{}

	for _, mod := range mods {
		if err := mod(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// BuildRequestMethod returns a requestBuilder that sets the given method.
func BuildRequestMethod(method string) RequestBuilder {
	return func(r *http.Request) error {
		r.Method = method

		return nil
	}
}

// BuildRequestServiceURL returns a requestBuilder that sets the URL to the ServiceURL
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

// BuildRequestCollectionURL returns a requestBuilder that sets the URL to the CollectionURL
// for a given Collection.
func BuildRequestCollectionURL(c *Client, collection Collection) RequestBuilder {
	return func(r *http.Request) error {
		u, err := c.CollectionURL(collection)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// BuildRequestCollectionURLWithTitle returns a requestBuilder that sets the URL to the CollectionURL
// for a given Collection, but also checks that the Collection's Title is not empty.
func BuildRequestCollectionURLWithTitle(c *Client, collection Collection) RequestBuilder {
	return func(r *http.Request) error {
		if !collection.HasTitle() {
			return fmt.Errorf("Title is required")
		}

		return BuildRequestCollectionURL(c, collection)(r)
	}
}

// BuildRequestBodyValues returns a requestBuilder that sets the Body to the encoded url.Values for
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

// BuildRequestBodyValuesWithTitle returns a requestBuilder that sets the Body to the encoded url.Values
// for a given Titler. It checks that the Title is not empty.
func BuildRequestBodyValuesWithTitle(t Titler) RequestBuilder {
	return func(r *http.Request) error {
		if !t.HasTitle() {
			return fmt.Errorf("Title is required")
		}

		return BuildRequestBodyValues(t)(r)
	}
}

// BuildRequestBodyContentValues returns a requestBuilder that sets the Body to the encoded url.Values
// for a given ContentGetter.
func BuildRequestBodyContentValues(c ContentGetter) RequestBuilder {
	return func(r *http.Request) error {
		content := c.GetContent(c)
		if content == nil {
			return fmt.Errorf("unable to GetContent")
		}

		return BuildRequestBodyValues(content)(r)
	}
}

// BuildRequestOutputModeJSON returns a requestBuilder that sets the URL's RawQuery to output_mode=json.
// It checks that the URL is already set, so it must be applied after setting the URL. It overwrites
// any existing RawQuery Values.
func BuildRequestOutputModeJSON() RequestBuilder {
	return func(r *http.Request) error {
		if r.URL == nil {
			return fmt.Errorf("unable to set output mode on empty URL")
		}

		r.URL.RawQuery = url.Values{
			"output_mode": []string{"json"},
		}.Encode()

		return nil
	}
}

// BuildRequestAuthenticate returns a requestBuilder that authenticates a request for a given Client.
func BuildRequestAuthenticate(c *Client) RequestBuilder {
	return func(r *http.Request) error {
		return c.Authenticator.AuthenticateRequest(c, r)
	}
}
