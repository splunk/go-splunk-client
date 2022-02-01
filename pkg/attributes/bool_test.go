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
			"empty",
			`{}`,
			struct{ Value Bool }{},
		},
		{
			"zero",
			`{"value":false}`,
			struct{ Value Bool }{Bool{explicit: true}},
		},
		{
			"non-zero",
			`{"value":true}`,
			struct{ Value Bool }{Bool{value: true, explicit: true}},
		},
	}

	tests.test(t)
}

func TestBool_EncodeValues(t *testing.T) {
	tests := queryValuesTestCases{
		{
			"implicit zero",
			struct{ Value Bool }{},
			url.Values{},
		},
		{
			"explicit zero",
			struct{ Value Bool }{Bool{explicit: true}},
			url.Values{"Value": []string{"false"}},
		},
		{
			"non-zero",
			struct{ Value Bool }{Value: Bool{value: true}},
			url.Values{"Value": []string{"true"}},
		},
	}

	tests.test(t)
}
