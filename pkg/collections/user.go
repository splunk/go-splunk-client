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

package collections

import (
	"github.com/splunk/go-sdk/pkg/attributes"
	"github.com/splunk/go-sdk/pkg/client"
)

// UserContent defines the content of a User object.
type UserContent struct {
	// Read/Write
	DefaultApp            string           `url:"defaultApp"`
	Email                 string           `url:"email"`
	Password              string           `url:"password,omitempty"`
	RealName              string           `url:"realname"`
	RestartBackgroundJobs bool             `url:"restart_background_jobs"`
	Roles                 attributes.Roles `url:"roles"`
	TZ                    string           `url:"tz"`

	// Read-only fields are populated by results returned by the Splunk API, but
	// are not settable by Create or Update operations.
	Capabilities attributes.Capabilities `url:"-"`
	Type         string                  `url:"-"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Content
}

// User defines a Splunk user.
type User struct {
	client.Title `json:"name" url:"name"`
	UserContent  `json:"content"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.GlobalNamespace
	client.Endpoint `endpoint:"authentication/users"`
}
