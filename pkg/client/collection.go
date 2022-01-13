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

import "net/http"

type collection interface {
	servicePath(interface{}) (string, error)
	endpointPath(...string) (string, error)
}

// collectionServicePath returns the service path for a collection. It is abstracted into
// its own function only because it is odd to pass the collection back to one of its own
// methods, which is the case due to servicePath being an inherited method by way of a collection
// having service as an anonymous field.
func collectionServicePath(coll collection) (string, error) {
	return coll.servicePath(coll)
}

func collectionReadRequest(client *Client, coll collection) (*http.Request, error) {
	servicePath, err := collectionServicePath(coll)
	if err != nil {
		return nil, err
	}

	url, err := client.urlForPath(servicePath)
	if err != nil {
		return nil, err
	}

	r := &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}

	if err := client.authenticateRequest(r); err != nil {
		return nil, err
	}

	return r, nil
}
