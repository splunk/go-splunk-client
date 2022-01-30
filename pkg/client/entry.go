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

// entryPath returns the path for the given CRUDLer. If the Entry has an
// empty Title, a valid path will be returned with the Title component being empty,
// because an Entry's path doesn't require a non-empty Title to be valid.
func entryPath(entry Entry) (string, error) {
	servicePath, err := servicePath(entry)
	if err != nil {
		return "", err
	}

	return paths.Join(servicePath, entry.TitleValue()), nil
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

// Read performs a Read action for the given Entry. It returns
// the Entry that was read.
func Read[E Entry](client *Client, entry E) (E, error) {
	readEntry := new(E)

	if err := client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURLWithTitle(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseCode(http.StatusNotFound, HandleResponseJSONMessagesCustomError(ErrorNotFound)),
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntry(readEntry),
		),
	); err != nil {
		return *new(E), err
	}

	return *readEntry, nil
}

// Update performs an Update action for the given Entry. It
// returns the Entry that resulted from the update.
func Update[E Entry](client *Client, entry E) (E, error) {
	updatedEntry := new(E)

	if err := client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodPost),
			BuildRequestEntryURLWithTitle(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestBodyValuesContent(entry),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntry(updatedEntry),
		),
	); err !=  nil {
		return *new(E), err
	}

	return *updatedEntry, nil
}

// Delete performs a Delete action for the given Entry. It
// returns a list of the remaining Entry items after the deletion.
func Delete[E Entry](client *Client, entry E) ([]E, error) {
	remainingEntries := make([]E, 0)

	if err := client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodDelete),
			BuildRequestEntryURLWithTitle(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntries(&remainingEntries),
		),
	); err != nil {
		return nil, err
	}

	return remainingEntries, nil
}

// List returns a list of the given type of Entry by performing a List
// action for its Entry URL, which may or may not have a value for its
// Title.
func List[E Entry](client *Client, entry E) ([]E, error) {
	listedEntries := make([]E, 0)

	if err := client.RequestAndHandle(
		ComposeRequestBuilder(
			BuildRequestMethod(http.MethodGet),
			BuildRequestEntryURL(client, entry),
			BuildRequestOutputModeJSON(),
			BuildRequestAuthenticate(client),
		),
		ComposeResponseHandler(
			HandleResponseRequireCode(http.StatusOK, HandleResponseJSONMessagesError()),
			HandleResponseEntries(&listedEntries),
		),
	); err != nil {
		return nil, err
	}

	return listedEntries, nil
}
