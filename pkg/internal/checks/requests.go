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

package checks

import (
	"net/http"
	"reflect"
	"testing"
)

// CheckRequestFunc functions perform a check against an http.Request.
type CheckRequestFunc func(*http.Request, *testing.T)

// ComposeCheckRequestFunc returns a new CheckRequestFunc from an arbitrary number of other
// CheckRequestFunc functions.
func ComposeCheckRequestFunc(checks ...CheckRequestFunc) CheckRequestFunc {
	return func(r *http.Request, t *testing.T) {
		for _, check := range checks {
			check(r, t)
		}
	}
}

// CheckRequestHeaderKeyValue checks that an http.Request's header has a given value for
// a given key.
func CheckRequestHeaderKeyValue(key string, value ...string) CheckRequestFunc {
	return func(r *http.Request, t *testing.T) {
		if r.Header == nil {
			t.Errorf("CheckRequestHeaderKeyValue: Header not set")
			return
		}

		got, ok := r.Header[key]
		if !ok {
			t.Errorf("CheckRequestHeaderKeyValue: Key %s not set", key)
			return
		}

		if !reflect.DeepEqual(got, value) {
			t.Errorf("CheckRequestHeaderKeyValue: Key %s = %#v, want %#v", key, got, value)
			return
		}
	}
}
