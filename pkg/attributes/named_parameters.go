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

package attributes

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NamedParameters struct {
	Name       string
	Status     string
	Parameters Parameters
}

func (params NamedParameters) StatusBool() (value bool, ok bool) {
	value, err := strconv.ParseBool(params.Status)
	ok = err == nil

	return
}

type NamedParametersCollection []NamedParameters

func (collection NamedParametersCollection) EnabledNames() []string {
	var enabled []string

	for _, params := range collection {
		if isEnabled, ok := params.StatusBool(); ok && isEnabled {
			enabled = append(enabled, params.Name)
		}
	}

	return enabled
}

func UnmarshalJSONForNamedParametersCollection(data []byte, dest interface{}) error {
	destVPtr := reflect.ValueOf(dest)
	if destVPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("attempted UnmarshalJSONForNamedParametersCollection on non-pointer type: %T", dest)
	}

	destV := destVPtr.Elem()
	destT := destV.Type()

	if destT.Kind() != reflect.Struct {
		return fmt.Errorf("attempted UnmarshalJSONForNamedParametersCollection on non-struct type: %T", dest)
	}

	for i := 0; i < destT.NumField(); i++ {
		fieldF := destT.Field(i)
		if !fieldF.IsExported() {
			continue
		}

		fieldTag := fieldF.Tag.Get("npc")
		if fieldTag == "" {
			continue
		}

		var collection NamedParametersCollection
		if fieldF.Type != reflect.TypeOf(collection) {
			return fmt.Errorf("attempted UnmarshalJSONForNamedParametersCollection on non-NamedParametersCollection type %T for field %s", destV.Field(i).Interface(), fieldF.Name)
		}

		var allParams Parameters
		if err := json.Unmarshal(data, &allParams); err != nil {
			return err
		}

		newParams := allParams.withDottedName(fieldTag)

		newCollection := newParams.namedParametersCollection()
		newCollectionV := reflect.ValueOf(newCollection)

		destV.Field(i).Set(newCollectionV)
	}

	return nil
}
