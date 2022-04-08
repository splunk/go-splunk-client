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

func ExampleExplicit_strings() {
	type knowledgeObject struct {
		Name   string                        `values:"name"`
		Values attributes.Explicit[[]string] `values:"values,omitempty"`
	}

	myObject := knowledgeObject{
		Name: "my_knowledge_object",
	}
	// myObjectURLValues will not have a value for Values as it has not been set
	myObjectURLValues, _ := values.Encode(myObject)
	fmt.Printf("myObjectURLValues without explicitly set Values: %s\n", myObjectURLValues)

	myObject.Values = attributes.NewExplicit([]string{})
	// myObjectURLValues will have a value of an empty string for Values as it has been explicitly set empty
	myObjectURLValues, _ = values.Encode(myObject)
	fmt.Printf("myObjectURLValues with explicitly set (empty) Values: %s\n", myObjectURLValues)

	myObject.Values = attributes.NewExplicit([]string{"valueA", "valueB", "valueC"})
	// myObjectURLValues will have multiple values for Values
	myObjectURLValues, _ = values.Encode(myObject)
	fmt.Printf("myObjectURLValues with explicitly set (non-empty) Values: %s\n", myObjectURLValues)

	// Output: myObjectURLValues without explicitly set Values: map[name:[my_knowledge_object]]
	// myObjectURLValues with explicitly set (empty) Values: map[name:[my_knowledge_object] values:[]]
	// myObjectURLValues with explicitly set (non-empty) Values: map[name:[my_knowledge_object] values:[valueA valueB valueC]]
}
