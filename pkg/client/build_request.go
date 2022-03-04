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
	"net/url"
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
			return wrapError(ErrorValues, err, err.Error())
		}

		r.Body = io.NopCloser(strings.NewReader(v.Encode()))

		return nil
	}
}

// BuildRequestOutputModeJSON returns a RequestBuilder that sets the URL's RawQuery to output_mode=json.
// It checks that the URL is already set, so it must be applied after setting the URL. It overwrites
// any existing RawQuery Values.
func BuildRequestOutputModeJSON() RequestBuilder {
	return func(r *http.Request) error {
		if r.URL == nil {
			return wrapError(ErrorNilValue, nil, "unable to set output mode on nil URL")
		}

		if r.URL.RawQuery != "" {
			return wrapError(ErrorOverwriteValue, nil, "attempted to set output_mode after RawQuery already set")
		}

		r.URL.RawQuery = url.Values{
			"output_mode": []string{"json"},
		}.Encode()

		return nil
	}
}

// BuildRequestBodyValuesWithTitle returns a RequestBuilder that sets the Body to the encoded url.Values
// for a given Titler. It checks that the Title is not empty.
func BuildRequestBodyValuesWithTitle(t Titler) RequestBuilder {
	return func(r *http.Request) error {
		if t.Title() == "" {
			return wrapError(ErrorMissingTitle, nil, "attempted to set request body values of Titler with an empty Title value")
		}

		return BuildRequestBodyValues(t)(r)
	}
}

// BuildRequestBodyValuesContent returns a RequestBuilder that sets the Body to the encoded url.Values
// for a given ContentGetter. The interface returned by c.GetContent(c) will be used to for the resulting
// values.
func BuildRequestBodyValuesContent(c ContentGetter) RequestBuilder {
	return func(r *http.Request) error {
		return BuildRequestBodyValues(c.GetContent(c))(r)
	}
}

// BuildRequestCollectionURL returns a RequestBuilder that sets the URL to the EntryURL
// for a given Entry.
func BuildRequestEntryURL(c *Client, entry Entry) RequestBuilder {
	return func(r *http.Request) error {
		u, err := c.EntryURL(entry)
		if err != nil {
			return err
		}

		r.URL = u

		return nil
	}
}

// BuildRequestEntryURLWithTitle returns a RequestBuilder that sets the URL to the EntryURL
// for a given Entry, but also checks that the Collection's Title is not empty.
func BuildRequestEntryURLWithTitle(c *Client, entry Entry) RequestBuilder {
	return func(r *http.Request) error {
		if entry.Title() == "" {
			return wrapError(ErrorMissingTitle, nil, "attempted to get URLWithTitle of Entry with empty Title")
		}

		return BuildRequestEntryURL(c, entry)(r)
	}
}

func BuildRequestEntryACLURL(c *Client, entry EntryAccessController) RequestBuilder {
	return func(r *http.Request) error {
		url, err := c.EntryACLUrl(entry)
		if err != nil {
			return err
		}

		r.URL = url

		return nil
	}
}

func BuildRequestAccessControllerBodyValues(entry AccessController) RequestBuilder {
	return func(r *http.Request) error {
		return BuildRequestBodyValues(entry.ACLValues())(r)
	}
}

// BuildRequestAuthenticate returns a RequestBuilder that authenticates a request for a given Client.
func BuildRequestAuthenticate(c *Client) RequestBuilder {
	return func(r *http.Request) error {
		return c.Authenticator.AuthenticateRequest(c, r)
	}
}
