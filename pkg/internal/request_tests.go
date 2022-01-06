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

package internal

import (
	"net/http"
	"testing"
)

// TestRequestFunc describes functions that perform tests against http.Request objects.
type TestRequestFunc func(name string, r *http.Request, t *testing.T)

// ComposeTestRequests composes a new TestRequestFunc from any number of other TestRequestFunc items.
func ComposeTestRequests(tests ...TestRequestFunc) TestRequestFunc {
	return func(name string, r *http.Request, t *testing.T) {
		for _, test := range tests {
			test(name, r, t)
		}
	}
}

// TestRequestMethod tests that an http.Request has a given Method.
func TestRequestMethod(wantMethod string) TestRequestFunc {
	return func(name string, r *http.Request, t *testing.T) {
		if r.Method != wantMethod {
			t.Errorf("%s TestRequestMethod: got %s, want %s", name, r.Method, wantMethod)
		}
	}
}

// TestRequestURL tests that an http.Request has a given URL.
func TestRequestURL(wantURL string) TestRequestFunc {
	return func(name string, r *http.Request, t *testing.T) {
		gotURL := r.URL.String()
		if gotURL != wantURL {
			t.Errorf("%s TestRequestURL: got\n%s, want\n%s", name, gotURL, wantURL)
		}
	}
}

// TestRequestHasAuth tests that an http.Request has an Authorization header.
func TestRequestHasAuth() TestRequestFunc {
	return func(name string, r *http.Request, t *testing.T) {
		authFound := false
		if r.Header != nil {
			if _, ok := r.Header["Authorization"]; ok {
				authFound = true
			}
		}

		if !authFound {
			t.Errorf("%s TestRequestHasAuth: request lacking Authorization", name)
		}
	}
}

// TestRequestBody tests that an http.Request has the given Body content.
func TestRequestBody(wantBody string) TestRequestFunc {
	return func(name string, r *http.Request, t *testing.T) {
		gotBody := string(MustReadAll(r.Body))
		if gotBody != wantBody {
			t.Errorf("%s Body =\n%s, want\n%s", name, gotBody, wantBody)
		}
	}
}
