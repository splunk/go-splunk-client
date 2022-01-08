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
	"testing"
)

func TestMessageElement_Unmarshal(t *testing.T) {
	tests := xmlUnmarshalerTestCases{
		{
			input:            "<msg>Message</msg>",
			gotInterfacePtr:  &message{},
			wantInterfacePtr: &message{Value: "Message"},
			wantError:        false,
		},
		{
			input:            `<msg code="200OK">Message</msg>`,
			gotInterfacePtr:  &message{},
			wantInterfacePtr: &message{Code: "200OK", Value: "Message"},
			wantError:        false,
		},
	}

	tests.test(t)
}

func TestMessages_Unmarshal(t *testing.T) {
	type responseType struct {
		Messages messages
	}

	tests := xmlUnmarshalerTestCases{
		{
			input: `
			<response>
				<messages>
					<msg code="200OK"></msg>
					<msg code="401Unauthorized">Unauthorized</msg>
				</messages>
			</response>`,
			gotInterfacePtr: &responseType{},
			wantInterfacePtr: &responseType{
				Messages: messages{
					Items: []message{
						{Code: "200OK", Value: ""},
						{Code: "401Unauthorized", Value: "Unauthorized"},
					},
				},
			},
			wantError: false,
		},
	}

	tests.test(t)
}

func TestMessages_firstAndOnly(t *testing.T) {
	tests := []struct {
		name        string
		input       messages
		wantOk      bool
		wantMessage message
	}{
		{
			name:        "empty",
			input:       messages{},
			wantOk:      false,
			wantMessage: message{},
		},
		{
			name: "multiple",
			input: messages{
				Items: []message{
					{Value: "first"},
					{Value: "second"},
				},
			},
			wantOk:      false,
			wantMessage: message{},
		},
		{
			name: "exactly one",
			input: messages{
				Items: []message{
					{Value: "only"},
				},
			},
			wantOk:      true,
			wantMessage: message{Value: "only"},
		},
	}

	for _, test := range tests {
		gotMessage, gotOk := test.input.firstAndOnly()

		if gotOk != test.wantOk {
			t.Errorf("%s firstAndOnly returned ok? %v", test.name, gotOk)
		}

		if gotMessage != test.wantMessage {
			t.Errorf("%s firstAndOnly got\n%#v, want\n%#v", test.name, gotMessage, test.wantMessage)
		}
	}
}
