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

type Roles struct {
	Namespace
	Entries []Role `json:"entry"`
}

func (r Roles) collectionPath() string {
	return "authorization/roles"
}

func (r Roles) EntryWithTitle(title string) (Entry, bool) {
	for _, entry := range r.Entries {
		if entry.Title() == title {
			return entry, true
		}
	}

	return nil, false
}
