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

	createdUser, err := client.Create(c, entry.User{
		Title: "newuser",
		UserContent: entry.UserContent{
			Password: attributes.NewString("changedit"),
			RealName: attributes.NewString("New User"),
			Roles:    attributes.NewStrings("user"),
		},
	})
	if err != nil {
		log.Fatalf("unable to create user: %s", err)
	}
	fmt.Printf("created user: %s\n", createdUser.Title)
	fmt.Printf("  real name: %s\n", createdUser.RealName)
	fmt.Printf("  roles: %s\n", createdUser.Roles)

	readUser, err := client.Read(c, entry.User{Title: "newuser"})
	if err != nil {
		log.Fatalf("unable to read user: %s", err)
	}
	fmt.Printf("read user: %s\n", readUser.Title)
	fmt.Printf("  real name: %s\n", readUser.RealName)
	fmt.Printf("  roles: %s\n", readUser.Roles)

	createdUser.RealName = attributes.NewString("Updated User")
	updatedUser, err := client.Update(c, createdUser)
	if err != nil {
		log.Fatalf("unable to update user: %s", err)
	}
	fmt.Printf("updated user: %s\n", updatedUser.Title)
	fmt.Printf("  real name: %s\n", updatedUser.RealName)
	fmt.Printf("  roles: %s\n", updatedUser.Roles)

	remainingUsers, err := client.Delete(c, createdUser)
	if err != nil {
		log.Fatalf("unable to delete user: %s", err)
	}
	for _, remainingUser := range remainingUsers {
		fmt.Printf("remaining user: %s\n", remainingUser.Title)
		fmt.Printf("  real name: %s\n", remainingUser.RealName)
		fmt.Printf("  roles: %s\n", remainingUser.Roles)
	}

	listedUsers, err := client.List(c, entry.User{})
	if err != nil {
		log.Fatalf("unable to list users: %s", err)
	}
	for _, listedUser := range listedUsers {
		fmt.Printf("listed user: %s\n", listedUser.Title)
		fmt.Printf("  real name: %s\n", listedUser.RealName)
		fmt.Printf("  roles: %s\n", listedUser.Roles)
	}
}
