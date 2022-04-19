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

package entry

import (
	"net/url"

	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/client"
)

// SavedSearchContent defines the content of a Savedsearch object.
type SavedSearchContent struct {
	Search  attributes.Explicit[string]          `json:"search" values:"search,omitzero"`
	Actions attributes.NamedParametersCollection `json:"-" named_parameters_collection:"action" values:"action,omitzero"`
}

// AddURLValues implements custom additional encoding to url.Values.
func (content SavedSearchContent) AddURLValues(key string, v *url.Values) error {
	// The Splunk REST API returns savedsearch action status like "action.email=1", but doesn't honor that format
	// for setting the action statuses. To set an action status, you must pass "actions=action1,action2" formatted
	// values. Here we iterate through the enabled actions and add a url.Values entry for all enabled actions.
	//
	// If Actions is empty (not nil), we "clear" the enabled actions list by setting a single empty value for "actions".

	if content.Actions != nil && len(content.Actions) == 0 {
		v.Add("actions", "")
	}

	for _, enabledActionName := range content.Actions.EnabledNames() {
		v.Add("actions", enabledActionName)
	}

	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling.
func (content *SavedSearchContent) UnmarshalJSON(data []byte) error {
	type contentAlias SavedSearchContent
	var newAliasedContent contentAlias

	if err := attributes.UnmarshalJSONForNamedParametersCollections(data, &newAliasedContent); err != nil {
		return err
	}

	*content = SavedSearchContent(newAliasedContent)

	return nil
}

// SavedSearch defines a Splunk savedsearch.
type SavedSearch struct {
	ID      client.ID          `service:"saved/searches" selective:"create"`
	Content SavedSearchContent `json:"content" values:",anonymize"`
}
