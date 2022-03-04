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

import "testing"

func Test_dottedParameterNameParts(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantName      string
		wantParamName string
	}{
		{
			"empty",
			"",
			"",
			"",
		},
		{
			"name only",
			"testname",
			"testname",
			"",
		},
		{
			"name and param name",
			"testname.testparam",
			"testname",
			"testparam",
		},
		{
			"name and dotted param name",
			"testname.testparamA.testparamB",
			"testname",
			"testparamA.testparamB",
		},
	}

	for _, test := range tests {
		gotName, gotParamName := dottedParameterNameParts(test.input)

		if (gotName != test.wantName) || (gotParamName != test.wantParamName) {
			t.Errorf("%s Test_dottedParameterNameParts() got\n(%q, %q), want\n(%q, %q)", test.name, gotName, gotParamName, test.wantName, test.wantParamName)
		}
	}
}
