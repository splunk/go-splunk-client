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

package main

import (
	"fmt"
	"log"

	"github.com/splunk/go-sdk/pkg/attributes"
	"github.com/splunk/go-sdk/pkg/authenticators"
	"github.com/splunk/go-sdk/pkg/client"
	"github.com/splunk/go-sdk/pkg/entry"
)

func main() {
	c := &client.Client{
		URL: "https://localhost:8089",
		Authenticator: &authenticators.Password{
			Username: "admin",
			Password: "changeme",
		},
		TLSInsecureSkipVerify: true,
	}

	ssCreated := entry.SavedSearch{
		ID: client.ID{IDFields: client.IDFields{
			// User:  "nobody",
			// App:   "system",
			Title: "new_savedsearch",
		}},
		ACL: client.ACL{
			Permissions: client.Permissions{
				Read:  attributes.NewStrings(),
				Write: attributes.NewStrings(),
			},
			Owner:   attributes.NewString("nobody"),
			Sharing: attributes.NewString("app"),
		},
		SavedSearchContent: entry.SavedSearchContent{
			Search: attributes.NewString("index=_internal"),
			Actions: entry.SavedSearchActions{
				{
					Name: "testaction",
					Parameters: map[string]string{
						"testparam": "testvalue",
					},
				},
			},
			// Args: entry.SavedSearchArgs{
			// 	"test.param1": "testvalue",
			// },
		},
	}
	if err := client.Create(c, ssCreated); err != nil {
		log.Fatalf("unable to create savedsearch: %s", err)
	}
	// if err := client.UpdateACL(c, ssCreated); err != nil {
	// 	log.Fatalf("unable to update ACL: %s", err)
	// }

	ssRead := entry.SavedSearch{
		ID: client.ID{IDFields: client.IDFields{Title: "new_savedsearch"}},
	}

	if err := client.Read(c, &ssRead); err != nil {
		log.Fatalf("unable to read savedsearch: %s", err)
	}
	fmt.Printf("read savedsearch: %s\n", ssRead.ID)
	for _, action := range ssRead.Actions {
		fmt.Printf("  action %s:\n", action.Name)
		fmt.Printf("    disabled? %v\n", action.Disabled)
		fmt.Printf("    params:\n")
		for key, value := range action.Parameters {
			fmt.Printf("      %s: %v\n", key, value)
		}
	}
	fmt.Printf("  args:\n")
	// for key, value := range ssRead.Args {
	// 	fmt.Printf("    %s: %s\n", key, value)
	// }
	// fmt.Printf("  %#v\n", ssRead)

	if err := client.Delete(c, ssCreated); err != nil {
		log.Fatalf("unable to delete savedsearch: %s", err)
	}
}
