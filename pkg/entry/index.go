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
	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/client"
)

// IndexContent is the content for an Index.
type IndexContent struct {
	ColdToFrozenDir        attributes.Explicit[string] `json:"coldToFrozenDir"        values:"coldToFrozenDir,omitempty"`
	ColdToFrozenScript     attributes.Explicit[string] `json:"coldToFrozenScript"     values:"coldToFrozenScript,omitempty"`
	DataType               attributes.Explicit[string] `json:"datatype"               values:"datatype,omitempty"               selective:"create"`
	Disabled               attributes.Explicit[bool]   `json:"disabled"               values:"disabled,omitempty"               selective:"read"`
	FrozenTimePeriodInSecs attributes.Explicit[int]    `json:"frozenTimePeriodInSecs" values:"frozenTimePeriodInSecs,omitempty"`
	HomePath               attributes.Explicit[string] `json:"homePath"               values:"homePath,omitempty"               selective:"create"`
	MaxDataSize            attributes.Explicit[string] `json:"maxDataSize"            values:"maxDataSize,omitempty"`
	MaxHotBuckets          attributes.Explicit[string] `json:"maxHotBuckets"          values:"maxHotBuckets,omitempty"`
	MaxHotIdleSecs         attributes.Explicit[int]    `json:"maxHotIdleSecs"         values:"maxHotIdleSecs,omitempty"`
	MaxHotSpanSecs         attributes.Explicit[int]    `json:"maxHotSpanSecs"         values:"maxHotSpanSecs,omitempty"`
	MaxTotalDataSizeMB     attributes.Explicit[int]    `json:"maxTotalDataSizeMB"     values:"maxTotalDataSizeMB,omitempty"`
	MaxWarmDBCount         attributes.Explicit[int]    `json:"maxWarmDBCount"         values:"maxWarmDBCount,omitempty"`
	QuarantineFutureSecs   attributes.Explicit[int]    `json:"quarantineFutureSecs"   values:"quarantineFutureSecs,omitempty"`
	QuarantinePastSecs     attributes.Explicit[int]    `json:"quarantinePastSecs"     values:"quarantinePastSecs,omitempty"`
	ThawedPath             attributes.Explicit[string] `json:"thawedPath"             values:"thawedPath,omitempty"            selective:"create"`
}

// Index is a Splunk Index.
type Index struct {
	ID      client.ID    `selective:"create" service:"data/indexes"`
	Content IndexContent `json:"content" values:",anonymize"`
}
