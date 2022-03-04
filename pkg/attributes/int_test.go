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

func TestInt_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        struct{ Value Int }{},
		},
		{
			name:        "zero",
			inputString: `{"value":0}`,
			want:        struct{ Value Int }{Int{explicit: true}},
		},
		{
			name:        "non-zero",
			inputString: `{"value":1}`,
			want:        struct{ Value Int }{Int{value: 1, explicit: true}},
		},
	}

	tests.test(t)
}

func TestInt_EncodeValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			name:  "implicit zero",
			input: struct{ Value Int }{},
			want:  url.Values{},
		},
		{
			name:  "explicit zero",
			input: struct{ Value Int }{Int{explicit: true}},
			want:  url.Values{"Value": []string{"0"}},
		},
		{
			name:  "non-zero",
			input: struct{ Value Int }{Value: Int{value: 1}},
			want:  url.Values{"Value": []string{"1"}},
		},
	}

	tests.test(t)
}
