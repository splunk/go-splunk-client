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

package values_test

import (
	"fmt"

	"github.com/splunk/go-splunk-client/pkg/values"
)

type Action struct {
	Name       string `values:"-"`
	Enabled    bool
	Parameters map[string]string
}

func (a Action) GetURLKey(parentKey, childKey string) (string, error) {
	// the key for Action is <parentKey>.<a.Name>
	// this makes Enabled a value at this path,
	// and Parameters at <parentKey>.<a.Name>.<Parameter.key>
	return fmt.Sprint(parentKey, ".", a.Name), nil
}

type Search struct {
	Name    string   `values:"name"`
	Actions []Action `values:"action"`
}

func Example_encodeSearch() {
	search := Search{
		Name: "my_search",
		Actions: []Action{
			{
				Name:    "email",
				Enabled: true,
				Parameters: map[string]string{
					"subject":   "Something happened!",
					"recipient": "joeuser@example.com",
				},
			},
		},
	}

	v, _ := values.Encode(search)
	fmt.Println(v.Encode())
	// Output: action.email=true&action.email.recipient=joeuser%40example.com&action.email.subject=Something+happened%21&name=my_search
}
