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

func TestString_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        struct{ Value String }{},
		},
		{
			name:        "empty",
			inputString: `{"value":""}`,
			want:        struct{ Value String }{String{explicit: true}},
		},
		{
			name:        "non-empty",
			inputString: `{"value":"this string is not empty"}`,
			want:        struct{ Value String }{String{value: "this string is not empty", explicit: true}},
		},
	}

	tests.test(t)
}

func TestString_EncodeValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			"implicit empty",
			struct{ Value String }{},
			url.Values{},
		},
		{
			"explicit empty",
			struct{ Value String }{String{explicit: true}},
			url.Values{"Value": []string{""}},
		},
		{
			"non-empty",
			struct{ Value String }{Value: String{value: "this string is not empty"}},
			url.Values{"Value": []string{"this string is not empty"}},
		},
	}

	tests.test(t)
}
