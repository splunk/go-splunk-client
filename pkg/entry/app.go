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

// AppContent is the content for an App.
type AppContent struct {
	Auth            attributes.Explicit[string] `json:"auth"             values:"auth,omitzero"             selective:"create"`
	Author          attributes.Explicit[string] `json:"author"           values:"author,omitzero"`
	Configured      attributes.Explicit[bool]   `json:"configured"       values:"configured,omitzero"`
	Description     attributes.Explicit[string] `json:"description"      values:"description,omitzero"`
	ExplicitAppname attributes.Explicit[string] `json:"explicit_appname" values:"explicit_appname,omitzero" selective:"create"`
	Filename        attributes.Explicit[bool]   `json:"filename"         values:"filename,omitzero"         selective:"create"`
	Label           attributes.Explicit[string] `json:"label"            values:"label,omitzero"`
	Session         attributes.Explicit[string] `json:"session"          values:"session,omitzero"          selective:"create"`
	Version         attributes.Explicit[string] `json:"version"          values:"version,omitzero"`
	Update          attributes.Explicit[bool]   `json:"update"           values:"update,omitzero"           selective:"create"`
	Visible         attributes.Explicit[bool]   `json:"visible"          values:"visible,omitzero"`
}

// App is a Splunk app.
type App struct {
	ID      client.ID  `selective:"create" service:"apps/local"`
	Content AppContent `json:"content" values:",anonymize"`
}
