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

	"github.com/splunk/go-sdk/pkg/messages"
)

// responseHandler defines a function that performs an action on an http.Response.
type responseHandler func(*http.Response) error

// composeResponseHandler creates a new responseHandleFunc that runs each responseHandler
// provided as an argument.
func composeResponseHandler(handlers ...responseHandler) responseHandler {
	return func(r *http.Response) error {
		for _, check := range handlers {
			if err := check(r); err != nil {
				return err
			}
		}

		return nil
	}
}

// handleResponse runs each responseHandler against the provided http.Response.
func handleResponse(r *http.Response, handlers ...responseHandler) error {
	for _, check := range handlers {
		if err := check(r); err != nil {
			return err
		}
	}

	return nil
}

// handleResponseRequireCode returns a responseHandler that returns an error if
// the http.Response has a StatusCode that doesn't match the provided code.
func handleResponseRequireCode(code int) responseHandler {
	return func(r *http.Response) error {
		if r.StatusCode != code {
			return fmt.Errorf("got status code %d, expected %d", r.StatusCode, code)
		}

		return nil
	}
}

// handleResponseRequireCodeWithMessage returns a responseHandler that returns
// an error containing returned Message's Value if the http.Response has a StatusCode
// that doesn't match the provided code.
func handleResponseRequireCodeWithMessage(code int) responseHandler {
	return func(r *http.Response) error {
		if err := handleResponseRequireCode(code)(r); err != nil {
			msg := messages.Messages{}
			d := json.NewDecoder(r.Body)
			if err := d.Decode(&msg); err != nil {
				return err
			}

			m, ok := msg.FirstAndOnly()
			if !ok {
				return fmt.Errorf("unexpected status code (%d) didn't return exactly one Message: %v", r.StatusCode, msg.Items)
			}

			return fmt.Errorf("unexpected status code (%d): %s", r.StatusCode, m.Value)
		}

		return nil
	}
}

// handleResponseEntries returns a responseHandler that parses the http.Response Body
// into the list of Collections provided.
func handleResponseEntries[C Collection](entries *[]C) responseHandler {
	return func(r *http.Response) error {
		collectionResponse := struct{
			Items []C `json:"entry"`
		}{}

		d := json.NewDecoder(r.Body)
		if err := d.Decode(&collectionResponse); err != nil {
			return err
		}

		*entries = collectionResponse.Items
	
		return nil
	}
}

// handleResponseEntry returns a responseHandler that parses the http.Response Body
// into the given Collection.
func handleResponseEntry[C Collection](entry *C) responseHandler {
	return func(r *http.Response) error {
		entries := make([]C, 0)

		if err := handleResponseEntries(&entries)(r); err != nil {
			return err
		}

		if len(entries) != 1 {
			return fmt.Errorf("expected exactly 1 entry, got %d", len(entries))
		}

		*entry = entries[0]

		return nil
	}
}
