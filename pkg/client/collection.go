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
	"go-sdk/pkg/client/internal/immutable"
	"net/http"
	"reflect"
)

type collection interface {
	collectionPath(collectionAttributes) (string, error)
	entryPath(collectionAttributes, immutable.Name) (string, error)
}

type collectionAttributes struct {
	path      string
	namespace Namespace
}

func pathForCollection(v reflect.Value) (string, error) {
	t := v.Type()

	f, ok := t.FieldByName("path")
	if !ok {
		return "", fmt.Errorf("collection path field missing")
	}

	tag := f.Tag.Get("collection")
	if tag == "" {
		return "", fmt.Errorf("collection path field missing collection tag")
	}

	return tag, nil
}

func namespaceForCollection(v reflect.Value) (Namespace, error) {
	nsV := v.FieldByName("Namespace")
	ns := Namespace{}

	if nsV.IsValid() {
		if nsV.Type() != reflect.TypeOf(Namespace{}) {
			return Namespace{}, fmt.Errorf("collectionForInterface() passed struct with Namespace field that is not a Namespace{}")
		}

		ns = nsV.Interface().(Namespace)
	}

	return ns, nil
}

func collectionForInterface(i interface{}) (collectionAttributes, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Struct {
		return collectionAttributes{}, fmt.Errorf("collectionForInterface() passed a non-struct")
	}

	newCollection := collectionAttributes{}

	path, err := collectionPath(v)
	if err != nil {
		return collectionAttributes{}, err
	}
	newCollection.path = path

	ns, err := collectionNamespace(v)
	if err != nil {
		return collectionAttributes{}, err
	}
	newCollection.namespace = ns

	return newCollection, nil
}

func (c collectionAttributes) readRequest(client *Client) (*http.Request, error) {
	u, err := client.urlForPath(c.namespace, c.path)
	if err != nil {
		return nil, err
	}

	r := &http.Request{
		URL: u,
	}

	if err := client.Authenticator.authenticateRequest(client, r); err != nil {
		return nil, err
	}

	return r, nil
}
