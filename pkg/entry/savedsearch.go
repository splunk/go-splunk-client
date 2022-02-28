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
	AlertDigestMode                   attributes.Bool   `json:"alert.digest_mode"                     url:"alert.digest_mode"`
	AlertExpires                      attributes.Int    `json:"alert.expires"                         url:"alert.expires"`
	AlertSeverity                     attributes.String `json:"alert.severity"                        url:"alert.severity"`
	AlertSuppress                     attributes.Bool   `json:"alert.suppress"                        url:"alert.suppress"`
	AlertSuppressFields               attributes.String `json:"alert.suppress.fields"                 url:"alert.suppress.fields"`
	AlertSuppressGroupName            attributes.String `json:"alert.suppress.group_name"             url:"alert.suppress.group_name"`
	AlertSuppressPeriod               attributes.Int    `json:"alert.suppress.period"                 url:"alert.suppress.period"`
	AlertTrack                        attributes.String `json:"alert.track"                           url:"alert.track"`
	AlertComparator                   attributes.String `json:"alert_comparator"                      url:"alert_comparator"`
	AlertCondition                    attributes.String `json:"alert_condition"                       url:"alert_condition"`
	AlertThreshold                    attributes.Int    `json:"alert_threshold"                       url:"alert_threshold"`
	AlertType                         attributes.String `json:"alert_type"                            url:"alert_type"`
	AllowSkew                         attributes.String `json:"allow_skew"                            url:"allow_skew"`
	AutoSummarize                     attributes.Bool   `json:"auto_summarize"                        url:"auto_summarize"`
	AutoSummarizeCommand              attributes.String `json:"auto_summarize.command"                url:"auto_summarize.command"`
	AutoSummarizeCronSchedule         attributes.String `json:"auto_summarize.cron_schedule"          url:"auto_summarize.cron_schedule"`
	AutoSummarizeDispatchEarliestTime attributes.String `json:"auto_summarize.dispatch.earliest_time" url:"auto_summarize.dispatch.earliest_time"`
	AutoSummarizeDispatchLatestTime   attributes.String `json:"auto_summarize.dispatch.latest_time"   url:"auto_summarize.dispatch.latest_time"`
	AutoSummarizeDispatchTimeFormat   attributes.String `json:"auto_summarize.dispatch.time_format"   url:"auto_summarize.dispatch.time_format"`
	AutoSummarizeDispatchTtl          attributes.String `json:"auto_summarize.dispatch.ttl"           url:"auto_summarize.dispatch.ttl"`
	AutoSummarizeMaxConcurrent        attributes.Int    `json:"auto_summarize.max_concurrent"         url:"auto_summarize.max_concurrent"`
	AutoSummarizeMaxDisabledBuckets   attributes.Int    `json:"auto_summarize.max_disabled_buckets"   url:"auto_summarize.max_disabled_buckets"`
	AutoSummarizeMaxSummaryRatio      attributes.Int    `json:"auto_summarize.max_summary_ratio"      url:"auto_summarize.max_summary_ratio"`
	AutoSummarizeMaxSummarySize       attributes.Int    `json:"auto_summarize.max_summary_size"       url:"auto_summarize.max_summary_size"`
	AutoSummarizeMaxTime              attributes.Int    `json:"auto_summarize.max_time"               url:"auto_summarize.max_time"`
	AutoSummarizeSuspendPeriod        attributes.String `json:"auto_summarize.suspend_period"         url:"auto_summarize.suspend_period"`
	AutoSummarizeTimespan             attributes.String `json:"auto_summarize.timespan"               url:"auto_summarize.timespan"`
	CronSchedule                      attributes.String `json:"cron_schedule"                         url:"cron_schedule"`
	Description                       attributes.String `json:"description"                           url:"description"`
	Disabled                          attributes.Bool   `json:"disabled"                              url:"disabled"`
	DispatchAs                        attributes.String `json:"dispatchAs"                            url:"dispatchAs"`
	Displayview                       attributes.String `json:"displayview"                           url:"displayview"`
	DurableBackfillType               attributes.String `json:"durable.backfill_type"                 url:"durable.backfill_type"`
	DurableLagTime                    attributes.Int    `json:"durable.lag_time"                      url:"durable.lag_time"`
	DurableMaxBackfillIntervals       attributes.Int    `json:"durable.max_backfill_intervals"        url:"durable.max_backfill_intervals"`
	DurableTrackTimeType              attributes.String `json:"durable.track_time_type"               url:"durable.track_time_type"`
	IsScheduled                       attributes.Bool   `json:"is_scheduled"                          url:"is_scheduled"`
	IsVisible                         attributes.Bool   `json:"is_visible"                            url:"is_visible"`
	MaxConcurrent                     attributes.Int    `json:"max_concurrent"                        url:"max_concurrent"`
	Name                              attributes.String `json:"name"                                  url:"name"`
	NextScheduledTime                 attributes.String `json:"next_scheduled_time"                   url:"next_scheduled_time"`
	QualifiedSearch                   attributes.String `json:"qualifiedSearch"                       url:"qualifiedSearch"`
	RealtimeSchedule                  attributes.Bool   `json:"realtime_schedule"                     url:"realtime_schedule"`
	RequestUiDispatchApp              attributes.String `json:"request.ui_dispatch_app"               url:"request.ui_dispatch_app"`
	RequestUiDispatchView             attributes.String `json:"request.ui_dispatch_view"              url:"request.ui_dispatch_view"`
	RestartOnSearchpeerAdd            attributes.Bool   `json:"restart_on_searchpeer_add"             url:"restart_on_searchpeer_add"`
	RunNTimes                         attributes.Int    `json:"run_n_times"                           url:"run_n_times"`
	RunOnStartup                      attributes.Bool   `json:"run_on_startup"                        url:"run_on_startup"`
	SchedulePriority                  attributes.String `json:"schedule_priority"                     url:"schedule_priority"`
	ScheduleWindow                    attributes.String `json:"schedule_window"                       url:"schedule_window"`
	Search                            attributes.String `json:"search"                                url:"search"`
	Vsid                              attributes.String `json:"vsid"                                  url:"vsid"`
	WorkloadPool                      attributes.String `json:"workload_pool"                         url:"workload_pool"`

	Actions SavedSearchActions `json:"-"` // Actions not unmarshalled as part of SavedSearchContent

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
