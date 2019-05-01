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
	"strings"

	"github.com/fatih/structs"
	wordwrap "github.com/mitchellh/go-wordwrap"
	yaml "gopkg.in/yaml.v2"
)

func toYAMLMapSlice(s interface{}) yaml.MapSlice {

	var out yaml.MapSlice

	for _, field := range structs.Fields(s) {

		if !field.IsExported() {
			continue
		}

		yamlName, omit := splitTags(field.Tag("yaml"))
		if yamlName == "" {
			continue
		}

		if (field.IsZero() || field.Value() == nil) && omit {
			continue
		}

		var v interface{}
		if yamlName == "description" {
			v = wordwrap.WrapString(field.Value().(string), 80)
		} else {
			v = field.Value()
		}

		item := yaml.MapItem{
			Key:   yamlName,
			Value: v,
		}

		out = append(out, item)
	}

	return out
}

func splitTags(tag string) (string, bool) {

	if tag == "" || tag == "-" {
		return "", false
	}

	parts := strings.Split(tag, ",")

	if len(parts) == 1 {
		return parts[0], false
	}

	return parts[0], parts[1] == "omitempty"
}
