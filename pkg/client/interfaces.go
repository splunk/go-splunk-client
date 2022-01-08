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
	"strings"

	"github.com/splunk/go-sdk/pkg/internal"
)

// titler represents types that have gettable/settable titles.
type titler interface {
	getTitle() string
	setTitle(string)
}

// validator represents types that can validate themselves.
type validator interface {
	validate() error
}

// reader represents types that define read requests.
type reader interface {
	titler
	validator
	readURL(*Client) (*url.URL, error)
}

func requestForRead(r reader, c *Client) (*http.Request, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	readUrl, err := r.readURL(c)
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

func Read(r reader, c *Client) error {
	request, err := requestForRead(r, c)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// fmt.Printf("response: %s\n", internal.MustReadAll(response.Body))

	responseValue := struct {
		Title   internal.Name `json:"title"`
		Content interface{}    `json:"content"`
	}{
		Content: &r,
	}

	if err := json.Unmarshal([]byte{}, &response); err != nil {
		return err
	}

	fmt.Printf("%#v\n", responseValue)
	// r.setTitle(responseValue.Title.value)

	return nil
}
