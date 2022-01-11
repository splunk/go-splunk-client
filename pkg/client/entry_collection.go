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
	"log"
	"reflect"
)

// entryCollection types represent collections of entries.
type entryCollection interface {
	collectionPath() string
	firstAndOnlyEntry() (User, error)
}

func firstAndOnlyEntry(entries interface{}) (entry, error) {
	entriesV := reflect.ValueOf(entries)
	if entriesV.Kind() != reflect.Slice {
		return nil, fmt.Errorf("entryCollection.Entries is not a slice")
	}

	if entriesV.Len() == 0 {
		return nil, fmt.Errorf("no entries present")
	}

	if entriesV.Len() > 1 {
		return nil, fmt.Errorf("more than one entry present, which should never happen")
	}

	foundEntryV := entriesV.Index(0)

	entryType := reflect.TypeOf((*entry)(nil)).Elem()
	if !foundEntryV.Type().Implements(entryType) {
		return nil, fmt.Errorf("non-entry value found")
	}

	return foundEntryV.Interface().(entry), nil
}
