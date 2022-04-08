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
	"testing"
)

type testStrings struct {
	Values Explicit[[]string] `values:",omitempty"`
}

func TestStrings_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "absent",
			inputString: "{}",
			want:        testStrings{},
		},
		{
			name:        "empty list",
			inputString: `{"values":[]}`,
			want: testStrings{
				Values: NewExplicit([]string{}),
			},
		},
		{
			name:        "populated list",
			inputString: `{"values":["one","two"]}`,
			want: testStrings{
				Values: NewExplicit([]string{"one", "two"}),
			},
		},
	}

	tests.test(t)
}

func TestStrings_SetURLValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			name:  "zero value",
			input: testStrings{},
			want:  url.Values{},
		},
		{
			name: "explicitly empty",
			input: testStrings{
				Values: NewExplicit([]string{}),
			},
			want: url.Values{"Values": []string{""}},
		},
		{
			name: "implicit values",
			input: testStrings{
				Values: NewExplicit([]string{"one", "two"}),
			},
			want: url.Values{"Values": []string{"one", "two"}},
		},
	}

	tests.test(t)
}
