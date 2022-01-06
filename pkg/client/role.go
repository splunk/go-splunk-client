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
	"strings"

	"github.com/google/go-querystring/query"
)

// Role represents a Splunk role.
type Role struct {
	title        title
	Capabilities Capabilities `url:"capabilities"`
}

// validate returns an error if Role is invalid. It is invalid if it:
// * has an empty name
func (r Role) validate() error {
	if r.title.value == "" {
		return fmt.Errorf("invalid role, has no title")
	}

	return nil
}

// NewRoleWithTitle returns a new Role with the given title.
func NewRoleWithTitle(t string) *Role {
	return &Role{
		title: title{value: t},
	}
}

// Title returns the title of a Role.
func (r *Role) Title() string {
	return r.title.value
}

// requestForRead returns the http.Request to perform the read action for a Role.
func (r *Role) requestForRead(c *Client) (*http.Request, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("authorization/roles/%s", r.Title())
	url, err := c.urlForPath(path, GlobalNamespace)
	if err != nil {
		return nil, err
	}

	rValues, err := query.Values(r)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader(rValues.Encode())),
	}

	if err := c.authenticateRequest(request); err != nil {
		return nil, err
	}

	return request, nil
}
