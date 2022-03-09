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

func TestBool_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        struct{ Value Bool }{},
		},
		{
			name:        "zero",
			inputString: `{"value":false}`,
			want:        struct{ Value Bool }{Bool{explicit: true}},
		},
		{
			name:        "non-zero",
			inputString: `{"value":true}`,
			want:        struct{ Value Bool }{Bool{value: true, explicit: true}},
		},
	}

	tests.test(t)
}

func TestBool_EncodeValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			name:  "implicit zero",
			input: struct{ Value Bool }{},
			want:  url.Values{},
		},
		{
			name:  "explicit zero",
			input: struct{ Value Bool }{Bool{explicit: true}},
			want:  url.Values{"Value": []string{"false"}},
		},
		{
			name:  "non-zero",
			input: struct{ Value Bool }{Value: Bool{value: true}},
			want:  url.Values{"Value": []string{"true"}},
		},
	}

	tests.test(t)
}
