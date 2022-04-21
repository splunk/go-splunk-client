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

type testInt struct {
	Value Explicit[int] `values:",omitzero"`
}

func TestInt_UnmarshalJSON(t *testing.T) {
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "empty",
			InputString: `{}`,
			Want:        testInt{},
		},
		{
			Name:        "zero",
			InputString: `{"value":0}`,
			Want:        testInt{Value: NewExplicit(0)},
		},
		{
			Name:        "non-zero",
			InputString: `{"value":1}`,
			Want:        testInt{Value: NewExplicit(1)},
		},
	}

	tests.Test(t)
}

func TestInt_SetURLValues(t *testing.T) {
	tests := checks.QueryValuesTestCases{
		{
			Name:  "implicit zero",
			Input: testInt{},
			Want:  url.Values{},
		},
		{
			Name:  "explicit zero",
			Input: testInt{Value: NewExplicit(0)},
			Want:  url.Values{"Value": []string{"0"}},
		},
		{
			Name:  "non-zero",
			Input: testInt{Value: NewExplicit(1)},
			Want:  url.Values{"Value": []string{"1"}},
		},
	}

	tests.Test(t)
}
