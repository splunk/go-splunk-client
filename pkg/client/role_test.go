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

import (
	"net/http"
	"testing"

	"github.com/splunk/go-sdk/pkg/internal"
)

func TestRole_UrlValues(t *testing.T) {
	tests := urlValuesTestCases{
		{
			name:           "empty role",
			input:          Role{},
			valuesTestFunc: internal.TestURLValuesEncoded("capabilities="),
		},
	}

	tests.test(t)
}

func TestRole_requestForRead(t *testing.T) {
	tests := readerTestCases{
		{
			name:        "missing title",
			inputClient: &dummyAuthenticatedClient,
			inputReader: &Role{},
			wantError:   true,
		},
		{
			name:        "basic",
			inputClient: &dummyAuthenticatedClient,
			inputReader: NewRoleWithTitle("testrole"),
			wantError:   false,
			requestTestFunc: internal.ComposeTestRequests(
				internal.TestRequestMethod(http.MethodPost),
				internal.TestRequestHasAuth(),
				internal.TestRequestURL("https://localhost:8089/services/authorization/roles/testrole"),
				internal.TestRequestBody("capabilities="),
			),
		},
	}

	tests.test(t)
}
