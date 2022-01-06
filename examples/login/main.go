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
	"log"

	"github.com/splunk/go-sdk/pkg/client"
)

func main() {
	client := client.Client{
		URL:                   "https://localhost:8089",
		Username:              "admin",
		Password:              "changedit",
		TLSInsecureSkipVerify: true,
	}

	if err := client.Login(); err != nil {
		log.Fatalf("unable to login: %s", err)
	}
}