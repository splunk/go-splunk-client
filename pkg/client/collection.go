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

// Collection is the interface that describes types that are collections in the REST API.
// Types that satisfy this interface meet the Service, Titler, and ContentGetter interfaces.
type Collection interface {
	Service
	Titler
	ContentGetter
}

// collectionPath returns the path for the given Collection. If the Collection has an
// empty Title, a valid path will be returned with the Title component being empty,
// because a collection's path doesn't require a non-empty Title to be valid.
func collectionPath(collection Collection) (string, error) {
	servicePath, err := servicePath(collection)
	if err != nil {
		return "", err
	}

	return paths.Join(servicePath, collection.TitleValue()), nil
}

// CollectionClient returns a list of the given type of Collection by performing a List
// action for its Service URL.
func CollectionList[C Collection](c *Client, collection C) ([]C, error) {
	entries := make([]C, 0)

	if err := c.performOperation(
		composeRequestBuilder(
			buildRequestMethod(http.MethodGet),
			buildRequestCollectionURL(c, collection),
			buildRequestOutputModeJSON(),
			buildRequestAuthenticate(c),
		),
		composeResponseHandler(
			handleResponseRequireCode(http.StatusOK),
			handleResponseEntries(&entries),
		),
	); err != nil {
		return nil, err
	}

	return entries, nil
}

// CollectionCreate performs a Create action for the given Collection item. It returns
// the Collection item that was created.
func CollectionCreate[C Collection](c *Client, collection C) (C, error) {
	entry := new(C)

	if err := c.performOperation(
		composeRequestBuilder(
			buildRequestMethod(http.MethodPost),
			buildRequestServiceURL(c, collection),
			buildRequestOutputModeJSON(),
			buildRequestBodyValuesWithTitle(collection),
			buildRequestAuthenticate(c),
		),
		composeResponseHandler(
			handleResponseRequireCodeWithMessage(http.StatusCreated),
			handleResponseEntry(entry),
		),
	); err !=  nil {
		return *new(C), err
	}

	return *entry, nil
}

// CollectionRead performs a Read action for the given Collection item. It returns
// the Collection item that was read.
func CollectionRead[C Collection](c *Client, collection C) (C, error) {
	entry := new(C)

	if err := c.performOperation(
		composeRequestBuilder(
			buildRequestMethod(http.MethodGet),
			buildRequestCollectionURLWithTitle(c, collection),
			buildRequestOutputModeJSON(),
			buildRequestAuthenticate(c),
		),
		composeResponseHandler(
			handleResponseRequireCodeWithMessage(http.StatusOK),
			handleResponseEntry(entry),
		),
	); err != nil {
		return *new(C), nil
	}

	return *entry, nil
}

// CollectionUpdate performs an Update action for the given Collection item. It
// returns the Collection item that resulted from the update.
func CollectionUpdate[C Collection](c *Client, collection C) (C, error) {
	entry := new(C)

	if err := c.performOperation(
		composeRequestBuilder(
			buildRequestMethod(http.MethodPost),
			buildRequestCollectionURLWithTitle(c, collection),
			buildRequestOutputModeJSON(),
			buildRequestBodyContentValues(collection),
			buildRequestAuthenticate(c),	
		),
		composeResponseHandler(
			handleResponseRequireCode(http.StatusOK),
			handleResponseEntry(entry),	
		),
	); err !=  nil {
		return *new(C), err
	}

	return *entry, nil
}

// CollectionDelete performs a Delete action for the given Collection item. It
// returns a list of the remaining Collection items after the deletion.
func CollectionDelete[C Collection](c *Client, collection C) ([]C, error) {
	entries := make([]C, 0)

	if err := c.performOperation(
		composeRequestBuilder(
			buildRequestMethod(http.MethodDelete),
			buildRequestCollectionURLWithTitle(c, collection),
			buildRequestOutputModeJSON(),
			buildRequestAuthenticate(c),
		),
		composeResponseHandler(
			handleResponseRequireCode(http.StatusOK),
			handleResponseEntries(&entries),
		),
	); err != nil {
		return nil, err
	}

	return entries, nil
}
