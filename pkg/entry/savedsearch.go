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
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/splunk/go-sdk/pkg/attributes"
	"github.com/splunk/go-sdk/pkg/client"
)

// SavedSearchAction represents a single configurable action of a SavedSearch.
type SavedSearchAction struct {
	// action.<name>
	Name string
	// if Disabled=false, Name will be included in the "actions" field
	Disabled bool
	// action.<name>.<key> = <value>
	Parameters map[string]string
}

// parameterValues returns a map of formatted action.<name>.<key>:<value> values.
func (action SavedSearchAction) parameterValues() map[string]string {
	paramMap := make(map[string]string)

	for key, value := range action.Parameters {
		paramName := fmt.Sprintf("action.%s.%s", action.Name, key)
		paramMap[paramName] = value
	}

	return paramMap
}

// SavedSearchActions is a collection of SavedSearchAction objects.
type SavedSearchActions []SavedSearchAction

// UnmarshalJSON implements custom JSON unmarshaling.
func (actions *SavedSearchActions) UnmarshalJSON(data []byte) error {
	newActionsMap := map[string]SavedSearchAction{}

	// first unmarshal into a generic key/interface map so we can examine each key/value separately
	valuesMap := map[string]interface{}{}
	if err := json.Unmarshal(data, &valuesMap); err != nil {
		return err
	}

	// create a new SavedSearchAction for each enabled action
	actionsI, ok := valuesMap["actions"]
	if ok {
		actionsString := actionsI.(string)

		// "actions" is a comma-separate list of action names that are enabled
		for _, actionName := range strings.Split(actionsString, ",") {
			// if the list was empty, it's possible we have a zero-length value, which should just be skipped
			if actionName != "" {
				newActionsMap[actionName] = SavedSearchAction{Name: actionName}
			}
		}
	}

	// iterate through all keys and values in the unmarshaled map
	for key, valueI := range valuesMap {
		keySegments := strings.Split(key, ".")
		if keySegments[0] != "action" {
			continue
		}

		// silently skipping action configurations that aren't action.<name>.<param>
		if len(keySegments) < 3 {
			continue
		}

		// action.<name>.<param> = <value>
		actionName := keySegments[1]
		// <param> can contain dots, so re-join them for actionParam
		actionParam := strings.Join(keySegments[2:], ".")

		// every found <name> should have a corresponding SavedSearchAction, so create one if it doesn't already exist
		_, ok = newActionsMap[actionName]
		if !ok {
			newActionsMap[actionName] = SavedSearchAction{
				Name: actionName,
				// enabled actions created by handling the "actions" key, so any actions created here are disabled
				Disabled: true,
			}
		}

		// init Parameters if necessary
		if newActionsMap[actionName].Parameters == nil {
			newAction := newActionsMap[actionName]
			newAction.Parameters = map[string]string{}
			newActionsMap[actionName] = newAction
		}

		// set the Parameter value to the string representation of the found value (only string/bool/int/float types are expected)
		switch reflect.TypeOf(valueI).Kind() {
		default:
			return fmt.Errorf("unknown value type %T for action.%s.%s", valueI, actionName, actionParam)
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			newActionsMap[actionName].Parameters[actionParam] = fmt.Sprintf("%v", valueI)
		}
	}

	// get a sorted list of keys in newActionsMap, because map iteration is random and we want to be
	// deterministic at least for testing purposes
	newActionsMapKeys := make([]string, 0, len(newActionsMap))
	for key := range newActionsMap {
		newActionsMapKeys = append(newActionsMapKeys, key)
	}
	sort.Strings(newActionsMapKeys)

	// create new SavedSearchActions from sorted list of action names
	newActions := make(SavedSearchActions, 0, len(newActionsMap))
	for _, newActionMapKey := range newActionsMapKeys {
		newActions = append(newActions, newActionsMap[newActionMapKey])
	}

	*actions = newActions

	return nil
}

// EncodeValues implements custom encoding to url.Values.
func (actions SavedSearchActions) EncodeValues(key string, v *url.Values) error {
	if actions == nil {
		return nil
	}

	actionNames := make([]string, 0, len(actions))

	for _, action := range actions {
		if !action.Disabled {
			actionNames = append(actionNames, action.Name)
		}

		for key, value := range action.parameterValues() {
			v.Set(key, value)
		}
	}

	v.Set("actions", strings.Join(actionNames, ","))

	return nil
}

// SavedSearchContent defines the content of a Savedsearch object.
type SavedSearchContent struct {
	// Read/Write
	Actions SavedSearchActions `json:"-"` // Actions not unmarshalled as part of SavedSearchContent
	
	Search attributes.String `json:"search" url:"search"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Content
}

// SavedSearch defines a Splunk savedsearch.
type SavedSearch struct {
	client.ID
	SavedSearchContent `json:"content"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Endpoint `endpoint:"saved/searches" json:"endpoint"`
}

// UnmarshalJSON implements custom JSON umarshaling for SavedSearch.
func (search *SavedSearch) UnmarshalJSON(data []byte) error {
	// SavedSearch has nested object SavedSearchActions, but the API returns the configurations
	// for them as dotted attribute names directly under "content". The returned JSON has content
	// like this:
	//
	//	"action.email":true,
	//	"action.email.allow_empty_attachment":"1",
	//	"action.email.allowedDomainList":"",
	//	...
	//	"search":"index=_internal",
	//
	// All of the "action.<name>" fields are handled by custom unmarshaling of SavedSearchActions,
	// while the rest is handled by standard unmarshaling.

	// first unmarshal a type idential to SavedSearch, but with standard unmarshaling.
	type ss SavedSearch
	newSS := ss{}
	if err := json.Unmarshal(data, &newSS); err != nil {
		return fmt.Errorf("newSS: %s", err)
	}

	// then unmarshal a custom type that only handles SavedSearchActions
	newActions := struct {
		Actions SavedSearchActions `json:"content"`
	}{}
	if err := json.Unmarshal(data, &newActions); err != nil {
		return fmt.Errorf("newActions: %s", err)
	}

	// add the custom unmarshaled SavedSearchActions to the vanilla unmarshaled value
	newSS.Actions = newActions.Actions

	*search = SavedSearch(newSS)

	return nil
}
