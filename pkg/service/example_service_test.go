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

package service_test

import (
	"fmt"
	"net/url"

	"github.com/splunk/go-splunk-client/pkg/service"
)

// Title is a string that uniquely identifies a resource.
type Title string

// GetServicePath implements custom GetServicePath encoding. It returns the given
// path back, which has the effect of using the Title field's struct tag as
// its GetServicePath.
func (t Title) GetServicePath(path string) (string, error) {
	return path, nil
}

// GetEntryPath implements custom GetEntryPath encoding. It returns the url-encoded
// value of the Title with the given path preceding it.
func (t Title) GetEntryPath(path string) (string, error) {
	servicePath, err := t.GetServicePath(path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", servicePath, url.PathEscape(string(t))), nil
}

// Artist represents a musical artist.
type Artist struct {
	// Name is a Title which determines both the ServicePath (music/artists)
	// and EntryPath (music/artists/<Title>).
	Name Title `service:"music/artists"`
}

func Example() {
	newArtist := Artist{Name: "The Refreshments"}

	servicePath, _ := service.ServicePath(newArtist)
	fmt.Printf("path to create Artist: %s\n", servicePath)
	entryPath, _ := service.EntryPath(newArtist)
	fmt.Printf("path to existing Artist: %s\n", entryPath)
	// Output: path to create Artist: music/artists
	// path to existing Artist: music/artists/The%20Refreshments
}
