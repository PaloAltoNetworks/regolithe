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

package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"go.aporeto.io/regolithe/cmd/rego/static"
	"go.aporeto.io/regolithe/spec"
)

func convertType(t spec.AttributeType) string {
	switch t {
	case spec.AttributeTypeString:
		return "string"
	case spec.AttributeTypeInt:
		return "integer"
	case spec.AttributeTypeFloat:
		return "float"
	case spec.AttributeTypeBool:
		return "boolean"
	case spec.AttributeTypeTime:
		return "time"
	case spec.AttributeTypeObject:
		return "object"
	case spec.AttributeTypeRefList, spec.AttributeTypeList:
		return "$list"
	case spec.AttributeTypeRefMap:
		return "$map"
	case spec.AttributeTypeRef:
		return "$ref"
	case spec.AttributeTypeExt:
		return "$external"
	default:
	}
	return "UNRECOGNIZED_TYPE/" + string(t)
}

func convertRegexp(str string, required bool) string {

	escaped := strings.Replace(strings.Replace(str, `\`, `\\`, -1), `"`, `\"`, -1)
	if required {
		return escaped
	}
	return fmt.Sprintf("(%s)?", escaped)
}

func isNil(target interface{}) bool {
	return target == nil
}

func jsonStringify(target interface{}) string {
	data, err := json.Marshal(target)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(data)
}

func stripFirstLevelBrackets(str string) string {
	return strings.TrimSuffix(strings.TrimPrefix(str, "{"), "}")
}

func makeTemplate(p string) (*template.Template, error) {

	data, err := static.Asset(p)
	if err != nil {
		return nil, err
	}

	return template.New(path.Base(p)).Funcs(functions).Parse(string(data))
}

func writeFile(path string, data []byte) error {

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to write file: %s", f.Name())
	}

	// #nosec G307
	defer f.Close() // nolint: errcheck
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("unable to write file: %s", f.Name())
	}

	return nil
}
