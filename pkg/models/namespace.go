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

import (
	"fmt"
	"path"
)

// Namespace represents a Splunk namespace/context. If both User/App are empty, the
// global context will be used. Both User and App must be unset or set. It is invalid
// for one to be set without the other.
type Namespace struct {
	User string
	App  string
}

// GlobalNamespace has an empty User/App.
var GlobalNamespace = Namespace{}

// validate returns an error if Namespace is invalid. It is invalid if either User or App
// is set without the other.
func (c Namespace) validate() error {
	if c.App != "" && c.User != "" {
		return nil
	} else if c.App == "" && c.User == "" {
		return nil
	}

	return fmt.Errorf("invalid Context, neither or both App/User must be set")
}

// Path returns the relative URL path for the Namespace.
func (c Namespace) Path() (string, error) {
	if err := c.validate(); err != nil {
		return "", err
	}

	// absence of either field indicates global context
	if c.App == "" {
		return "services", nil
	}

	return path.Join("servicesNS", c.User, c.App), nil
}
