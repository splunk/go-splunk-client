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
	"path"
	"strings"
)

// Namespace represents a Splunk namespace/context. If both User/App are empty, the
// global context will be used. Both User and App must be unset or set. It is invalid
// for one to be set without the other.
type Namespace struct {
	User string
	App  string
}

// validate returns an error if Namespace is invalid. It is invalid if either User or App
// is set without the other.
func (c Namespace) validate() error {
	// neither or both App and User must be set
	// this is true if each of their (== "") tests are the same
	if (c.User == "") == (c.App == "") {
		return nil
	}

	return wrapError(ErrorNamespace, nil, "invalid Context, neither or both App and User must be set")
}

// NamespacePath returns the namespace path.
func (c Namespace) NamespacePath() (string, error) {
	if err := c.validate(); err != nil {
		return "", err
	}

	// absence of either field indicates global context
	if c.App == "" {
		return "services", nil
	}

	return path.Join("servicesNS", c.User, c.App), nil
}

// namespaceForID returns a new Namespace by parsing an ID URL. It returns an error if unable to
// properly parse the ID string.
func namespaceForID(id string) (Namespace, error) {
	reverseStrings := func(input []string) []string {
		output := make([]string, 0, len(input))

		for i := len(input) - 1; i >= 0; i-- {
			output = append(output, input[i])
		}

		return output
	}

	// we work backwards in the id path segments to find the namespace, so reverse idParts
	idParts := strings.Split(id, "/")
	idParts = reverseStrings(idParts)

	// remove title from idParts (in case it happens to be named services or servicesNS)
	idParts = idParts[1:]

	for i, idPart := range idParts {
		// if we've made it to the "services" segment we have an empty Namespace
		if idPart == "services" {
			return Namespace{}, nil
		}

		if idPart == "servicesNS" {
			// if we've made it to the "servicesNS" segment we have a user/app namespace

			if i < 2 {
				return Namespace{}, wrapError(ErrorNamespace, nil, "unable to parse id into Namespace, servicesNS found without user/app: %s", id)
			}

			return Namespace{
				User: idParts[i-1],
				App:  idParts[i-2],
			}, nil
		}
	}

	return Namespace{}, wrapError(ErrorNamespace, nil, "unable to parse id into Namespace, missing services or servicesNS: %s", id)
}

// UnmarshalJSON implements custom JSON unmarshaling of a plain string ID
// value into a Namespace.
func (c *Namespace) UnmarshalJSON(data []byte) error {
	value := new(string)

	if err := json.Unmarshal(data, value); err != nil {
		return wrapError(ErrorNamespace, err, err.Error())
	}

	ns, err := namespaceForID(*value)
	if err != nil {
		return err
	}

	*c = ns

	return nil
}
