// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Pluralize pluralizes the given word.
func Pluralize(word string) string {

	if len(word) == 0 {
		return word
	}

	lastChar := word[len(word)-1:]
	if lastChar == "s" {
		return word
	}

	if len(word) >= 2 {
		secondLastChar := word[len(word)-2 : len(word)-1]

		if lastChar == "y" &&
			secondLastChar != "a" &&
			secondLastChar != "e" &&
			secondLastChar != "i" &&
			secondLastChar != "o" &&
			secondLastChar != "u" &&
			secondLastChar != "y" {
			return word[:len(word)-1] + "ies"
		}
	}

	return word + "s"
}

func makeSchemaValidationError(message string, res []gojsonschema.ResultError) []error {

	out := make([]error, len(res))
	for i := range res {
		out[i] = fmt.Errorf("%s: schema error: %s", message, res[i].String())
	}

	return out
}

func formatValidationErrors(errs []error) error {

	if len(errs) == 0 {
		return nil
	}

	out := make([]string, len(errs))
	for i := range errs {
		out[i] = errs[i].Error()
	}

	sort.Strings(out)

	return errors.New(strings.Join(out, "\n"))
}

func sortVersionStrings(versions []string) []string {

	vs := make([]int, len(versions))

	for i := range versions {

		v := versions[i]
		currentVersion, err := versionToInt(v)
		if err != nil {
			panic(fmt.Sprintf("invalid version '%s'", v))
		}

		vs[i] = currentVersion
	}

	sort.Ints(vs)

	out := make([]string, len(vs))
	for i := range vs {
		out[i] = fmt.Sprintf("v%d", vs[i])
	}

	return out
}

func sortAttributes(attrs []*Attribute) {

	sort.Slice(attrs, func(i int, j int) bool {
		return strings.Compare(attrs[i].Name, attrs[j].Name) == -1
	})
}

func sortParameters(params []*Parameter) {

	sort.Slice(params, func(i int, j int) bool {
		return strings.Compare(params[i].Name, params[j].Name) == -1
	})
}

func versionToInt(version string) (int, error) {

	return strconv.Atoi(strings.TrimPrefix(version, "v"))
}
