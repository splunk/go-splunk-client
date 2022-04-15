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

type testBool struct {
	Value Explicit[bool] `values:",omitempty"`
}

func TestBool_UnmarshalJSON(t *testing.T) {
	tests := jsonUnmarshalTestCases{
		{
			name:        "empty",
			inputString: `{}`,
			want:        testBool{},
		},
		{
			name:        "zero",
			inputString: `{"value":false}`,
			want:        testBool{Value: NewExplicit(false)},
		},
		{
			name:        "non-zero",
			inputString: `{"value":true}`,
			want:        testBool{Value: NewExplicit(true)},
		},
	}

	tests.test(t)
}

func TestBool_SetURLValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			name:  "implicit zero",
			input: testBool{},
			want:  url.Values{},
		},
		{
			name: "explicit zero",
			input: testBool{
				Value: NewExplicit(false),
			},
			want: url.Values{"Value": []string{"false"}},
		},
		{
			name: "non-zero",
			input: testBool{
				Value: NewExplicit(true),
			},
			want: url.Values{"Value": []string{"true"}},
		},
	}

	tests.test(t)
}
