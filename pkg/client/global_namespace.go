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

// GlobalNamespace represents the global namespace in Splunk. This type should be used
// for types that need to define a namespace, but non-global namespaces don't make sense
// (such as for User).
type GlobalNamespace struct{}

// NamespacePath returns the namespace path.
func (c GlobalNamespace) NamespacePath() (string, error) {
	return Namespace{}.NamespacePath()
}

// SetNamespace is a required method for the NamespacePather interface, but for a GlobalNamespace
// doesn't actually do anything other than check that user and app are empty.
func (c *GlobalNamespace) SetNamespace(user string, app string) error {
	if user != "" || app != "" {
		return wrapError(ErrorNamespace, nil, "attempted to set non-empty namespace on a GlobalNamespace")
	}

	return nil
}
