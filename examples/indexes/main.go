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
	c := client.Client{
		URL: "https://localhost:8089",
		Authenticator: &authenticators.Password{
			Username: "admin",
			Password: "changeme",
		},
		TLSInsecureSkipVerify: true,
	}

	newIndex := entry.Index{
		ID: client.ID{
			Title: "new_index",
		},
	}

	// create index
	if err := client.Create(&c, newIndex); err != nil {
		log.Fatalf("unable to create index: %s", err)
	}
	if err := client.Read(&c, &newIndex); err != nil {
		log.Fatalf("unable to read index: %s", err)
	}
	fmt.Printf("created index: %#v\n", newIndex)

	// update index
	newIndex.Content.FrozenTimePeriodInSecs = attributes.NewExplicit(86400)
	if err := client.Update(&c, newIndex); err != nil {
		log.Fatalf("unable to update index: %s", err)
	}
	if err := client.Read(&c, &newIndex); err != nil {
		log.Fatalf("unable to read index: %s", err)
	}
	fmt.Printf("updated index: %#v\n", newIndex)

	// delete index
	if err := client.Delete(&c, newIndex); err != nil {
		log.Fatalf("unable to delete index: %s", err)
	}
}
