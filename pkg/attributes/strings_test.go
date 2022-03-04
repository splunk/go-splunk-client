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

func TestStrings_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "absent",
			inputString: "{}",
			want:        struct{ Values Strings }{},
		},
		{
			name:        "empty list",
			inputString: `{"values":[]}`,
			want: struct{ Values Strings }{
				Strings{
					values:   []string{},
					explicit: true,
				},
			},
		},
		{
			name:        "populated list",
			inputString: `{"values":["one","two"]}`,
			want: struct{ Values Strings }{
				Strings{
					values:   []string{"one", "two"},
					explicit: true,
				},
			},
		},
	}

	tests.test(t)
}

func TestStrings_EncodeValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			"zero value",
			struct{ Value Strings }{},
			url.Values{},
		},
		{
			"explicitly empty",
			struct{ Value Strings }{
				Strings{
					explicit: true,
				},
			},
			url.Values{"Value": []string{""}},
		},
		{
			"implicit values",
			struct{ Value Strings }{
				Strings{
					values: []string{"one", "two"},
				},
			},
			url.Values{"Value": []string{"one", "two"}},
		},
	}

	tests.test(t)
}
