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

// requestBuilder defines a function that performs an operation on an http.Request.
type requestBuilder func(*http.Request) error

// composeRequestBuilder creates a new buildRequestFunc that performs each buildRequestFunc
// provided as an argument.
func composeRequestBuilder(mods ...requestBuilder) requestBuilder {
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
func buildRequest(mods ...requestBuilder) (*http.Request, error) {
	r := &http.Request{}

	for _, mod := range mods {
		if err := mod(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// buildRequestMethod returns a requestBuilder that sets the given method.
func buildRequestMethod(method string) requestBuilder {
	return func(r *http.Request) error {
		r.Method = method

		return nil
	}
}

// buildRequestServiceURL returns a requestBuilder that sets the URL to the ServiceURL
// for a given Service.
func buildRequestServiceURL(c *Client, service Service) requestBuilder {
	return func(r *http.Request) error {
		u, err := c.ServiceURL(service)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// buildRequestCollectionURL returns a requestBuilder that sets the URL to the CollectionURL
// for a given Collection.
func buildRequestCollectionURL(c *Client, collection Collection) requestBuilder {
	return func(r *http.Request) error {
		u, err := c.CollectionURL(collection)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// buildRequestCollectionURLWithTitle returns a requestBuilder that sets the URL to the CollectionURL
// for a given Collection, but also checks that the Collection's Title is not empty.
func buildRequestCollectionURLWithTitle(c *Client, collection Collection) requestBuilder {
	return func(r *http.Request) error {
		if !collection.HasTitle() {
			return fmt.Errorf("Title is required")
		}

		return buildRequestCollectionURL(c, collection)(r)
	}
}

// buildRequestBodyValues returns a requestBuilder that sets the Body to the encoded url.Values for
// a given interface.
func buildRequestBodyValues(i interface{}) requestBuilder {
	return func(r *http.Request) error {
		v, err := query.Values(i)
		if err != nil {
			return err
		}

		r.Body = io.NopCloser(strings.NewReader(v.Encode()))

		return nil
	}
}

// buildRequestBodyValuesWithTitle returns a requestBuilder that sets the Body to the encoded url.Values
// for a given Titler. It checks that the Title is not empty.
func buildRequestBodyValuesWithTitle(t Titler) requestBuilder {
	return func(r *http.Request) error {
		if !t.HasTitle() {
			return fmt.Errorf("Title is required")
		}

		return buildRequestBodyValues(t)(r)
	}
}

// buildRequestBodyContentValues returns a requestBuilder that sets the Body to the encoded url.Values
// for a given ContentGetter.
func buildRequestBodyContentValues(c ContentGetter) requestBuilder {
	return func(r *http.Request) error {
		content := c.GetContent(c)
		if content == nil {
			return fmt.Errorf("unable to GetContent")
		}

		return buildRequestBodyValues(content)(r)
	}
}

// buildRequestOutputModeJSON returns a requestBuilder that sets the URL's RawQuery to output_mode=json.
// It checks that the URL is already set, so it must be applied after setting the URL. It overwrites
// any existing RawQuery Values.
func buildRequestOutputModeJSON() requestBuilder {
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

// buildRequestAuthenticate returns a requestBuilder that authenticates a request for a given Client.
func buildRequestAuthenticate(c *Client) requestBuilder {
	return func(r *http.Request) error {
		return c.Authenticator.AuthenticateRequest(c, r)
	}
}
