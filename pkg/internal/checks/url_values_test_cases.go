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
	"net/url"
	"reflect"
	"testing"

	"github.com/splunk/go-splunk-client/pkg/values"
)

// queryValuesTestCase defines a test case for query.Values.
type queryValuesTestCase struct {
	name      string
	input     interface{}
	want      url.Values
	wantError bool
}

// test runs the test.
func (test queryValuesTestCase) test(t *testing.T) {
	got, err := values.Encode(test.input)
	gotError := err != nil

	if gotError != test.wantError {
		t.Errorf("%s values.Encode returned error? %v: %s", test.name, gotError, err)
	}

	if !reflect.DeepEqual(got, test.want) {
		t.Errorf("%s values.Encode got\n%#v, want\n%#v", test.name, got, test.want)
	}
}

// queryValuesTestCases is a collection of queryValuesTestCase tests.
type queryValuesTestCases []queryValuesTestCase

// test runs the test defined for each item in the collection.
func (tests queryValuesTestCases) test(t *testing.T) {
	for _, test := range tests {
		test.test(t)
	}
}
