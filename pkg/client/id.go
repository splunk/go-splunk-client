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

import "strings"

// ID represents a Splunk object ID URL for a specific object.
type ID struct {
	IDFields IDFields `json:"id"`
}

// ParseIDURL returns a new ID by parsing an ID URL.
func ParseIDURL(idURL string) (ID, error) {
	fields, err := ParseIDFields(idURL)
	if err != nil {
		return ID{}, err
	}

	newID := ID{IDFields: fields}

	if newID.String() != idURL {
		return ID{}, wrapError(ErrorID, nil, "parsed ID string value (%s) doesn't match given URL (%s). This should be reported as a bug.", newID.String(), idURL)
	}

	return ID{IDFields: fields}, nil
}

// SetIDFromURL sets the ID's value to match what is parsed from the given ID URL.
func (id *ID) SetIDFromURL(idURL string) error {
	newID, err := ParseIDURL(idURL)
	if err != nil {
		return err
	}

	*id = newID

	return nil
}

// NamespacePath returns the namespace path for the given ID, or an error if the path could not be
// determined by parsing the ID. Even if an error is returned, the invalid calculated path will be returned.
func (id ID) NamespacePath() (string, error) {
	var path string

	// absence of both user/app indicates global context
	if (id.IDFields.User == "") && (id.IDFields.App == "") {
		path = "services"
	} else {
		path = strings.Join([]string{"servicesNS", id.IDFields.User, id.IDFields.App}, "/")
	}

	return path, id.IDFields.validate()
}

// Title returns the title value of an ID.
func (id ID) Title() string {
	return id.IDFields.Title
}

// String returns a string representation of the ID. If ID.Title is empty the resulting string
// will have a trailing slash. The string representation of an ID should not be assumed to be
// valid, as the NamespacePath component is not error checked here.
func (id ID) String() string {
	nsPath, _ := id.NamespacePath()

	return strings.Join(
		[]string{
			id.IDFields.baseURL,
			nsPath,
			id.IDFields.endpoint,
			id.IDFields.Title,
		},
		"/",
	)
}

// IDOpt is a function that performs an operation on an ID object.
type IDOpt func(*ID)

// IDOptApplyer is the interface that describes types that implement ApplyIDOpt.
type IDOptApplyer interface {
	ApplyIDOpt(...IDOpt)
}

// ApplyIDOpt applies the given IDOpt functions.
func (id *ID) ApplyIDOpt(opts ...IDOpt) {
	for _, opt := range opts {
		opt(id)
	}
}

// SetTitle returns an IDOpt that sets Title on an ID object.
func SetTitle(title string) IDOpt {
	return func(id *ID) {
		id.IDFields.Title = title
	}
}

// SetNamespace returns an IDOpt that sets User and App on an ID object.
func SetNamespace(user string, app string) IDOpt {
	return func(id *ID) {
		id.IDFields.User = user
		id.IDFields.App = app
	}
}
