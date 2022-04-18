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

package entry

import (
	"testing"

	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/go-splunk-client/pkg/internal/checks"
)

func TestStanza_UnmarshalJSON(t *testing.T) {
	tests := checks.JSONUnmarshalTestCases{
		{
			Name:        "empty",
			InputString: `{}`,
			Want:        Stanza{},
		},
		{
			Name:        "id",
			InputString: `{"id":"https://localhost:8089/services/configs/conf-testfile/teststanza"}`,
			Want: Stanza{
				ID: client.ConfID{
					Namespace: client.Namespace{},
					File:      "testfile",
					Stanza:    "teststanza",
				},
			},
			WantError: false,
		},
		{
			Name:        "disabled",
			InputString: `{"content":{"disabled":true}}`,
			Want: Stanza{
				Content: StanzaContent{
					Disabled: attributes.NewExplicit(true),
				},
			},
		},
		{
			Name:        "disabled non-bool",
			InputString: `{"content":{"disabled":"true"}}`,
			Want:        Stanza{},
			WantError:   true,
		},
		{
			Name:        "ignore eai:",
			InputString: `{"content":{"eai:whatever":"anyvalue"}}`,
			Want:        Stanza{},
		},
		{
			Name:        "stanza values",
			InputString: `{"content":{"keyA":"valueA","keyB":"valueB"}}`,
			Want: Stanza{
				Content: StanzaContent{
					Values: map[string]string{
						"keyA": "valueA",
						"keyB": "valueB",
					},
				},
			},
		},
	}

	tests.Test(t)
}

func TestStanza_EncodeURLValues(t *testing.T) {
	tests := checks.QueryValuesTestCases{
		{
			Name:      "empty",
			Input:     Stanza{},
			WantError: true,
		},
		{
			Name: "has id",
			Input: Stanza{
				ID: client.ConfID{
					Stanza: "teststanza",
				},
			},
			Want: map[string][]string{
				"name": {"teststanza"},
			},
		},
		{
			Name: "disabled=false",
			Input: Stanza{
				ID: client.ConfID{
					Stanza: "teststanza",
				},
				Content: StanzaContent{
					Disabled: attributes.NewExplicit(false),
				},
			},
			Want: map[string][]string{
				"name":     {"teststanza"},
				"disabled": {"false"},
			},
		},
		{
			Name: "stanza values",
			Input: Stanza{
				ID: client.ConfID{
					Stanza: "teststanza",
				},
				Content: StanzaContent{
					Values: map[string]string{
						"keyA": "valueA",
						"keyB": "valueB",
					},
				},
			},
			Want: map[string][]string{
				"name": {"teststanza"},
				"keyA": {"valueA"},
				"keyB": {"valueB"},
			},
		},
	}

	tests.Test(t)
}
