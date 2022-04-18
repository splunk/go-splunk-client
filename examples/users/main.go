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

	if err := c.Create(entry.User{
		ID: client.ID{
			Title: "newuser",
		},
		Content: entry.UserContent{
			Password: attributes.NewExplicit("changedit"),
			RealName: attributes.NewExplicit("New User"),
			Roles:    attributes.NewExplicit([]string{"user"}),
		},
	}); err != nil {
		log.Fatalf("unable to create user: %s", err)
	}

	createdUser := entry.User{
		ID: client.ID{
			Title: "newuser",
		},
	}

	if err := c.Read(&createdUser); err != nil {
		log.Fatalf("unable to read user: %s", err)
	}
	fmt.Printf("read user: %s\n", createdUser.ID)
	fmt.Printf("  real name: %s\n", createdUser.Content.RealName)
	fmt.Printf("  roles: %s\n", createdUser.Content.Roles)

	createdUser.Content.RealName = attributes.NewExplicit("Updated User")
	if err := c.Update(createdUser); err != nil {
		log.Fatalf("unable to update user: %s", err)
	}

	if err := c.Delete(createdUser); err != nil {
		log.Fatalf("unable to delete user: %s", err)
	}

	listedUsers := []entry.User{}
	err := c.List(&listedUsers)
	if err != nil {
		log.Fatalf("unable to list users: %s", err)
	}
	for _, listedUser := range listedUsers {
		fmt.Printf("listed user: %s\n", listedUser.ID)
		fmt.Printf("  real name: %s\n", listedUser.Content.RealName)
		fmt.Printf("  roles: %s\n", listedUser.Content.Roles)
	}
}
