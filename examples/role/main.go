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

	"github.com/splunk/go-sdk/pkg/client"
)

func main() {
	c := &client.Client{
		URL:                   "https://localhost:8089",
		Username:              "admin",
		Password:              "changedit",
		TLSInsecureSkipVerify: true,
	}
	c.Login()

	roles := client.Roles{}
	if err := client.RefreshCollection(&roles, c); err != nil {
		log.Fatalf("unable to refresh roles: %s", err)
	}
	fmt.Printf("%#v\n", roles)
	for _, role := range roles.Entries {
		fmt.Printf("%s\n", role.Title())
	}

	role, ok := roles.EntryWithTitle("user")
	if !ok {
		log.Fatalf("unable to get entry")
	}
	fmt.Printf("%#v\n", role)

	newRole := client.NewRoleWithTitle("admin")
	if err := client.RefreshEntry(newRole, c); err != nil {
		log.Fatalf("unable to refresh: %s", err)
	}
	fmt.Printf("%#v\n", newRole)

	// role := client.NewRoleWithTitle("admin")
	// fmt.Printf("pre:\n%#v\n", role)
	// if err := client.Read(role, c); err != nil {
	// 	log.Fatalf("unable to read role: %s", err)
	// }

	// fmt.Printf("post:\n%#v\n", role)
}
