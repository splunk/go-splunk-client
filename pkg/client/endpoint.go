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
	"reflect"
)

// Endpoint is a type that permits a type to define its Endpoint path via struct tag. It
// is expected to be added as an anonymous member to an object type. Being anonymous
// simplifies that type adhering to the Service interface.
//
// To be valid, the `endpoint` tag for this field must include a value that defines the
// Endpoint's API path.
//
// For an imaginary Splunk object type whose Endpoint path is `things/widgets`, the
// Endpoint field would be used as such:
//
//   type Widget struct {
//     Endpoint `endpoint:"things/widgets"`
//   }
type Endpoint struct{}

// endpointPath returns the "endpoint" tag's value for a given interface. It returns an error
// if the given interface:
//
// * is not a struct
// * doesn't have a field named "Endpoint"
// * has a field named "Endpoint" that isn't of type Endpoint
// * has no endpoint tag on the Endpoint field
func (e Endpoint) endpointPath(i interface{}) (string, error) {
	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		v := reflect.ValueOf(i)
		t = reflect.Indirect(v).Type()
	}

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("unable to determine endpoint for non-struct")
	}

	f, ok := t.FieldByName("Endpoint")
	if !ok {
		return "", fmt.Errorf("unable to determine endpoint without Endpoint field")
	}

	if f.Type != reflect.TypeOf(e) {
		// it's not entirely valid to require the field named service to be of the service type,
		// but by requiring it we avoid potential confusion, and promote the field being an anonymous member,
		// which is the intention, so that servicePath() is an inherited method.
		return "", fmt.Errorf("unable to determine endpoint of non-Endpoint field")
	}

	tag := f.Tag.Get("endpoint")
	if tag == "" {
		return "", fmt.Errorf("unable to determine endpoint without endpoint tag")
	}

	return tag, nil
}
