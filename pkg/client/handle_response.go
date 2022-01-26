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
	"encoding/xml"
	"fmt"
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
		return xml.NewDecoder(r.Body).Decode(i)
	}
}

// HandleResponseXMLMessagesError returns a ResponseHandler that decodes an http.Response's Body
// as an XML document of Messages and returns the Messages as an error.
func HandleResponseXMLMessagesError() ResponseHandler {
	return func(r *http.Response) error {
		response := struct {
			Messages messages.Messages
		}{}

		if err := HandleResponseXML(&response)(r); err != nil {
			return err
		}

		return fmt.Errorf(response.Messages.String())
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
