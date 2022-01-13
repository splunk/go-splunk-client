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

// service is a type that permits a type to define its service path via struct tag. It
// is expected to be added as an anonymous member to an object type. Being anonymous
// simplifies that type adhering to the servicer interface.
//
// To be valid, the `service` tag for this field must include a value that defines the
// service's API path.
//
// For an imaginary Splunk object type whose service path is `things/widgets`, the
// service field would be used as such:
//
//   type Widget struct {
//     service `service:"things/widgets"`
//   }
type service struct{}

// servicePath returns the "service" tag's value for a given interface. It returns an error
// if the given interface:
//
// * is not a struct
// * doesn't have a field named "service"
// * has a field named "service" that isn't of type service
// * has no service tag on the service field
func (s service) servicePath(i interface{}) (string, error) {
	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		v := reflect.ValueOf(i)
		t = reflect.Indirect(v).Type()
	}

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("unable to determine servicePath for non-struct")
	}

	f, ok := t.FieldByName("service")
	if !ok {
		return "", fmt.Errorf("unable to determine servicePath without service field")
	}

	if f.Type != reflect.TypeOf(s) {
		// it's not entirely valid to require the field named service to be of the service type,
		// but by requiring it we avoid potential confusion, and promote the field being an anonymous member,
		// which is the intention, so that servicePath() is an inherited method.
		return "", fmt.Errorf("unable to determine servicePath of non-service field")
	}

	tag := f.Tag.Get("service")
	if tag == "" {
		return "", fmt.Errorf("unable to determine servicePath without service tag")
	}

	return tag, nil
}
