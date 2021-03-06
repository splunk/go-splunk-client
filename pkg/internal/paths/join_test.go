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

package paths

import "testing"

func Test_Join(t *testing.T) {
	tests := []struct {
		name       string
		inputPaths []string
		want       string
	}{
		{
			"empty parts",
			[]string{"", "", ""},
			"//",
		},
		{
			"trimmed leading/trailing slashes",
			[]string{"/one/", "/two", "three/"},
			"one/two/three",
		},
	}

	for _, test := range tests {
		got := Join(test.inputPaths...)

		if got != test.want {
			t.Errorf("%s Join got %s, want %s", test.name, got, test.want)
		}
	}
}
