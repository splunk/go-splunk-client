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

// ID is the URL ID of a Splunk object.
type ID string

// GetNamespaceUserApp returns the parsed User and App from the ID. It returns an error
// if unable to parse the ID int
func (id ID) GetIDParts() (user string, app string, endpoint string, title string, err error) {
	reverseStrings := func(input []string) []string {
		output := make([]string, 0, len(input))

		for i := len(input) - 1; i >= 0; i-- {
			output = append(output, input[i])
		}

		return output
	}

	// we work backwards in the id path segments to find the namespace, so reverse idParts
	idParts := strings.Split(string(id), "/")
	idParts = reverseStrings(idParts)

	// store title
	foundTitle := idParts[0]
	// remove title
	idParts = idParts[1:]

	// 5 as caplen as a guess of how many segments an endpoint may have
	foundEndpointParts := make([]string, 0, 5)
	for i, idPart := range idParts {
		// if we've made it to the "services" segment we have an empty Namespace
		if idPart == "services" {
			user = ""
			app = ""
			endpoint = strings.Join(reverseStrings(foundEndpointParts), "/")
			title = foundTitle

			return
		}

		if idPart == "servicesNS" {
			// if we've made it to the "servicesNS" segment we have a user/app namespace

			if i < 2 {
				err = wrapError(ErrorNamespace, nil, "unable to parse id into Namespace, servicesNS found without user/app: %s", id)
				return
			}

			user = idParts[i-1]
			app = idParts[i-2]
			// get rid of the user/app pieces we previously found
			foundEndpointParts = foundEndpointParts[:len(foundEndpointParts)-2]
			endpoint = strings.Join(reverseStrings(foundEndpointParts), "/")
			title = foundTitle

			return
		}

		foundEndpointParts = append(foundEndpointParts, idPart)
	}

	err = wrapError(ErrorNamespace, nil, "unable to parse id into Namespace, missing services or servicesNS: %s", id)
	return
}
