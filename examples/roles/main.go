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

	if err := client.Create(c, entry.Role{
		ID: client.ID{
			IDFields: client.IDFields{
				Title: "new_role",
			},
		},
		RoleContent: entry.RoleContent{
			SrchDiskQuota: attributes.NewInt(1),
			Capabilities:  attributes.NewStrings("search", "rtsearch"),
		},
	}); err != nil {
		log.Fatalf("unable to create role: %s", err)
	}

	createdRole := entry.Role{}
	if err := createdRole.SetIDFromURL("https://localhost:8089/services/authorization/roles/new_role"); err != nil {
		log.Fatalf("unable to set ID from URL: %s", err)
	}

	if err := client.Read(c, &createdRole); err != nil {
		log.Fatalf("unable to read role: %s", err)
	}
	fmt.Printf("capabilities: %v\n", createdRole.Capabilities)
	fmt.Printf("created role: %#v\n", createdRole)

	// here we explicitly set SrchDiskQuota to 0
	updateRole := entry.Role{
		ID: client.ID{
			IDFields: client.IDFields{
				Title: "new_role",
			},
		},
	}
	updateRole.SrchDiskQuota.Set(0)
	updateRole.Capabilities.Set()
	if err := client.Update(c, updateRole); err != nil {
		log.Fatalf("unable to update role: %s", err)
	}

	if err := client.Read(c, &createdRole); err != nil {
		log.Fatalf("unable to read role: %s", err)
	}
	fmt.Printf("capabilities: %v\n", createdRole.Capabilities)

	if err := client.Delete(c, createdRole); err != nil {
		log.Fatalf("unable to delete role: %s", err)
	}
}
