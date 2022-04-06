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

// Entry is the interface that describes types that are support Create, Read, Update,
// Delete, List operations. Types that satisfy this interface meet the Service, Titler,
// and ContentGetter interfaces.
type Entry interface {
	ContentGetter
}

// Create performs a Create action for the given Entry.
func Create(client *Client, entry Entry) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodPost),
			BuildRequestServiceURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValues(entry),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Created, HandleResponseJSONMessagesError()),
		),
	)
}

// Read performs a Read action for the given Entry. It modifies entry in-place,
// so entry must be a pointer.
func Read(client *Client, entry Entry) error {
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
func Update(client *Client, entry Entry) error {
	var codes service.StatusCodes

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestGetServiceStatusCodes(entry, &codes),
			BuildRequestMethod(http.MethodPost),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesContent(entry),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(codes.Updated, HandleResponseJSONMessagesError()),
		),
	)
}

// Delete performs a Delete action for the given Entry.
func Delete(client *Client, entry Entry) error {
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

func listModified(client *Client, entries interface{}, modifier interface{}) error {
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
func ListNamespace(client *Client, entries interface{}, ns Namespace) error {
	return listModified(client, entries, ns)
}

// ListNamespace populates entries in place for an ID.
func ListID(client *Client, entries interface{}, id ID) error {
	return listModified(client, entries, id)
}

// ListNamespace populates entries in place without any ID or Namespace context.
func List(client *Client, entries interface{}) error {
	return listModified(client, entries, nil)
}
