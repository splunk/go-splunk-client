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

func TestEndpoint_getEndpointConfig(t *testing.T) {
	tests := []struct {
		name               string
		input              endpointConfigGetter
		wantError          bool
		wantErrorCode      ErrorCode
		wantEndpointConfig endpointConfig
	}{
		{
			"Endpoint field missing service tag",
			struct {
				Endpoint
			}{},
			true,
			ErrorEndpoint,
			endpointConfig{},
		},
		{
			"path only",
			struct {
				Endpoint `endpoint:"test/endpoint"`
			}{},
			false,
			ErrorUndefined,
			endpointConfig{
				path:         "test/endpoint",
				codeNotFound: 404,
			},
		},
		{
			"path and notfound",
			struct {
				Endpoint `endpoint:"test/endpoint,notfound:400"`
			}{},
			false,
			ErrorUndefined,
			endpointConfig{
				path:         "test/endpoint",
				codeNotFound: 400,
			},
		},
	}

	s := Endpoint{}
	for _, test := range tests {
		gotEndpointConfig, err := s.getEndpointConfig(test.input)

		testError(test.name, err, test.wantError, test.wantErrorCode, t)

		if gotEndpointConfig != test.wantEndpointConfig {
			t.Errorf("%s getEndpointConfig got:\n%#v, want\n%#v", test.name, gotEndpointConfig, test.wantEndpointConfig)
		}
	}
}
