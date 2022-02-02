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
	"net/http"

	"github.com/splunk/go-sdk/pkg/messages"
)

// ResponseHandler defines a function that performs an action on an http.Response.
type ResponseHandler func(*http.Response) error

// ComposeResponseHandler creates a new ResponseHandler that runs each ResponseHandler
// provided as an argument.
func ComposeResponseHandler(handlers ...ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		for _, handler := range handlers {
			if err := handler(r); err != nil {
				return err
			}
		}

		return nil
	}
}

// HandleResponseXML returns a ResponseHandler that decodes an http.Response's Body
// as XML to the given interface.
func HandleResponseXML(i interface{}) ResponseHandler {
	return func(r *http.Response) error {
		if err := xml.NewDecoder(r.Body).Decode(i); err != nil {
			wrapError(ErrorResponseBody, err, "unable to decode response XML: %s", err)
		}

		return nil
	}
}

func HandleResponseXMLMessagesCustomError(code ErrorCode) ResponseHandler {
	return func(r *http.Response) error {
		response := struct {
			Messages messages.Messages
		}{}

		if err := HandleResponseXML(&response)(r); err != nil {
			return err
		}

		return wrapError(code, nil, "response contained message: %s", response.Messages.String())
	}
}

// HandleResponseXMLMessagesError returns a ResponseHandler that decodes an http.Response's Body
// as an XML document of Messages and returns the Messages as an error.
func HandleResponseXMLMessagesError() ResponseHandler {
	return func(r *http.Response) error {
		return HandleResponseXMLMessagesCustomError(ErrorSplunkMessage)(r)
	}
}

// HandleResponseJSON returns a ResponseHandler that decodes an http.Response's Body
// as JSON to the given interface.
func HandleResponseJSON(i interface{}) ResponseHandler {
	return func(r *http.Response) error {
		if err := json.NewDecoder(r.Body).Decode(i); err != nil {
			return wrapError(ErrorResponseBody, err, "unable to decode response JSON: %s", err)
		}

		return nil
	}
}

func HandleResponseJSONMessagesCustomError(code ErrorCode) ResponseHandler {
	return func(r *http.Response) error {
		msg := messages.Messages{}
		if err := HandleResponseJSON(&msg)(r); err != nil {
			return err
		}

		return wrapError(code, nil, "response contained message: %s", msg.String())
	}
}

// HandleResponseJSONMessagesError returns a ResponseHandler that decode's an http.Response's Body
// as a JSON document of Messages and returns the Messages as an error.
func HandleResponseJSONMessagesError() ResponseHandler {
	return func(r *http.Response) error {
		return HandleResponseJSONMessagesCustomError(ErrorSplunkMessage)(r)
	}
}

func HandleResponseCode(code int, errorResponseHandler ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		if r.StatusCode != code {
			return nil
		}

		return errorResponseHandler(r)
	}
}

// HandleResponseEntryNotFound returns a ResponseHandler that calls errorResponseHandler if the response
// status code is equal to Entry's endpointConfig.codeNotFound.
func HandleResponseEntryNotFound(entry Entry, errorResponseHandler ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		config, err := entry.getEndpointConfig(entry)
		if err != nil {
			return err
		}

		return HandleResponseCode(config.codeNotFound, errorResponseHandler)(r)
	}
}

// HandleResponseRequireCode returns a ResponseHandler that checks for a given StatusCode. If
// the http.Response has a different StatusCode, the provided ResponseHandler will be called
// to return the appopriate error message.
func HandleResponseRequireCode(code int, errorResponseHandler ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		if r.StatusCode == code {
			return nil
		}

		return errorResponseHandler(r)
	}
}

// HandleResponseEntries returns a ResponseHandler that parses the http.Response Body
// into the list of Entry reference provided.
func HandleResponseEntries[E Entry](entries *[]E) ResponseHandler {
	return func(r *http.Response) error {
		entriesResponse := struct{
			Entries []E `json:"entry"`
		}{}

		d := json.NewDecoder(r.Body)
		if err := d.Decode(&entriesResponse); err != nil {
			return wrapError(ErrorResponseBody, err, "unable to decide JSON: %s", err)
		}

		*entries = entriesResponse.Entries
	
		return nil
	}
}

// HandleResponseEntry returns a responseHaResponseHandlerndler that parses the http.Response Body
// into the given Entry.
func HandleResponseEntry[E Entry](entry *E) ResponseHandler {
	return func(r *http.Response) error {
		entries := make([]E, 0)

		if err := HandleResponseEntries(&entries)(r); err != nil {
			return err
		}

		if len(entries) != 1 {
			return wrapError(ErrorResponseBody, nil, "expected exactly 1 entry, got %d", len(entries))
		}

		*entry = entries[0]

		return nil
	}
}
