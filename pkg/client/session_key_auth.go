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

import (
	"fmt"
	"net/http"
)

// SessionKeyAuth passes SessionKey to the REST API via the Authorization header.
type SessionKeyAuth struct {
	SessionKey string `xml:"sessionKey"`
}

// authenticateRequest adds the SessionKey to the http.Request's Header.
func (s *SessionKeyAuth) authenticateRequest(c *Client, r *http.Request) error {
	if s.SessionKey == "" {
		return fmt.Errorf("attempted to authenticate request with empty SessionKey")
	}

	if r.Header == nil {
		r.Header = http.Header{}
	}

	r.Header.Add("Authorization", fmt.Sprintf("Splunk %s", s.SessionKey))

	return nil
}
