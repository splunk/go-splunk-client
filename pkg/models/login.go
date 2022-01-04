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

// Login represents the fields used to authenticate to Splunk.
type Login struct {
	Username string `url:"username"`
	Password string `url:"password"`
}

// LoginResponseElement represents the <response> element for an authentication request's reseponse.
type LoginResponseElement struct {
	Messages   MessagesElement `xml:"messages"`
	SessionKey string          `xml:"sessionKey"`
}
