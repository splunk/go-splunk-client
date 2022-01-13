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
	"io"
	"net/http"
	"reflect"
	"strings"
)

// entry types represent a single entry in an entryCollection.
type entry interface {
	// servicePath is provided by service
	servicePath(interface{}) (string, error)
	// endpointPath is provided by Namespace
	endpointPath(...string) (string, error)
	// title is provided by Title
	title() string
}

func entryCollectionPath(e entry) (string, error) {
	servicePath, err := e.servicePath(e)
	if err != nil {
		return "", err
	}

	return e.endpointPath(servicePath)
}

func entryPath(e entry) (string, error) {
	servicePath, err := e.servicePath(e)
	if err != nil {
		return "", err
	}

	if e.title() == "" {
		return "", fmt.Errorf("entry has no Title")
	}

	return e.endpointPath(servicePath, e.title())
}

func entryCollectionReadRequest(c *Client, e entry) (*http.Request, error) {
	path, err := entryCollectionPath(e)
	if err != nil {
		return nil, err
	}

	u, err := c.urlForPath(path)
	if err != nil {
		return nil, err
	}

	r := &http.Request{
		URL: u,
		Method: http.MethodGet,
		Body: io.NopCloser(strings.NewReader("output_mode=json")),
	}

	if err := c.authenticateRequest(r); err != nil {
		return nil, err
	}

	return r, nil
}

func entryReadRequest(c *Client, e entry) (*http.Request, error) {
	r, err := entryCollectionReadRequest(c, e)
	if err != nil {
		return nil, err
	}

	path, err := entryPath(e)
	if err != nil {
		return nil, err
	}

	u, err := c.urlForPath(path)
	if err != nil {
		return nil, err
	}

	r.URL = u

	return r, nil
}

func handleEntryReadResponse[E entry](e E, r *http.Response) error {
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("entry read failed: %s", r.Status)
	}

	responseEntry := struct{
		Entries []E `json:"entry"`
	}{}

	j := json.NewDecoder(r.Body)
	if err := j.Decode(&responseEntry); err != nil {
		return err
	}

	if len(responseEntry.Entries) == 0 {
		return fmt.Errorf("not found")
	}

	if len(responseEntry.Entries) > 1 {
		return fmt.Errorf("too many entries found, this should not ever happen")
	}

	foundEntry := responseEntry.Entries[0]
	foundEntryV := reflect.Indirect(reflect.ValueOf(foundEntry))

	eV := reflect.Indirect(reflect.ValueOf(e))
	eV.Set(foundEntryV)

	return nil
}

// firstAndOnlyEntry returns the only entry in a slice of entry objects. If the items
// in the given interface aren't of the entry type or if too many or too few items are
// present, an error is returned.
func firstAndOnlyEntry(entries interface{}) (entry, error) {
	entriesV := reflect.ValueOf(entries)
	if entriesV.Kind() != reflect.Slice {
		return nil, fmt.Errorf("entryCollection.Entries is not a slice")
	}

	if entriesV.Len() == 0 {
		return nil, fmt.Errorf("no entries present")
	}

	if entriesV.Len() > 1 {
		return nil, fmt.Errorf("more than one entry present, which should never happen")
	}

	foundEntryV := entriesV.Index(0)

	entryType := reflect.TypeOf((*entry)(nil)).Elem()
	if !foundEntryV.Type().Implements(entryType) {
		return nil, fmt.Errorf("non-entry value found")
	}

	return foundEntryV.Interface().(entry), nil
}

func entriesAsType[E entry](entries []entry, entryType E) ([]E, error) {
	newEntries := make([]E, len(entries))

	for i, e := range entries {
		fmt.Printf("entry\n")
		if reflect.TypeOf(e) != reflect.TypeOf(entryType) {
			return nil, fmt.Errorf("unable to convert list item of type %T to type %T", e, entryType)
		}

		newEntries[i] = e.(E)
	}
	
	return newEntries, nil
}

// func readEntry[E entry](c *Client, e E) error {
func readEntry[E entry](c *Client, e E) error {
	req, err := entryReadRequest(c, e)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleEntryReadResponse(e, resp)
}
