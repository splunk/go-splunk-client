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
	"encoding/json"
	"net/url"

	"github.com/splunk/go-splunk-client/pkg/internal/paths"
)

// ID represents a Splunk object ID URL for a specific object.
type ID struct {
	Namespace Namespace

	// Title is the ID's title component. It is the name of the Splunk object.
	Title string
}

// parseID returns a new ID by parsing the ID URL string.
func parseID(idURL string) (ID, error) {
	newNS, remnants, err := parseNamespace(idURL)
	if err != nil {
		return ID{}, err
	}

	if len(remnants) < 1 {
		return ID{}, wrapError(ErrorID, nil, "client: parseNamespace didn't return a remnant for ID.Title")
	}

	return ID{
		Namespace: newNS,
		Title:     remnants[len(remnants)-1],
	}, nil
}

// Parse sets the ID's value to match what is parsed from the given ID URL.
func (id *ID) Parse(idURL string) error {
	newID, err := parseID(idURL)
	if err != nil {
		return err
	}

	*id = newID

	return nil
}

// GetServicePath implements custom GetServicePath encoding. It returns its Namespace's
// service path.
func (id ID) GetServicePath(path string) (string, error) {
	return id.Namespace.GetServicePath(path)
}

// GetEntryPath implements custom GetEntryPath encoding. It returns the url-encoded
// value of the ID's Title with the service path preceding it.
func (id ID) GetEntryPath(path string) (string, error) {
	servicePath, err := id.GetServicePath(path)
	if err != nil {
		return "", err
	}

	return paths.Join(servicePath, url.PathEscape(id.Title)), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for IDFields.
func (id *ID) UnmarshalJSON(data []byte) error {
	idString := ""
	if err := json.Unmarshal(data, &idString); err != nil {
		return wrapError(ErrorID, err, "client: unable to unmarshal %q as string", data)
	}

	if err := id.Parse(idString); err != nil {
		return err
	}

	return nil
}

// EncodeValues implements custom url.Query encoding of ID. It adds a field "name" for the ID's
// Title. If the Title value is empty, it returns an error, as there are no scenarios where an ID
// object is expected to be POSTed with an empty Title.
func (id ID) EncodeValues(key string, v *url.Values) error {
	if id.Title == "" {
		return wrapError(ErrorID, nil, "client: attempted EncodeValues on ID with empty Title")
	}

	v.Add("name", id.Title)

	return nil
}
