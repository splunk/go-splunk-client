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

package entry

import (
	"github.com/splunk/go-sdk/pkg/attributes"
	"github.com/splunk/go-sdk/pkg/client"
)

// SAMLGroupContent defines the content for a SAMLGroup.
type SAMLGroupContent struct {
	Roles attributes.Strings `json:"roles" url:"roles"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Content
}

// SAMLGroup defines a SAML group mapping.
type SAMLGroup struct {
	client.Title     `json:"name" url:"name"`
	SAMLGroupContent `json:"content"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.GlobalNamespace
	// This endpoint returns a 400 if unable to find the given SAML Group.
	client.Endpoint `endpoint:"admin/SAML-groups,notfound:400"`
}
