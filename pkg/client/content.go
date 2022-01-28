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
)

// ContentGetter types implement the GetContent() function to return the
// interface that represents their content.
type ContentGetter interface {
	GetContent(ContentGetter) ContentGetter
}

// Content is an empty data type that provides the content method to satisfy
// the contenter interface. It is the prescribed means of having a struct
// properly adhere to the Contenter type.
//
// Example:
//   type SomeTypeContent struct {
//       Data string
//       // Content is an anonymous field, satisfies the ContentGetter interface
//       // and identifies SomeTypeContent as the content for types that have
//       // SomeTypeContent as an anonymous field.
//       client.Content
//   }
//   type SomeType struct {
//       Name string
//       // SomeTypeContent is this type's "content"
//       SomeTypeContent
//   }
type Content struct{}

// GetContent returns the provided Contenter's content. Its content is determined
// by looking through the Contenter's struct fields to find the outermost level
// that has Content as an anonymous field.
func (c Content) GetContent(contenter ContentGetter) ContentGetter {
	t := reflect.TypeOf(contenter)

	// if i is a Content object, return nil
	if t == reflect.TypeOf(Content{}) {
		return nil
	}

	// get value pointed to, if pointer was passed
	if t.Kind() == reflect.Ptr {
		v := reflect.Indirect(reflect.ValueOf(contenter))
		t = v.Type()
	}

	// content always expected to be a struct
	if t.Kind() != reflect.Struct {
		return nil
	}

	// look for Content as an anonymous struct member, because that identifes
	// which level should be considered the content.
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.Anonymous {
			continue
		}

		if f.Type == reflect.TypeOf(c) {
			return contenter
		}
	}

	contenterType := reflect.TypeOf((*ContentGetter)(nil)).Elem()
	var foundContenter ContentGetter

	// if we haven't found the Content field yet, iterate through all struct
	// fields and recurse to find the embedded item that has it.
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		fieldV := reflect.ValueOf(contenter).Field(i)
		if fieldV.Type().Implements(contenterType) {
			fieldContent := c.GetContent(fieldV.Interface().(ContentGetter))
			if fieldContent != nil {
				// multiple fields having their own Content fields is ambiguous as to which
				// should be used, so return nil indicating none found. this really shouldn't
				// be a likely situation, because Go won't expose ambigious methods from anonymous
				// fields, but checking for it here just in case.
				if foundContenter != nil {
					return nil
				}

				foundContenter = fieldContent
			}
		}
	}

	// foundContenter may be nil, if non-ambiguous Content not found
	return foundContenter
}
