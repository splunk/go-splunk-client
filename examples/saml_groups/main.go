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

	err := client.Create(c, entry.SAMLGroup{
		ID: client.ID{
			Title: "new_saml_group",
		},
		SAMLGroupContent: entry.SAMLGroupContent{
			Roles: attributes.NewExplicit([]string{"admin"}),
		},
	})
	if err != nil {
		log.Fatalf("unable to create new SAML group: %s", err)
	}

	readSAMLGroup := entry.SAMLGroup{
		ID: client.ID{
			Title: "new_saml_group",
		},
	}
	if err := client.Read(c, &readSAMLGroup); err != nil {
		if clientErr, ok := err.(client.Error); ok {
			if clientErr.Code == client.ErrorNotFound {
				log.Fatalf("SAML group not found: %s", clientErr)
			}
		}
		log.Fatalf("unable to read SAML group: %s", err)
	}
	fmt.Printf("read SAML group: %#v\n", readSAMLGroup)

	if err := client.Delete(c, entry.SAMLGroup{
		ID: client.ID{
			Title: "new_saml_group",
		},
	}); err != nil {
		log.Fatalf("unable to delete SAML group: %s", err)
	}
}
