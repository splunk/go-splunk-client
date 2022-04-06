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

package deepset_test

import (
	"fmt"

	"github.com/splunk/go-splunk-client/pkg/deepset"
)

type Namespace struct {
	User string
	App  string
}

type ID struct {
	Title     string
	Namespace Namespace
}

type View struct {
	ID      ID
	Content string
}

func Example_set() {
	myView := View{}

	// returned error ignored here
	_ = deepset.Set(&myView, Namespace{User: "admin", App: "search"})

	fmt.Printf("ID User: %s, App: %s", myView.ID.Namespace.User, myView.ID.Namespace.App)
	// Output: ID User: admin, App: search
}
