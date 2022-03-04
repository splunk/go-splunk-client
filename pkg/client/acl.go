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
	"github.com/splunk/go-sdk/pkg/attributes"
	"github.com/splunk/go-sdk/pkg/internal/paths"
)

// Permissions represents the read/write permissions of a Splunk object.
type Permissions struct {
	Read  attributes.Strings `json:"read" url:"perms.read"`
	Write attributes.Strings `json:"write" url:"perms.write"`
}

// AccessController is the interface describing types that implement ACLPath and ACLValues.
type AccessController interface {
	ACLPath() string
	ACLValues() ACL
}

// EntryAccessControler is the interface describing types that implement both Entry and AccessController
// interfaces.
type EntryAccessController interface {
	Entry
	AccessController
}

// entryACLPath returns the path for an EntryAccessController's ACL.
func entryACLPath(entry EntryAccessController) (string, error) {
	if entry.Title() == "" {
		return "", wrapError(ErrorACL, nil, "attempted to get EntryACLPath for Entry with empty Title")
	}

	eP, err := entryPath(entry)
	if err != nil {
		return "", err
	}

	return paths.Join(eP, entry.ACLPath()), nil
}

// ACL represents the ACL of a Splunk object.
type ACL struct {
	Permissions `json:"perms"`
	Owner       attributes.String `json:"owner" url:"owner"`
	Sharing     attributes.String `json:"sharing" url:"sharing"`
}

// ACLPath returns the relative path for ACL management.
func (acl ACL) ACLPath() string {
	return "acl"
}

// ACLValues returns the object that should be used by EncodeValues for ACL management.
func (acl ACL) ACLValues() ACL {
	return acl
}
