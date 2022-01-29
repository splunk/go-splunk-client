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

package client

import "testing"

func TestEndpoint_endpoint(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		wantError     bool
		wantErrorCode ErrorCode
		wantEndpoint  string
	}{
		{
			"non-struct",
			"this string is not a struct",
			true,
			ErrorEndpoint,
			"",
		},
		{
			"struct has no field named Endpoint",
			struct {
				path Endpoint `endpoint:"test/endpoint"`
			}{},
			true,
			ErrorEndpoint,
			"",
		},
		{
			"Endpoint field isn't Endpoint type",
			struct {
				Endpoint string `endpoint:"test/endpoint"`
			}{},
			true,
			ErrorEndpoint,
			"",
		},
		{
			"Endpoint field missing service tag",
			struct {
				Endpoint
			}{},
			true,
			ErrorEndpoint,
			"",
		},
		{
			"everything in place",
			struct {
				Endpoint `endpoint:"test/endpoint"`
			}{},
			false,
			ErrorUndefined,
			"test/endpoint",
		},
	}

	s := Endpoint{}
	for _, test := range tests {
		gotEndpoint, err := s.endpointPath(test.input)

		testError(test.name, err, test.wantError, test.wantErrorCode, t)

		if gotEndpoint != test.wantEndpoint {
			t.Errorf("%s got %s, want %s", test.name, gotEndpoint, test.wantEndpoint)
		}
	}
}
