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

// Title represents the title of an object in Splunk.
type Title string

// TitleValue returns the string value of the Title.
func (t Title) TitleValue() string {
	return string(t)
}

func (t Title) HasTitle() bool {
	return t != ""
}

// Titler defines methods that a type must implement to be a titled object.
type Titler interface {
	// Title returns the string representation of the object's Title.
	TitleValue() string
	// HasTitle returns a boolean indicating if there is a non-empty Title.
	HasTitle() bool
}
