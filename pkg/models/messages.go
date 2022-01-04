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

package models

// MessageElement represents the <msg> element of a <messages> entry.
type MessageElement struct {
	Message string `xml:",chardata"`
	Code    string `xml:"code,attr"`
}

// MessagesElement represents the <messages> element of a <response> entry.
type MessagesElement struct {
	MessageElements []MessageElement `xml:"msg"`
}
