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

package messages

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Messages represents the <messages> element of a <response> entry.
type Messages struct {
	XMLName string    `xml:"messages"`
	Items   []Message `xml:"msg" json:"messages"`
}

func NewMessagesFromResponse(r *http.Response) (Messages, error) {
	m := Messages{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&m); err != nil {
		return Messages{}, err
	}

	return m, nil
}

// FirstAndOnly returns the first message if exactly one message is present. Otherwise
// it returns an empty message and false.
func (m Messages) FirstAndOnly() (Message, bool) {
	if len(m.Items) != 1 {
		return Message{}, false
	}

	return m.Items[0], true
}

// String returns the string representation of Messages. If multiple Message items are
// present, they will be comma-separated.
func (m Messages) String() string {
	itemStrings := make([]string, len(m.Items))
	for i := 0; i < len(m.Items); i++ {
		itemStrings[i] = m.Items[i].String()
	}

	return strings.Join(itemStrings, ",")
}
