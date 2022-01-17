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