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
	"github.com/splunk/go-sdk/pkg/attributes"
	"github.com/splunk/go-sdk/pkg/client"
)

// SavedSearchContent defines the content of a Savedsearch object.
type SavedSearchContent struct {
	// Read/Write
	// Actions                attributes.String `json:"actions" url:"actions"`
	// AlertDigestMode        attributes.Bool   `json:"alert.digest_mode" url:"alert.digest_mode"`
	// AlertExpires           attributes.Int    `json:"alert.expires" url:"alert.expires"`
	// AlertSeverity          attributes.String `json:"alert.severity" url:"alert.severity"`
	// AlertSuppress          attributes.Bool   `json:"alert.suppress" url:"alert.suppress"`
	// AlertSuppressFields    attributes.String `json:"alert.suppress.fields" url:"alert.suppress.fields"`
	// AlertSuppressGroupName attributes.String `json:"alert.suppress.group_name" url:"alert.suppress.group_name"`
	// AlertSuppressPeriod    attributes.Int    `json:"alert.suppress.period" url:"alert.suppress.period"`
	// AlertTrack             attributes.String `json:"alert.track" url:"alert.track"`
	// AlertComparator        attributes.String `json:"alert_comparator" url:"alert_comparator"`
	// AlertCondition         attributes.String `json:"alert_condition" url:"alert_condition"`
	// AlertThreshold         attributes.Int    `json:"alert_threshold" url:"alert_threshold"`
	// AlertType              attributes.String `json:"alert_type" url:"alert_type"`
	// AllowSkew              attributes.String `json:"allow_skew" url:"allow_skew"`
	// Args* attributes.String `json:"args.*" url:"args.*"`
	// AutoSummarize                     attributes.Bool   `json:"auto_summarize" url:"auto_summarize"`
	// AutoSummarizeCommand              attributes.String `json:"auto_summarize.command" url:"auto_summarize.command"`
	// AutoSummarizeCronSchedule         attributes.String `json:"auto_summarize.cron_schedule" url:"auto_summarize.cron_schedule"`
	// AutoSummarizeDispatchEarliestTime attributes.String `json:"auto_summarize.dispatch.earliest_time" url:"auto_summarize.dispatch.earliest_time"`
	// AutoSummarizeDispatchLatestTime   attributes.String `json:"auto_summarize.dispatch.latest_time" url:"auto_summarize.dispatch.latest_time"`
	// AutoSummarizeDispatchTimeFormat   attributes.String `json:"auto_summarize.dispatch.time_format" url:"auto_summarize.dispatch.time_format"`
	// AutoSummarizeDispatchTtl          attributes.String `json:"auto_summarize.dispatch.ttl" url:"auto_summarize.dispatch.ttl"`
	// AutoSummarizeMaxConcurrent        attributes.Int    `json:"auto_summarize.max_concurrent" url:"auto_summarize.max_concurrent"`
	// AutoSummarizeMaxDisabledBuckets   attributes.Int    `json:"auto_summarize.max_disabled_buckets" url:"auto_summarize.max_disabled_buckets"`
	// AutoSummarizeMaxSummaryRatio      attributes.Int    `json:"auto_summarize.max_summary_ratio" url:"auto_summarize.max_summary_ratio"`
	// AutoSummarizeMaxSummarySize       attributes.Int    `json:"auto_summarize.max_summary_size" url:"auto_summarize.max_summary_size"`
	// AutoSummarizeMaxTime              attributes.Int    `json:"auto_summarize.max_time" url:"auto_summarize.max_time"`
	// AutoSummarizeSuspendPeriod        attributes.String `json:"auto_summarize.suspend_period" url:"auto_summarize.suspend_period"`
	// AutoSummarizeTimespan             attributes.String `json:"auto_summarize.timespan" url:"auto_summarize.timespan"`
	// CronSchedule                      attributes.String `json:"cron_schedule" url:"cron_schedule"`
	// Description                       attributes.String `json:"description" url:"description"`
	// Disabled                          attributes.Bool   `json:"disabled" url:"disabled"`
	// Dispatch* attributes.String `json:"dispatch.*" url:"dispatch.*"`
	// DispatchAllowPartialResults    attributes.Bool   `json:"dispatch.allow_partial_results" url:"dispatch.allow_partial_results"`
	// DispatchAutoCancel             attributes.Int    `json:"dispatch.auto_cancel" url:"dispatch.auto_cancel"`
	// DispatchAutoPause              attributes.Int    `json:"dispatch.auto_pause" url:"dispatch.auto_pause"`
	// DispatchBuckets                attributes.Int    `json:"dispatch.buckets" url:"dispatch.buckets"`
	// DispatchEarliestTime           attributes.String `json:"dispatch.earliest_time" url:"dispatch.earliest_time"`
	// DispatchIndexEarliest          attributes.String `json:"dispatch.index_earliest" url:"dispatch.index_earliest"`
	// DispatchIndexLatest            attributes.String `json:"dispatch.index_latest" url:"dispatch.index_latest"`
	// DispatchIndexedRealtime        attributes.Bool   `json:"dispatch.indexedRealtime" url:"dispatch.indexedRealtime"`
	// DispatchIndexedRealtimeOffset  attributes.Int    `json:"dispatch.indexedRealtimeOffset" url:"dispatch.indexedRealtimeOffset"`
	// DispatchIndexedRealtimeMinSpan attributes.Int    `json:"dispatch.indexedRealtimeMinSpan" url:"dispatch.indexedRealtimeMinSpan"`
	// DispatchLatestTime             attributes.String `json:"dispatch.latest_time" url:"dispatch.latest_time"`
	// DispatchLookups                attributes.Bool   `json:"dispatch.lookups" url:"dispatch.lookups"`
	// DispatchMaxCount               attributes.Int    `json:"dispatch.max_count" url:"dispatch.max_count"`
	// DispatchMaxTime                attributes.Int    `json:"dispatch.max_time" url:"dispatch.max_time"`
	// DispatchReduceFreq             attributes.Int    `json:"dispatch.reduce_freq" url:"dispatch.reduce_freq"`
	// DispatchRtBackfill             attributes.Bool   `json:"dispatch.rt_backfill" url:"dispatch.rt_backfill"`
	// DispatchRtMaximumSpan          attributes.Int    `json:"dispatch.rt_maximum_span" url:"dispatch.rt_maximum_span"`
	// DispatchSampleRatio            attributes.Int    `json:"dispatch.sample_ratio" url:"dispatch.sample_ratio"`
	// DispatchSpawnProcess           attributes.Bool   `json:"dispatch.spawn_process" url:"dispatch.spawn_process"`
	// DispatchTimeFormat             attributes.String `json:"dispatch.time_format" url:"dispatch.time_format"`
	// DispatchTtl                    attributes.Int    `json:"dispatch.ttl" url:"dispatch.ttl"`
	// DispatchAs                     attributes.String `json:"dispatchAs" url:"dispatchAs"`
	// Displayview                    attributes.String `json:"displayview" url:"displayview"`
	// DurableBackfillType            attributes.String `json:"durable.backfill_type" url:"durable.backfill_type"`
	// DurableLagTime                 attributes.Int    `json:"durable.lag_time" url:"durable.lag_time"`
	// DurableMaxBackfillIntervals    attributes.Int    `json:"durable.max_backfill_intervals" url:"durable.max_backfill_intervals"`
	// DurableTrackTimeType           attributes.String `json:"durable.track_time_type" url:"durable.track_time_type"`
	// IsScheduled                    attributes.Bool   `json:"is_scheduled" url:"is_scheduled"`
	// IsVisible                      attributes.Bool   `json:"is_visible" url:"is_visible"`
	// MaxConcurrent                  attributes.Int    `json:"max_concurrent" url:"max_concurrent"`
	// Name                           attributes.String `json:"name" url:"name"`
	// NextScheduledTime              attributes.String `json:"next_scheduled_time" url:"next_scheduled_time"`
	// QualifiedSearch                attributes.String `json:"qualifiedSearch" url:"qualifiedSearch"`
	// RealtimeSchedule               attributes.Bool   `json:"realtime_schedule" url:"realtime_schedule"`
	// RequestUiDispatchApp           attributes.String `json:"request.ui_dispatch_app" url:"request.ui_dispatch_app"`
	// RequestUiDispatchView          attributes.String `json:"request.ui_dispatch_view" url:"request.ui_dispatch_view"`
	// RestartOnSearchpeerAdd         attributes.Bool   `json:"restart_on_searchpeer_add" url:"restart_on_searchpeer_add"`
	// RunNTimes                      attributes.Int    `json:"run_n_times" url:"run_n_times"`
	// RunOnStartup                   attributes.Bool   `json:"run_on_startup" url:"run_on_startup"`
	// SchedulePriority               attributes.String `json:"schedule_priority" url:"schedule_priority"`
	// ScheduleWindow                 attributes.String `json:"schedule_window" url:"schedule_window"`
	Search attributes.String `json:"search" url:"search"`
	// Vsid                           attributes.String `json:"vsid" url:"vsid"`
	// WorkloadPool                   attributes.String `json:"workload_pool" url:"workload_pool"`

	// Read-only fields are populated by results returned by the Splunk API, but
	// are not settable by Create or Update operations.

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Content
}

// SavedSearch defines a Splunk savedsearch.
type SavedSearch struct {
	client.Namespace   `json:"links" url:"-"`
	client.ACL         `json:"acl" url:"-"`
	client.Title       `json:"name" url:"name"`
	SavedSearchContent `json:"content"`

	// Read-only fields are populated by results returned by the Splunk API, but
	// are not settable by Create or Update operations.
	client.ID `json:"id" url:"-"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Endpoint `endpoint:"saved/searches" json:"endpoint"`
}
