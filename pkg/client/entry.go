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
	"net/http"
)
// Entry is the interface that describes types that are support Create, Read, Update,
// Delete, List operations. Types that satisfy this interface meet the Service, Titler,
// and ContentGetter interfaces.
type Entry interface {
	Service
	Titler
	ContentGetter
}

// Create performs a Create action for the given Entry. It returns
// the Entry that was created.
func Create[E Entry](client *Client, entry E) (E, error) {
	createdEntry := new(E)

	if err := client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodPost),
			BuildRequestServiceURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesWithTitle(entry),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusCreated, HandleResponseJSONMessagesError()),
			HandleResponseEntry(createdEntry),
		),
	); err !=  nil {
		return *new(E), err
	}

	return *createdEntry, nil
}
