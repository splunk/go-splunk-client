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

package errors

import (
	"fmt"

	"github.com/splunk/go-sdk/pkg/messages"
)

type Kind int

const (
	ErrorUnhandled Kind = iota
	ErrorHTTPNotFound
)

type Error struct {
	Kind
	Wrapped error
	Message string
}

func Wrap(kind Kind, err error, messageF string, messageArgs ...interface{}) Error {
	return Error{
		Kind:    kind,
		Wrapped: err,
		Message: fmt.Sprintf(messageF, messageArgs...),
	}
}

func New(kind Kind, messageF string, messageArgs ...interface{}) Error {
	return Error{
		Kind:    kind,
		Message: fmt.Sprintf(messageF, messageArgs...),
	}
}

func FromMessages(kind Kind, messages messages.Messages) Error {
	return Error{
		Kind:    kind,
		Message: messages.String(),
	}
}

func (e Error) Error() string {
	return e.Message
}
