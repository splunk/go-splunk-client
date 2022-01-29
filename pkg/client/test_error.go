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

// testError performs a test against a given error. wantCode is only tested if err
// is not nil.
func testError(name string, err error, wantError bool, wantCode ErrorCode, t *testing.T) {
	gotError := err != nil

	if gotError != wantError {
		t.Errorf("%s returned error? %v", name, wantError)
	}

	if gotError {
		clientErr := err.(Error)
		if clientErr.Code != wantCode {
			t.Errorf("%s returned error code %d, want %d", name, clientErr.Code, wantCode)
		}
	}
}
