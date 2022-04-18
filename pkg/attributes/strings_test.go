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

	"github.com/splunk/go-splunk-client/pkg/internal/checks"
)

type testStrings struct {
	Values Explicit[[]string] `values:",omitempty"`
}

func TestStrings_UnmarshalJSON(t *testing.T) {
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "absent",
			InputString: "{}",
			Want:        testStrings{},
		},
		{
			Name:        "empty list",
			InputString: `{"values":[]}`,
			Want: testStrings{
				Values: NewExplicit([]string{}),
			},
		},
		{
			Name:        "populated list",
			InputString: `{"values":["one","two"]}`,
			Want: testStrings{
				Values: NewExplicit([]string{"one", "two"}),
			},
		},
	}

	tests.Test(t)
}

func TestStrings_SetURLValues(t *testing.T) {
	tests := checks.QueryValuesTestCases{
		{
			Name:  "zero value",
			Input: testStrings{},
			Want:  url.Values{},
		},
		{
			Name: "explicitly empty",
			Input: testStrings{
				Values: NewExplicit([]string{}),
			},
			Want: url.Values{"Values": []string{""}},
		},
		{
			Name: "implicit values",
			Input: testStrings{
				Values: NewExplicit([]string{"one", "two"}),
			},
			Want: url.Values{"Values": []string{"one", "two"}},
		},
	}

	tests.Test(t)
}
