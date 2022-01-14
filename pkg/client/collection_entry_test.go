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
	"go-sdk/pkg/client/internal/checks"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestCollectionEntry_entryReadRequest(t *testing.T) {
	tests := []struct {
		name      string
		input     collectionEntry
		wantError bool
		checks    checks.CheckRequestFunc
	}{
		{
			"no title, globalNamespace",
			struct {
				service `service:"widgets"`
				globalNamespace
				Title
			}{},
			false,
			checks.ComposeCheckRequestFunc(
				checks.CheckRequestURL("https://localhost:8089/services/widgets"),
			),
		},
		{
			"has title, Namespace",
			struct {
				service `service:"widgets"`
				Namespace
				Title
			}{
				Namespace: Namespace{User: "testuser", App: "testapp"},
				Title:     "testwidget",
			},
			false,
			checks.ComposeCheckRequestFunc(
				checks.CheckRequestURL("https://localhost:8089/servicesNS/testuser/testapp/widgets/testwidget"),
			),
		},
	}

	c := &Client{
		URL:           "https://localhost:8089",
		Authenticator: &SessionKeyAuth{SessionKey: "fake-session-key"},
	}

	for _, test := range tests {
		r, err := entryReadRequest(c, test.input)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s entryReadRequest returned error? %v", test.name, gotError)
		}

		test.checks(r, t)
	}
}

func TestCollectionEntry_entryReadResponseEntries(t *testing.T) {
	tests := []struct {
		name            string
		inputStatusCode int
		inputBody       string
		wantError       bool
		wantEntries     []User
	}{
		{
			"non-200 OK",
			http.StatusNotFound,
			"",
			true,
			[]User(nil),
		},
		{
			"no entries",
			http.StatusOK,
			`{"entry":[]}`,
			false,
			[]User{},
		},
		{
			"one entry",
			http.StatusOK,
			`{
				"entry":[
					{
						"name":"admin",
						"content":{
							"defaultApp":"search"
						}
					}
				]
			}`,
			false,
			[]User{
				{Title: "admin", Attributes: UserAttributes{DefaultApp: "search"}},
			},
		},
		{
			"two enties",
			http.StatusOK,
			`{
				"entry":[
					{
						"name":"admin",
						"content":{
							"defaultApp":"search"
						}
					},
					{
						"name":"joeuser",
						"content":{
							"defaultApp":"launcher"
						}
					}
				]
			}`,
			false,
			[]User{
				{Title: "admin", Attributes: UserAttributes{DefaultApp: "search"}},
				{Title: "joeuser", Attributes: UserAttributes{DefaultApp: "launcher"}},
			},
		},
	}

	for _, test := range tests {
		r := &http.Response{
			StatusCode: test.inputStatusCode,
			Body:       io.NopCloser(strings.NewReader(test.inputBody)),
		}

		gotEntries, err := entryReadResponseEntries(User{}, r)
		gotError := err != nil

		if gotError != test.wantError {
			t.Errorf("%s entryReadResponseEntries returned error? %v", test.name, gotError)
		}

		if !reflect.DeepEqual(gotEntries, test.wantEntries) {
			t.Errorf("%s entryReadResponseEntries got\n%#v, want\n%#v", test.name, gotEntries, test.wantEntries)
		}
	}
}
