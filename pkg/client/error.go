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
	"runtime/debug"
)

// ErrorCode identifies the type of error that was encountered.
type ErrorCode int

const (
	// ErrorUndefined is the zero-value, indicating that the error type
	// has not been defined.
	ErrorUndefined ErrorCode = iota

	// ErrorNamespace indicates an error with the Namespace.
	ErrorNamespace
)

// Error represents an encountered error. It adheres to the "error" interface,
// so will be returned as a standard error.
//
// Returned errors can be handled as this Error type:
//
//   if err := c.RequestAndHandle(...); err != nil {
// 	  if clientErr, ok := err.(client.Error) {
// 		  // check clientErr.Code to determine appropriate action
// 	  }
//   }
type Error struct {
	Code       ErrorCode
	Message    string
	Wrapped    error
	StackTrace string
}

// Wrap returns a new Error with the given code, error, and message. The error value
// may be nil if this is a new error.
func wrapError(code ErrorCode, err error, messagef string, messageArgs ...interface{}) Error {
	return Error{
		Code:       code,
		Message:    fmt.Sprintf(messagef, messageArgs...),
		Wrapped:    err,
		StackTrace: string(debug.Stack()),
	}
}

// Error returns the Error's message.
func (err Error) Error() string {
	return err.Message
}
