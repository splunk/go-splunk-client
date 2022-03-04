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

// UserContent defines the content of a User object.
type UserContent struct {
	// Read/Write
	DefaultApp            attributes.String  `url:"defaultApp"`
	Email                 attributes.String  `url:"email"`
	Password              attributes.String  `url:"password,omitempty"`
	RealName              attributes.String  `url:"realname"`
	RestartBackgroundJobs attributes.Bool    `url:"restart_background_jobs"`
	Roles                 attributes.Strings `url:"roles"`
	TZ                    attributes.String  `url:"tz"`

	// Read-only fields are populated by results returned by the Splunk API, but
	// are not settable by Create or Update operations.
	Capabilities attributes.Strings `url:"-"`
	Type         string             `url:"-"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Content
}

// User defines a Splunk user.
type User struct {
	client.ID
	UserContent `json:"content"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Endpoint `endpoint:"authentication/users"`
}
