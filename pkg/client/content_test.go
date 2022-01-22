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
	"reflect"
	"testing"
)

func TestContenter_Content(t *testing.T) {
	type InnerContent struct {
		data string
		Content
	}

	type InnerData struct {
		data string
	}

	type hasInnerContent struct {
		InnerContent
	}

	type hasOuterContent struct {
		InnerData
		Content
	}

	type hasBothContent struct {
		InnerContent
		Content
	}

	tests := []struct {
		name        string
		input       ContentGetter
		wantContent ContentGetter
	}{
		{
			"inner content",
			hasInnerContent{},
			InnerContent{},
		},
		{
			"outer content",
			hasOuterContent{},
			hasOuterContent{},
		},
		{
			"inner and outer content",
			hasBothContent{},
			hasBothContent{},
		},
	}

	for _, test := range tests {
		gotContent := test.input.GetContent(test.input)

		if !reflect.DeepEqual(gotContent, test.wantContent) {
			t.Errorf("%s content() got\n%#v, want\n%#v", test.name, gotContent, test.wantContent)
		}
	}
}
