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

package client

import (
	"fmt"

	"github.com/splunk/go-sdk/pkg/internal"
)

type RoleAttributes struct {
	Capabilities Capabilities `url:"capabilities"`
}

// Role represents a Splunk role.
type Role struct {
	internal.Name
	Attributes RoleAttributes `json:"content"`
}

func (r Role) entryCollection() EntryCollection {
	return &Roles{}
}

// validate returns an error if Role is invalid. It is invalid if it:
// * has an empty name
func (r Role) validate() error {
	if r.Name.Value == "" {
		return fmt.Errorf("invalid role, has no title")
	}

	return nil
}

// NewRoleWithTitle returns a new Role with the given title.
func NewRoleWithTitle(t string) *Role {
	return &Role{
		Name: internal.Name{Value: t},
	}
}
