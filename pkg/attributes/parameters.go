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
	"sort"
	"strings"
)

// Parameters is a map of parameter names to string values.
type Parameters map[string]string

func dottedParameterParts(fullFieldName string) (dottedFieldName string, dottedParamName string) {
	parts := strings.Split(fullFieldName, ".")

	return parts[0], strings.Join(parts[1:], ".")
}

func (p Parameters) dottedNames() []string {
	foundNamesMap := map[string]bool{}
	var foundNames []string

	for key := range p {
		fieldName, _ := dottedParameterParts(key)

		if _, ok := foundNamesMap[fieldName]; !ok {
			foundNames = append(foundNames, fieldName)
		}
		foundNamesMap[fieldName] = true
	}

	sort.Strings(foundNames)

	return foundNames
}

// withDottedName returns a new Parameters object containing the nested parameters
// for the given name. The new Parameters name field will have this name removed.
//
// For example, this input Parameters:
//
//   Parameters{
// 	   "action.email":    "true",
// 	   "action.email.to": "whoever@example.com",
//   }
//
// would result in the below Parameters for `withDottedName("action")`:
//
//   Parameters{
// 	   "email":    "true",
// 	   "email.to": "whoever@example.com",
//   }
func (p Parameters) withDottedName(name string) Parameters {
	var newParameters Parameters

	for key, value := range p {
		fieldName, fieldParamName := dottedParameterParts(key)

		if fieldName == name && fieldParamName != "" {
			if newParameters == nil {
				newParameters = Parameters{}
			}

			newParameters[fieldParamName] = value
		}
	}

	return newParameters
}

func (p Parameters) namedParametersWithDottedName(name string) NamedParameters {
	return NamedParameters{
		Name:       name,
		Status:     p[name],
		Parameters: p.withDottedName(name),
	}
}

func (p Parameters) namedParametersCollection() NamedParametersCollection {
	names := p.dottedNames()
	var newCollection NamedParametersCollection

	for _, name := range names {
		newCollection = append(newCollection, p.namedParametersWithDottedName(name))
	}

	return newCollection
}

func (p *Parameters) UnmarshalJSON(data []byte) error {
	interfaceMap := map[string]interface{}{}
	if err := json.Unmarshal(data, &interfaceMap); err != nil {
		return err
	}

	if len(interfaceMap) == 0 {
		return nil
	}

	newP := Parameters{}
	for key, value := range interfaceMap {
		switch reflect.TypeOf(value).Kind() {
		default:
			return fmt.Errorf("unable to unmarshal unhandled type %T into Parameters for key %s", value, key)
		case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			newP[key] = fmt.Sprintf("%v", value)
		}
	}

	*p = newP

	return nil
}

func UnmarshalJSONForParameters(data []byte, dest interface{}) error {
	destVPtr := reflect.ValueOf(dest)
	if destVPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("attempted UnmarshalJSONForParameters on non-pointer type: %T", dest)
	}

	destV := destVPtr.Elem()
	destT := destV.Type()

	if destT.Kind() != reflect.Struct {
		return fmt.Errorf("attempted UnmarshalJSONForParameters on non-struct type: %T", dest)
	}

	for i := 0; i < destT.NumField(); i++ {
		fieldF := destT.Field(i)
		if !fieldF.IsExported() {
			continue
		}

		fieldTag := fieldF.Tag.Get("parameters")
		if fieldTag == "" {
			continue
		}

		allParams := make(Parameters)
		if fieldF.Type != reflect.TypeOf(allParams) {
			return fmt.Errorf("attempted UnmarshalJSONForParameters on non-Parameters type %T for field %s", destV.Field(i).Interface(), fieldF.Name)
		}

		if err := json.Unmarshal(data, &allParams); err != nil {
			return err
		}

		newParams := allParams.withDottedName(fieldTag)
		newParamsV := reflect.ValueOf(newParams)

		destV.Field(i).Set(newParamsV)
	}

	return nil
}
