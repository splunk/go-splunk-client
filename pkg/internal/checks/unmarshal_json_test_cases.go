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

package attributes

import (
	"encoding/json"
	"reflect"
	"testing"
)

// jsonUnmarshalTestCase defines a test case for json.Unmarshal.
type jsonUnmarshalTestCase struct {
	name        string
	inputString string
	want        interface{}
	wantError   bool
}

// test runs the test case.
func (test jsonUnmarshalTestCase) test(t *testing.T) {
	// create a new pointer to a zero value of test.want
	gotT := reflect.TypeOf(test.want)
	if gotT == nil {
		t.Fatalf("%s attempted with nil want type", test.name)
	}
	gotV := reflect.New(gotT)
	gotP := gotV.Interface()

	// create a new pointer to a the same type as test.want,
	// and set its data to match test.want
	wantT := reflect.TypeOf(test.want)
	wantV := reflect.New(wantT)
	wantV.Elem().Set(reflect.ValueOf(test.want))
	wantP := wantV.Interface()

	err := json.Unmarshal([]byte(test.inputString), gotP)
	gotError := err != nil
	if gotError != test.wantError {
		t.Fatalf("%s json.Unmarshal returned error? %v (%s)", test.name, gotError, err)
	}

	if !reflect.DeepEqual(gotP, wantP) {
		t.Errorf("%s json.Unmarshal got\n%#v, want\n%#v", test.name, gotP, wantP)
	}
}

// jsonUnmarshalTestCases is a collection of jsonUnmarshalTestCases tests.
type jsonUnmarshalTestCases []jsonUnmarshalTestCase

// test runs the test defined for each item in the collection.
func (tests jsonUnmarshalTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}
