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
	"net/url"
	"reflect"
	"strings"
)

type namespace interface {
	namespacePath() (string, error)
}

type EntryCollection interface {
	namespace
	collectionPath() string
	EntryWithTitle(title string) (Entry, bool)
}

func entryCollectionReadRequest(e EntryCollection, c *Client) (*http.Request, error) {
	readUrl, err := c.urlForPath(e.collectionPath(), e)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		URL:    readUrl,
		Method: http.MethodGet,
		Body:   io.NopCloser(strings.NewReader(url.Values{"output_mode": []string{"json"}}.Encode())),
	}

	if err := c.authenticateRequest(request); err != nil {
		return nil, err
	}

	return request, nil
}

func entryReadRequest(e Entry, c *Client) (*http.Request, error) {
	entryPath := fmt.Sprintf("%s/%s", e.entryCollection().collectionPath(), e.Title())

	readUrl, err := c.urlForPath(entryPath, e.entryCollection())
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		URL:    readUrl,
		Method: http.MethodGet,
		Body:   io.NopCloser(strings.NewReader(url.Values{"output_mode": []string{"json"}}.Encode())),
	}

	if err := c.authenticateRequest(request); err != nil {
		return nil, err
	}

	return request, nil
}

func RefreshCollection(e EntryCollection, c *Client) error {
	r, err := entryCollectionReadRequest(e, c)
	if err != nil {
		return err
	}

	resp, err := c.do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := json.NewDecoder(resp.Body)
	if err := reader.Decode(e); err != nil {
		return err
	}
	// fmt.Printf("%s\n", string(internal.MustReadAll(resp.Body)))

	return nil
}

func refreshCollectionForEntry(eC EntryCollection, e Entry, c *Client) error {
	r, err := entryReadRequest(e, c)
	if err != nil {
		return err
	}

	resp, err := c.do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := json.NewDecoder(resp.Body)
	if err := reader.Decode(eC); err != nil {
		return err
	}
	// fmt.Printf("%s\n", string(internal.MustReadAll(resp.Body)))

	return nil
}

func fetchCollectionItemWithTitle(c EntryCollection, t string, e Entry) error {
	item, ok := c.EntryWithTitle(t)
	if !ok {
		return fmt.Errorf("item not found with title %s", t)
	}

	itemValue := reflect.ValueOf(item)
	iValue := reflect.ValueOf(e).Elem()
	if itemValue.Type() != iValue.Type() {
		return fmt.Errorf("found item type (%s) differs from passed item type (%s)", itemValue.Type(), iValue.Type())
	}

	iValue.Set(itemValue)

	return nil
}

type Entry interface {
	Title() string
	entryCollection() EntryCollection
}

func RefreshEntry(e Entry, c *Client) error {
	eC := e.entryCollection()

	if err := refreshCollectionForEntry(eC, e, c); err != nil {
		return err
	}

	// fmt.Printf("%#v\n", eC)

	if err := fetchCollectionItemWithTitle(eC, e.Title(), e); err != nil {
		return err
	}

	return nil
}
