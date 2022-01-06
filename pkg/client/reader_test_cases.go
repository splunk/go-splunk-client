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
	"testing"

	"github.com/splunk/go-sdk/pkg/internal"
)

// readerTestCase defines a test case for types that implement the "reader" interface.
type readerTestCase struct {
	name            string
	inputClient     *Client
	inputReader     reader
	wantError       bool
	requestTestFunc internal.TestRequestFunc
}

// test performs the defined test case.
func (test readerTestCase) test(t *testing.T) {
	gotRequest, err := test.inputReader.requestForRead(test.inputClient)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("%s returned error? %v", test.name, gotError)
	}

	if test.requestTestFunc != nil {
		test.requestTestFunc(test.name, gotRequest, t)
	}
}

// readerTestCases is a slice of readerTestCase items.
type readerTestCases []readerTestCase

// test runs the test for each item.
func (tests readerTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}
