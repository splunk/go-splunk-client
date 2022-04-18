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
	"reflect"

	"github.com/splunk/go-splunk-client/pkg/deepset"
	"github.com/splunk/go-splunk-client/pkg/service"
)

// Create performs a Create action for the given Entry.
func (client *Client) Create(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodPost),
			BuildRequestServiceURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesSelective(entry, "create"),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Created, HandleResponseJSONMessagesError()),
		),
	)
}

// Read performs a Read action for the given Entry. It modifies entry in-place,
// so entry must be a pointer.
func (client *Client) Read(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseCode(codes.NotFound, HandleResponseJSONMessagesCustomError(ErrorNotFound)),
			HandleResponseRequireCode(codes.Read, HandleResponseJSONMessagesError()),
			HandleResponseEntry(entry),
		),
	)
}

// Update performs an Update action for the given Entry.
func (client *Client) Update(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodPost),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesSelective(entry, "update"),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Updated, HandleResponseJSONMessagesError()),
		),
	)
}

// Delete performs a Delete action for the given Entry.
func (client *Client) Delete(entry interface{}) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodDelete),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Deleted, HandleResponseJSONMessagesError()),
		),
	)
}

func (client *Client) listModified(entries interface{}, modifier interface{}) error {
	entriesPtrV := reflect.ValueOf(entries)
	if entriesPtrV.Kind() != reflect.Ptr {
		return wrapError(ErrorPtr, nil, "client: List attempted on on-pointer value")
	}

	entriesV := reflect.Indirect(entriesPtrV)
	if entriesV.Kind() != reflect.Slice {
		return wrapError(ErrorSlice, nil, "client: List attempted on non-slice value")
	}
	entryT := entriesV.Type().Elem()
	entryI := reflect.New(entryT).Interface()

	if modifier != nil {
		if err := deepset.Set(entryI, modifier); err != nil {
			return err
		}
	}

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURL(client, entryI),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntries(entries),
		),
	)
}

// ListNamespace populates entries in place for a Namespace.
func (client *Client) ListNamespace(entries interface{}, ns Namespace) error {
	return client.listModified(entries, ns)
}

// ListNamespace populates entries in place for an ID.
func (client *Client) ListID(entries interface{}, id ID) error {
	return client.listModified(entries, id)
}

// ListNamespace populates entries in place without any ID or Namespace context.
func (client *Client) List(entries interface{}) error {
	return client.listModified(entries, nil)
}
