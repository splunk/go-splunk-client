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

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			log.Fatalf("line didn't have exactly two parts: %s\n", line)
		}

		fieldName, fieldType := parts[0], parts[1]

		tag := fmt.Sprintf(`json:"%s" url:"%s"`, fieldName, fieldName)

		re := regexp.MustCompile("[._]")
		fieldNameParts := re.Split(fieldName, -1)
		for i, fieldNamePart := range fieldNameParts {
			fieldNameParts[i] = strings.Title(fieldNamePart)
		}
		fieldName = strings.Join(fieldNameParts, "")

		switch fieldType {
		case "Bool", "Boolean":
			fieldType = "attributes.Bool"
		case "String":
			fieldType = "attributes.String"
		case "Number":
			fieldType = "attributes.Int"
		}

		fmt.Printf("%s %s `%s`\n", fieldName, fieldType, tag)
	}
}
