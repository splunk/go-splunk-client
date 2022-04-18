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

	newStanza := entry.Stanza{
		ID: client.ConfID{
			Namespace: client.Namespace{
				App:  "search",
				User: "nobody",
			},
			File:   "customfile",
			Stanza: "customstanza",
		},
		Content: entry.StanzaContent{
			Values: map[string]string{
				"customKeyA": "customValueA",
			},
		},
	}

	// create new stanza
	if err := client.Create(c, newStanza); err != nil {
		log.Fatalf("unable to create stanza: %s", err)
	}

	// read created stanza
	if err := client.Read(c, &newStanza); err != nil {
		log.Fatalf("unable to read stanza: %s", err)
	}
	fmt.Printf("created Stanza:\n%#v\n", newStanza)

	// update stanza
	newStanza.Content.Values = map[string]string{
		// customKeyA will untouched, as it won't be POSTed to the API
		// if a key needs to be cleared, it can be done so by setting its value to an empty string:
		// "customKeyA": "",

		// customKeyB will be added
		"customKeyB": "customValueB",
	}
	if err := client.Update(c, newStanza); err != nil {
		log.Fatalf("unable to update stanza: %s", err)
	}
	if err := client.Read(c, &newStanza); err != nil {
		log.Fatalf("unable to read stanza: %s", err)
	}
	fmt.Printf("updated Stanza:\n%#v\n", newStanza)

	// delete stanza
	if err := client.Delete(c, &newStanza); err != nil {
		log.Fatalf("unable to delete stanza: %s", err)
	}
}
