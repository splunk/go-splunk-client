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

	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/authenticators"
	"github.com/splunk/go-splunk-client/pkg/client"
	"github.com/splunk/go-splunk-client/pkg/entry"
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

	newSearch := entry.SavedSearch{
		ID: client.ID{
			Namespace: client.Namespace{
				User: "nobody",
				App:  "search",
			},
			Title: "mysavedsearch",
		},
		Content: entry.SavedSearchContent{
			Search: attributes.NewExplicit("| makeresults"),
			Actions: attributes.NamedParametersCollection{
				attributes.NamedParameters{
					Name:   "email",
					Status: attributes.NewExplicit("true"),
					Parameters: attributes.Parameters{
						"to":      "joeuser@example.com",
						"subject": "mysearch results",
					},
				},
			},
		},
	}

	// create search
	if err := c.Create(newSearch); err != nil {
		log.Fatalf("unable to create search: %s", err)
	}
	if err := c.Read(&newSearch); err != nil {
		log.Fatalf("unable to read search: %s", err)
	}
	fmt.Printf("created search:\n%#v\n", newSearch)

	// update search
	newSearch.Content.Actions = attributes.NamedParametersCollection{} // explicitly clear Actions, disabling all
	if err := c.Update(newSearch); err != nil {
		log.Fatalf("unable to update search: %s", err)
	}
	if err := c.Read(&newSearch); err != nil {
		log.Fatalf("unable to read search: %s", err)
	}
	fmt.Printf("updated search:\n%#v\n", newSearch)

	// delete search
	if err := c.Delete(newSearch); err != nil {
		log.Fatalf("unable to delete search: %s", err)
	}
}
