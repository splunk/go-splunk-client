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

type testBool struct {
	Value Explicit[bool] `values:",omitempty"`
}

func TestBool_UnmarshalJSON(t *testing.T) {
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "empty",
			InputString: `{}`,
			Want:        testBool{},
		},
		{
			Name:        "zero",
			InputString: `{"value":false}`,
			Want:        testBool{Value: NewExplicit(false)},
		},
		{
			Name:        "non-zero",
			InputString: `{"value":true}`,
			Want:        testBool{Value: NewExplicit(true)},
		},
	}

	tests.Test(t)
}

func TestBool_SetURLValues(t *testing.T) {
	tests := checks.QueryValuesTestCases{
		{
			Name:  "implicit zero",
			Input: testBool{},
			Want:  url.Values{},
		},
		{
			Name: "explicit zero",
			Input: testBool{
				Value: NewExplicit(false),
			},
			Want: url.Values{"Value": []string{"false"}},
		},
		{
			Name: "non-zero",
			Input: testBool{
				Value: NewExplicit(true),
			},
			Want: url.Values{"Value": []string{"true"}},
		},
	}

	tests.Test(t)
}
