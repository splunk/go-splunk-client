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
	"reflect"
	"strconv"
	"strings"
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

// endpointConfig is a calculated configuration of an endpoint.
type endpointConfig struct {
	// endpoint path
	path string
	// status code that indicates the resource was not found
	codeNotFound int
}

// endpointConfigGetter is the interface for types that include the getEndpointConfig
// method. This interface is intended to be satisfied by having Endpoint as an anonymous
// member of a struct.
type endpointConfigGetter interface {
	getEndpointConfig(endpointConfigGetter) (endpointConfig, error)
}

// endpointTag returns the "endpoint" tag's value for a given endpointConfigGetter. It
// returns an error if the given endpointConfigGetter:
//
// * is not a struct
// * doesn't have a field named "Endpoint"
// * has a field named "Endpoint" that isn't of type Endpoint
// * has no endpoint tag on the Endpoint field
//
// The first three error scenarios are easily avoided by having Endpoint as an anonymous
// field of the struct.
func (e Endpoint) endpointTag(i endpointConfigGetter) (string, error) {
	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		v := reflect.ValueOf(i)
		t = reflect.Indirect(v).Type()
	}

	if t.Kind() != reflect.Struct {
		return "", wrapError(ErrorEndpoint, nil, "unable to determine endpoint for non-struct")
	}

	f, ok := t.FieldByName("Endpoint")
	if !ok {
		return "", wrapError(ErrorEndpoint, nil, "unable to determine endpoint without Endpoint field")
	}

	if f.Type != reflect.TypeOf(e) {
		// it's not entirely valid to require the field named service to be of the service type,
		// but by requiring it we avoid potential confusion, and promote the field being an anonymous member,
		// which is the intention, so that servicePath() is an inherited method.
		return "", wrapError(ErrorEndpoint, nil, "unable to determine endpoint of non-Endpoint field")
	}

	tag := f.Tag.Get("endpoint")
	if tag == "" {
		return "", wrapError(ErrorEndpoint, nil, "unable to determine endpoint without endpoint tag")
	}

	return tag, nil
}

// getEndpointConfig returns the endpointConfig for a given endpointConfigGetter.
func (e Endpoint) getEndpointConfig(i endpointConfigGetter) (endpointConfig, error) {
	tag, err := e.endpointTag(i)
	if err != nil {
		return endpointConfig{}, err
	}

	tagParts := strings.Split(tag, ",")
	// because endpointTag checks for non-empty tag values, this shouldn't be a true
	// necessity, but it feels dirty not checking explicitly that tagParts has at least
	// one value before indexing it.
	if len(tagParts) == 0 {
		return endpointConfig{}, wrapError(ErrorEndpoint, nil, "Endpoint tag has no value")
	}

	newEndpointConfig := endpointConfig{
		path: tagParts[0],
		// defaults
		codeNotFound: 404,
	}

	for _, tagPart := range tagParts[1:] {
		configParts := strings.Split(tagPart, ":")
		if len(configParts) != 2 {
			return endpointConfig{}, wrapError(ErrorEndpoint, nil, "Endpoint tag has invalid configuration: %q", tagPart)
		}

		configPartCode, err := strconv.ParseInt(configParts[1], 10, 0)
		if err != nil {
			return endpointConfig{}, wrapError(ErrorEndpoint, nil, "Endpoint tag has non-integer value: %q", configParts[1])
		}

		switch configParts[0] {
		default:
			return endpointConfig{}, wrapError(ErrorEndpoint, nil, "Endpoint tag has unknown option: %q", configParts[0])
		case "notfound":
			newEndpointConfig.codeNotFound = int(configPartCode)
		}
	}

	return newEndpointConfig, nil
}
