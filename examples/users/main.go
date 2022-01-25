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
	"github.com/splunk/go-sdk/pkg/collections"
	"github.com/splunk/go-sdk/pkg/errors"
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

	// List
	users, err := client.CollectionList(c, collections.User{Title: "admin"})
	if err != nil {
		if clientErr, ok := err.(errors.Error); ok {
			if clientErr.Kind == errors.ErrorHTTPNotFound {
				log.Fatalf("not found")
			}
		}

		log.Fatalf("error: %s", err)
	}
	fmt.Printf("list users:\n")
	for _, user := range users {
		fmt.Printf("  %s: %s\n", user.Title, user.Roles)
	}

	// Create
	createUser := collections.User{
		Title: "newuser",
		UserContent: collections.UserContent{
			Email:    "newuser@example.com",
			Password: "newpassword",
			RealName: "New User",
			Roles:    attributes.Roles{"user"},
		},
	}
	createdUser, err := client.CollectionCreate(c, createUser)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("created user:\n")
	fmt.Printf("  %s: %s\n", createUser.Title, createdUser.Roles)

	// List after user creation
	users, err = client.CollectionList(c, collections.User{})
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("list users after creation:\n")
	for _, user := range users {
		fmt.Printf("  %s: %s\n", user.Title, user.Roles)
	}

	// Update
	createUser.Roles = attributes.Roles{"admin"}
	updatedUser, err := client.CollectionUpdate(c, createUser)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("updated user:\n")
	fmt.Printf("  %s: %s\n", updatedUser.Title, updatedUser.Roles)

	// Delete
	usersAfterDelete, err := client.CollectionDelete(c, createUser)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("users after deletion:\n")
	for _, user := range usersAfterDelete {
		fmt.Printf("  %s: %s\n", user.Title, user.Roles)
	}
}
