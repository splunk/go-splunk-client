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

package attributes_test

import (
	"fmt"

	"github.com/splunk/go-splunk-client/pkg/attributes"
	"github.com/splunk/go-splunk-client/pkg/values"
)

func ExampleExplicit_bool() {
	type knowledgeObject struct {
		Name     string                    `values:"name"`
		Disabled attributes.Explicit[bool] `values:"disabled,omitzero"`
	}

	myObject := knowledgeObject{
		Name: "my_knowledge_object",
	}
	// myObjectURLValues will not have a value for Disabled as it has not been set
	myObjectURLValues, _ := values.Encode(myObject)
	fmt.Printf("myObjectURLValues without explicitly set Disabled: %s\n", myObjectURLValues)

	myObject.Disabled = attributes.NewExplicit(false)
	// myObjectURLValues will have a value of false for Disabled as it has been set
	myObjectURLValues, _ = values.Encode(myObject)
	fmt.Printf("myObjectURLValues with explicitly set Disabled: %s\n", myObjectURLValues)
	// Output: myObjectURLValues without explicitly set Disabled: map[name:[my_knowledge_object]]
	// myObjectURLValues with explicitly set Disabled: map[disabled:[false] name:[my_knowledge_object]]
}
