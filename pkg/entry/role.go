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

// RoleContent defines the Content for a Role.
type RoleContent struct {
	Capabilities              attributes.Strings `json:"capabilities"              url:"capabilities"`
	CumulativeRTSrchJobsQuota attributes.Int     `json:"cumulativeRTSrchJobsQuota" url:"cumulativeRTSrchJobsQuota"`
	CumulativeSrchJobsQuota   attributes.Int     `json:"cumulativeSrchJobsQuota"   url:"cumulativeSrchJobsQuota"`
	DefaultApp                attributes.String  `json:"defaultApp"                url:"defaultApp"`
	ImportedRoles             attributes.Strings `json:"imported_roles"            url:"imported_roles"`
	RtSrchJobsQuota           attributes.Int     `json:"rtSrchJobsQuota"           url:"rtSrchJobsQuota"`
	SrchDiskQuota             attributes.Int     `json:"srchDiskQuota"             url:"srchDiskQuota"`
	SrchFilter                attributes.String  `json:"srchFilter"                url:"srchFilter"`
	SrchIndexesAllowed        attributes.Strings `json:"srchIndexesAllowed"        url:"srchIndexesAllowed"`
	SrchIndexesDefault        attributes.Strings `json:"srchIndexesDefault"        url:"srchIndexesDefault"`
	SrchJobsQuota             attributes.Int     `json:"srchJobsQuota"             url:"srchJobsQuota"`
	SrchTimeWin               attributes.Int     `json:"srchTimeWin"               url:"srchTimeWin"`

	// Read-only fields are populated by results returned by the Splunk API, but
	// are not settable by Create or Update operations.
	ImportedCapabilities       attributes.Strings `json:"imported_capabilities"       url:"-"`
	ImportedRtSrchJobsQuota    attributes.Int     `json:"imported_rtSrchJobsQuota"    url:"-"`
	ImportedRtSrchJObsQuota    attributes.Int     `json:"imported_rtSrchJObsQuota"    url:"-"`
	ImportedSrchDiskQuota      attributes.Int     `json:"imported_srchDiskQuota"      url:"-"`
	ImportedSrchFilter         attributes.String  `json:"imported_srchFilter"         url:"-"`
	ImportedSrchIndexesAllowed attributes.Strings `json:"imported_srchIndexesAllowed" url:"-"`
	ImportedSrchIndexesDefault attributes.Strings `json:"imported_srchIndexesDefault" url:"-"`
	ImportedSrchJobsQuota      attributes.Int     `json:"imported_srchJobsQuota"      url:"-"`
	ImportedSrchTimeWin        attributes.Int     `json:"imported_srchTimeWin"        url:"-"`

	// The below fields don't have values, and only exist to provide context to
	// the Splunk API.
	client.Content
}

// Role defines a Splunk role.
type Role struct {
	ID          client.ID `service:"authorization/roles"`
	RoleContent `json:"content"`
}
