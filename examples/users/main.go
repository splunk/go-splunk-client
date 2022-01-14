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
		TLSInsecureSkipVerify: true,
		Authenticator: &client.PasswordAuth{
			Username:     "admin",
			Password:     "changeme",
			UseBasicAuth: true,
		},
	}

	// user := client.User{Title: "admin"}
	// c.ReadEntry(u) // should ensure only one result

	users, err := client.ReadCollection(c, client.User{})
	if err != nil {
		log.Fatalf("unable to read users: %s", err)
	}

	for _, user := range users {
		fmt.Printf("user: %s\n", user.Title)
	}

	roles, err := client.ReadCollection(c, client.Role{Namespace: client.Namespace{User: "-", App: "-"}, Title: "admin"})
	if err != nil {
		log.Fatalf("unable to read roles: %s", err)
	}

	for _, role := range roles {
		fmt.Printf("role: %s: %v\n", role.Title, role.Attributes.Capabilities)
	}
}
