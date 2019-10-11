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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"text/template"

	"go.aporeto.io/regolithe/spec"
)

var functions = template.FuncMap{
	"convertType":             convertType,
	"convertRegexp":           convertRegexp,
	"jsonStringify":           jsonStringify,
	"isNil":                   isNil,
	"stripFirstLevelBrackets": stripFirstLevelBrackets,
}

func writeGlobalResources(set spec.SpecificationSet, outFolder string) error {

	if err := os.MkdirAll(outFolder, 0750); err != nil && !os.IsExist(err) {
		return err
	}

	tmpl, err := makeTemplate("templates/json-schema-restname.gotpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(
		&buf,
		struct {
			Name string
			Set  spec.SpecificationSet
		}{
			Name: set.Configuration().Name,
			Set:  set,
		}); err != nil {
		return fmt.Errorf("unable to generate global resource code: %s", err)
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return fmt.Errorf("unable to unmarshal model code: %s", err)
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal model code: %s", err)
	}

	return writeFile(path.Join(outFolder, "_models.json"), out)
}

func writeGlobalResourceLists(set spec.SpecificationSet, outFolder string) error {

	if err := os.MkdirAll(outFolder, 0750); err != nil && !os.IsExist(err) {
		return err
	}

	tmpl, err := makeTemplate("templates/json-schema-resourcename.gotpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(
		&buf,
		struct {
			Name string
			Set  spec.SpecificationSet
		}{
			Name: set.Configuration().Name,
			Set:  set,
		}); err != nil {
		return fmt.Errorf("unable to generate global resource lists code: %s", err)
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return fmt.Errorf("unable to unmarshal model code: %s", err)
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal model code: %s", err)
	}

	return writeFile(path.Join(outFolder, "_lists.json"), out)
}

func writeModel(set spec.SpecificationSet, outFolder string) error {

	if err := os.MkdirAll(outFolder, 0750); err != nil && !os.IsExist(err) {
		return err
	}

	tmpl, err := makeTemplate("templates/json-schema.gotpl")
	if err != nil {
		return err
	}

	for _, s := range set.Specifications() {
		var buf bytes.Buffer

		if err = tmpl.Execute(
			&buf,
			struct {
				Name string
				Spec spec.Specification
			}{
				Name: set.Configuration().Name,
				Spec: s,
			}); err != nil {
			return fmt.Errorf("unable to generate model code: %s", err)
		}

		data := map[string]interface{}{}
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			return fmt.Errorf("unable to unmarshal model code: %s", err)
		}

		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal model code: %s", err)
		}

		if err := writeFile(path.Join(outFolder, s.Model().RestName+".json"), out); err != nil {
			return err
		}
	}

	return nil
}
