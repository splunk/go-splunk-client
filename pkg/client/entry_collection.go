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

// entryCollection types represent collections of entries.
type entryCollection struct {
	path string
}

// func entryCollectionReadRequest(client *Client, collection entryCollection) (*http.Request, error) {
// 	ns, err := namespace(collection)
// 	if err != nil {
// 		return nil, err
// 	}

// 	url, err := client.urlForPath(ns, collection.collectionPath())
// 	if err != nil {
// 		return nil, err
// 	}

// 	r := &http.Request{
// 		Method: http.MethodGet,
// 		URL:    url,
// 	}

// 	if err := client.Authenticator.authenticateRequest(client, r); err != nil {
// 		return nil, err
// 	}

// 	return r, nil
// }
