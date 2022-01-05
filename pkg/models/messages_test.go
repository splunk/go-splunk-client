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

package models

import (
	"testing"

	"github.com/splunk/go-sdk/pkg/internal/testinghelpers"
)

func TestMessageElement_Unmarshal(t *testing.T) {
	tests := testinghelpers.XMLUnmarshalerTestCases{
		{
			Input:            "<msg>Message</msg>",
			GotInterfacePtr:  &MessageElement{},
			WantInterfacePtr: &MessageElement{Message: "Message"},
			WantError:        false,
		},
		{
			Input:            `<msg code="200OK">Message</msg>`,
			GotInterfacePtr:  &MessageElement{},
			WantInterfacePtr: &MessageElement{Code: "200OK", Message: "Message"},
			WantError:        false,
		},
	}

	tests.Test(t)
}

func TestMessagesElement_Unmarshal(t *testing.T) {
	tests := testinghelpers.XMLUnmarshalerTestCases{
		{
			Input: `
				<messages>
					<msg code="200OK"></msg>
					<msg code="401Unauthorized">Unauthorized</msg>
				</messages>`,
			GotInterfacePtr: &MessagesElement{},
			WantInterfacePtr: &MessagesElement{
				MessageElements: []MessageElement{
					{Code: "200OK", Message: ""},
					{Code: "401Unauthorized", Message: "Unauthorized"},
				},
			},
			WantError: false,
		},
	}

	tests.Test(t)
}
