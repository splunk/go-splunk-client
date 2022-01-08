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

// message represents the <msg> element of a <messages> entry.
type message struct {
	Value string `xml:",chardata"`
	Code  string `xml:"code,attr"`
}

// messages represents the <messages> element of a <response> entry.
type messages struct {
	XMLName string    `xml:"messages"`
	Items   []message `xml:"msg"`
}

// firstAndOnly returns the first message if exactly one message is present. Otherwise
// it returns an empty message and false.
func (m messages) firstAndOnly() (message, bool) {
	if len(m.Items) != 1 {
		return message{}, false
	}

	return m.Items[0], true
}
