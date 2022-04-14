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

package selective_test

import (
	"encoding/json"
	"fmt"

	"github.com/splunk/go-splunk-client/pkg/selective"
)

type Discography struct {
	Albums  []string `json:"albums" selective:"albums"`
	Singles []string `json:"singles" selective:"singles"`
}

type Artist struct {
	// Name defines no selective tag, so it is always selected
	Name string `json:"name"`

	// Discography is selected for both albums and singles
	Discography Discography `json:"discography" selective:"albums,singles"`

	// Members is selected only for personnel
	Members []string `json:"members" selective:"personnel"`
}

func ExampleEncode() {
	artist := Artist{
		Name: "The Refreshments",
		Discography: Discography{
			Albums: []string{
				"Fizzy Fuzzy Big and Buzzy",
				"The Bottle and Fresh Horses",
			},
			Singles: []string{
				"Yahoos and Triangles",
			},
		},
		Members: []string{
			"Roger Clyne",
			"Dusty Denham",
			"P.H. Naffah",
			"Buddy Edwards",
			"Brian David Blush",
		},
	}

	artistAlbumsOnly, _ := selective.Encode(artist, "albums")
	artistAlbumsOnlyJSON, _ := json.MarshalIndent(artistAlbumsOnly, "", "  ")

	fmt.Printf("%s", artistAlbumsOnlyJSON)
	// Output: {
	//   "name": "The Refreshments",
	//   "discography": {
	//     "albums": [
	//       "Fizzy Fuzzy Big and Buzzy",
	//       "The Bottle and Fresh Horses"
	//     ]
	//   }
	// }
}
