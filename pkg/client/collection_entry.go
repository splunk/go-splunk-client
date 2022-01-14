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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// collectionEntry defines the methods that a type must implement to be considered
// an entry within a collection.
type collectionEntry interface {
	collection
	titler
}

// entryPath returns the REST path to a collectionEntry.
func entryPath(e collectionEntry) (string, error) {
	collectionPath, err := collectionPath(e)
	if err != nil {
		return "", err
	}

	title, ok := e.title()
	if !ok {
		return "", fmt.Errorf("entryPath requires a title")
	}

	return urlJoin(collectionPath, title), nil
}

func entryReadRequest(c *Client, e collectionEntry) (*http.Request, error) {
	collPath, err := collectionPath(e)
	if err != nil {
		return nil, err
	}

	// ignore the boolean, if the title is unset we fetch the titleless collection.
	title, _ := e.title()

	u, err := c.urlForPath(collPath, title)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("output_mode", "json")
	v.Set("count", "0")
	u.RawQuery = v.Encode()

	r := &http.Request{
		Method: http.MethodGet,
		URL:    u,
	}

	if err := c.Authenticator.authenticateRequest(c, r); err != nil {
		return nil, err
	}

	return r, nil
}

func entryReadResponseEntries[E collectionEntry](e E, r *http.Response) ([]E, error) {
	if r.StatusCode != http.StatusOK {
		// TODO parse messages
		return nil, fmt.Errorf("read failed: %s", r.Status)
	}

	entriesResponse := struct{
		Entries []E `json:"entry"`
	}{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&entriesResponse); err != nil {
		return nil, err
	}

	return entriesResponse.Entries, nil
}

func readEntryCollection[E collectionEntry](c *Client, e E) ([]E, error) {
	req, err := entryReadRequest(c, e)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	entries, err := entryReadResponseEntries(e, resp)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
