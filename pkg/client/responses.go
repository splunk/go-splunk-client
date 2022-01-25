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
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/splunk/go-sdk/pkg/errors"
	"github.com/splunk/go-sdk/pkg/messages"
)

// ResponseHandler defines a function that performs an action on an http.Response.
type ResponseHandler func(*http.Response) error

// ComposeResponseHandler creates a new responseHandleFunc that runs each responseHandler
// provided as an argument.
func ComposeResponseHandler(handlers ...ResponseHandler) ResponseHandler {
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
func handleResponse(r *http.Response, handlers ...ResponseHandler) error {
	for _, check := range handlers {
		if err := check(r); err != nil {
			return err
		}
	}

	return nil
}

func HandleResponseXML(i interface{}) ResponseHandler {
	return func(r *http.Response) error {
		return xml.NewDecoder(r.Body).Decode(i)
	}
}

// HandleResponseMessages returns a responseHandler that returns an error representing
// the parsed Messages from the http.Response.
func HandleResponseMessages(kind errors.Kind) ResponseHandler {
	return func(r *http.Response) error {
		msg, err := messages.NewMessagesFromResponse(r)
		if err != nil {
			return err
		}

		return errors.FromMessages(kind, msg)
	}
}

// HandleResponseMessagesForCode returns a responseHandler that returns an error representing
// the parsed Messages from the http.Response if the response's status code matches the given
// code.
func HandleResponseMessagesForCode(code int, kind errors.Kind) ResponseHandler {
	return func(r *http.Response) error {
		if r.StatusCode == code {
			return HandleResponseMessages(kind)(r)
		}

		return nil
	}
}

// HandleResponseRequireCode returns a responseHandler that returns an error if
// the http.Response has a StatusCode that doesn't match the provided code.
func HandleResponseRequireCode(code int) ResponseHandler {
	return func(r *http.Response) error {
		if r.StatusCode != code {
			inner := HandleResponseMessages(errors.ErrorUnhandled)(r)

			return errors.Wrap(errors.ErrorUnhandled, inner, "got status code %d, expected %d", r.StatusCode, code)
		}

		return nil
	}
}

// HandleResponseRequireCodeWithMessage returns a responseHandler that returns
// an error containing returned Message's Value if the http.Response has a StatusCode
// that doesn't match the provided code.
func HandleResponseRequireCodeWithMessage(code int) ResponseHandler {
	return func(r *http.Response) error {
		if err := HandleResponseRequireCode(code)(r); err != nil {
			return HandleResponseMessages(errors.ErrorUnhandled)(r)
		}

		return nil
	}
}

func HandleResponseJSON(i interface{}) ResponseHandler {
	return func(r *http.Response) error {
		return json.NewDecoder(r.Body).Decode(i)
	}
}

// HandleResponseEntries returns a responseHandler that parses the http.Response Body
// into the list of Collections provided.
func HandleResponseEntries[C Collection](entries *[]C) ResponseHandler {
	return func(r *http.Response) error {
		collectionResponse := struct{
			Items []C `json:"entry"`
		}{}

		if err := HandleResponseJSON(&collectionResponse)(r); err != nil {
			return err
		}

		*entries = collectionResponse.Items
	
		return nil
	}
}

// HandleResponseEntry returns a responseHandler that parses the http.Response Body
// into the given Collection.
func HandleResponseEntry[C Collection](entry *C) ResponseHandler {
	return func(r *http.Response) error {
		entries := make([]C, 0)

		if err := HandleResponseEntries(&entries)(r); err != nil {
			return err
		}

		if len(entries) != 1 {
			return fmt.Errorf("expected exactly 1 entry, got %d", len(entries))
		}

		*entry = entries[0]

		return nil
	}
}
