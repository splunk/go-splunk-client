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

package internal

import (
	"encoding/xml"
	"reflect"
	"testing"
)

// XMLUnmarshalerTestCase represents an individual test case for unmarshaling from XML.
type XMLUnmarshalerTestCase struct {
	Input            string
	GotInterfacePtr  interface{}
	WantInterfacePtr interface{}
	WantError        bool
}

// Test performs the test for a XMLUnmarshalerTestCase definition.
func (test XMLUnmarshalerTestCase) Test(t *testing.T) {
	err := xml.Unmarshal([]byte(test.Input), test.GotInterfacePtr)
	gotError := err != nil

	if gotError != test.WantError {
		t.Errorf("xml.Unmarshal(%q) returned error? %v", test.Input, gotError)
	}

	if !reflect.DeepEqual(test.GotInterfacePtr, test.WantInterfacePtr) {
		t.Errorf("xml.Unmarshal(%q) = %#v, want %#v", test.Input, test.GotInterfacePtr, test.WantInterfacePtr)
	}
}

// XMLUnmarshalerTestCases is a slice of XMLUnmarshalerTestCase objects.
type XMLUnmarshalerTestCases []XMLUnmarshalerTestCase

// Test performs each XMLUnmarshalerTestCase in XMLUnmarshalerTestCases.
func (tests XMLUnmarshalerTestCases) Test(t *testing.T) {
	for _, test := range tests {
		test.Test(t)
	}
}
