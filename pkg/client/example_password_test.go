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

package client_test

import (
	"github.com/splunk/go-splunk-client/pkg/authenticators"
	"github.com/splunk/go-splunk-client/pkg/client"
)

func ExampleClient_passwordAuth() {
	_ = client.Client{
		URL: "https://splunk.example.com:8089",
		Authenticator: &authenticators.Password{
			Username: "admin",
			Password: "changeme",
		},
	}
}
