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
	"strings"
)

// IDFields represents the components of an ID.
type IDFields struct {
	// User and App define the ID's namespace.
	User string
	App  string

	// Title is the ID's title component. It is the name of the Splunk object.
	Title string

	// baseURL is the part of the ID leading up to services/ or servicesNS. Its value
	// is irrelevant to the operation of this library, and is present only to permit
	// complete representation of the contents of an ID.
	baseURL string

	// endpoint is the resource endpoint segment of the ID, such as authorization/users.
	// Its value is irrelevant to the operation of this library, and is present only to permit
	// complete representation of the contents of an ID.
	endpoint string
}

// ParseIDFields parses a string ID into IDFields.
func ParseIDFields(id string) (IDFields, error) {
	// start with a clean IDFields
	newFields := IDFields{}

	// an empty id is a "clean slate", with no namespace info (user/app), endpoint, or title
	if id == "" {
		return newFields, nil
	}

	reverseStrings := func(input []string) []string {
		output := make([]string, 0, len(input))

		for i := len(input) - 1; i >= 0; i-- {
			output = append(output, input[i])
		}

		return output
	}

	// we work backwards in the id path segments to find the namespace, so reverse idPartStrings
	idPartStrings := strings.Split(id, "/")
	idPartStrings = reverseStrings(idPartStrings)

	// store title
	newFields.Title = idPartStrings[0]
	// remove title from string parts
	idPartStrings = idPartStrings[1:]

	// 5 as caplen as a guess of how many segments an endpoint may have
	foundEndpointParts := make([]string, 0, 5)
	for i, idPartString := range idPartStrings {
		// if we've made it to the "services" segment we have an empty Namespace
		if idPartString == "services" {
			newFields.endpoint = strings.Join(foundEndpointParts, "/")
			// baseURL is the remaining segments, but in reverse order
			newFields.baseURL = strings.Join(reverseStrings(idPartStrings[i+1:]), "/")

			return newFields, nil
		}

		if idPartString == "servicesNS" {
			// if we've made it to the "servicesNS" segment we have a user/app namespace

			if i < 2 {
				return IDFields{}, wrapError(ErrorID, nil, "unable to parse ID, servicesNS found without user/app: %s", id)
			}

			newFields.User = idPartStrings[i-1]
			newFields.App = idPartStrings[i-2]
			// foundEndpointParts will have the user/app segments still, so remove them
			foundEndpointParts = foundEndpointParts[2:]
			newFields.endpoint = strings.Join(foundEndpointParts, "/")
			// baseURL is the remaining segments, but in reverse order
			newFields.baseURL = strings.Join(reverseStrings(idPartStrings[i+1:]), "/")

			return newFields, nil
		}

		// if still looking for services/servicesNS, we're still working with endpoint segments, so prepend with this new segment
		foundEndpointParts = append([]string{idPartString}, foundEndpointParts...)
	}

	return IDFields{}, wrapError(ErrorID, nil, "unable to parse ID, missing services or servicesNS: %s", id)
}

// UnmarshalJSON implements custom JSON unmarshaling for IDFields.
func (fields *IDFields) UnmarshalJSON(data []byte) error {
	idString := ""
	if err := json.Unmarshal(data, &idString); err != nil {
		return wrapError(ErrorID, err, "unable to unmarshal %q as string", data)
	}

	parsedFields, err := ParseIDFields(idString)
	if err != nil {
		return err
	}

	*fields = parsedFields

	return nil
}

// validate returns an error if an ID is unable to be parsed or the resulting namespace would be invalid.
func (fields IDFields) validate() error {
	if (fields.User == "") != (fields.App == "") {
		return wrapError(ErrorNamespace, nil, "invalid ID, user and app must both be empty or non-empty")
	}

	return nil
}

// EncodeValues implements custom url.Query encoding of IDFields. It adds a field "name" for the ID's
// Title. If the Title value is empty, it returns an error, as there are no scenarios where an ID
// object is expected to be POSTed with an empty Title.
func (id IDFields) EncodeValues(key string, v *url.Values) error {
	title := id.Title
	if title == "" {
		return wrapError(ErrorID, nil, "attempted encode ID with empty title")
	}

	v.Add("name", title)

	return nil
}
