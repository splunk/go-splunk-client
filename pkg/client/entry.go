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

	"github.com/splunk/go-sdk/pkg/internal/paths"
)

// Entry is the interface that describes types that are support Create, Read, Update,
// Delete, List operations. Types that satisfy this interface meet the Service, Titler,
// and ContentGetter interfaces.
type Entry interface {
	Service
	Titler
	ContentGetter
}

// EntryLister is the interface that describes types that satisfy both the Entry and IDOptApplyer
// interfaces. This interface is a superset of Entry as IDOptApplyer requires a pointer, but Entry
// does not.
type EntryLister interface {
	Entry
	IDOptApplyer
}

// entryPath returns the path for the given CRUDLer. If the Entry has an
// empty Title, a valid path will be returned with the Title component being empty,
// because an Entry's path doesn't require a non-empty Title to be valid.
func entryPath(entry Entry) (string, error) {
	servicePath, err := servicePath(entry)
	if err != nil {
		return "", err
	}

	return paths.Join(servicePath, entry.Title()), nil
}

// Create performs a Create action for the given Entry.
func Create(client *Client, entry Entry) error {
	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodPost),
			BuildRequestServiceURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesWithTitle(entry),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusCreated, HandleResponseJSONMessagesError()),
		),
	)
}

// Read performs a Read action for the given Entry. It modifies entry in-place,
// so entry must be a pointer.
func Read(client *Client, entry Entry) error {
	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURLWithTitle(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseEntryNotFound(entry, HandleResponseJSONMessagesCustomError(ErrorNotFound)),
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntry(entry),
		),
	)
}

// Update performs an Update action for the given Entry.
func Update(client *Client, entry Entry) error {
	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodPost),
			BuildRequestEntryURLWithTitle(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesContent(entry),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
		),
	)
}

// Delete performs a Delete action for the given Entry.
func Delete(client *Client, entry Entry) error {
	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodDelete),
			BuildRequestEntryURLWithTitle(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
		),
	)
}

// List returns a list of the given type of Entry by performing a List
// action for an ID with the given ID Field Options applied.
func List(client *Client, entries interface{}, idFieldOpts ...IDOpt) error {
	entriesPtrV := reflect.ValueOf(entries)
	if entriesPtrV.Kind() != reflect.Ptr {
		return wrapError(ErrorPtr, nil, "List attempted on on-pointer value")
	}

	entriesV := reflect.Indirect(entriesPtrV)
	if entriesV.Kind() != reflect.Slice {
		return wrapError(ErrorSlice, nil, "List attempted on non-slice value")
	}
	entryT := entriesV.Type().Elem()
	entryPtrV := reflect.New(entryT)
	entryPtrI := entryPtrV.Interface()
	entryPtrEntry, ok := entryPtrI.(EntryLister)
	if !ok {
		entryI := reflect.Indirect(entryPtrV).Interface()
		return wrapError(ErrorSlice, nil, "List attempted on slice of non-EntryLister type %T", entryI)
	}

	entryPtrEntry.ApplyIDOpt(idFieldOpts...)

	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURL(client, entryPtrEntry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntries(entries),
		),
	)
}

func UpdateACL(client *Client, entry EntryAccessController) error {
	return client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodPost),
			BuildRequestEntryACLURL(client, entry),
			BuildRequestAccessControllerBodyValues(entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
		),
	)
}
