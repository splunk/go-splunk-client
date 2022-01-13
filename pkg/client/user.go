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

// UserAttributes are the Attributes of a User.
type UserAttributes struct {
	DefaultApp string `url:"defaultApp"`
}

// User represents a Splunk user.
type User struct {
	service `service:"authentication/users"`
	globalNamespace
	Title
	Attributes UserAttributes `json:"content"`
}
