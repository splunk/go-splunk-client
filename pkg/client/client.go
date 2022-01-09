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
	"net/url"
	"strings"
)

// Client defines how to connect and authenticate to the Splunk REST API.
type Client struct {
	URL string
}

// urlForPath returns a url.URL for the given Namespace and path components.
func (c Client) urlForPath(ns Namespace, path ...string) (*url.URL, error) {
	// parts will hold the Client URL, Namespace, and all path components, capacity set to accomodate
	parts := make([]string, 0, len(path)+2)

	nsPart, err := ns.path()
	if err != nil {
		return nil, err
	}

	parts = append(parts, strings.Trim(c.URL, "/"), nsPart)

	for _, part := range path {
		parts = append(parts, strings.Trim(part, "/"))
	}

	pathURL := strings.Join(parts, "/")

	return url.Parse(pathURL)
}
