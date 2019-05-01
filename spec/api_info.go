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
	"os"

	"github.com/xeipuuv/gojsonschema"
	"go.aporeto.io/regolithe/schema"
	yaml "gopkg.in/yaml.v2"
)

// An APIInfo holds general information about the API.
type APIInfo struct {
	Prefix  string `yaml:"prefix,omitempty"     json:"prefix,omitempty"`
	Root    string `yaml:"root,omitempty"       json:"root,omitempty"`
	Version int    `yaml:"version,omitempty"    json:"version,omitempty"`
}

// NewAPIInfo returns a new APIInfo.
func NewAPIInfo() *APIInfo {
	return &APIInfo{}
}

// LoadAPIInfo loads an APIInfo from the given file.
func LoadAPIInfo(path string) (*APIInfo, error) {

	apiinfo := NewAPIInfo()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)

	if err = decoder.Decode(apiinfo); err != nil {
		return nil, err
	}

	if err := apiinfo.Validate(); err != nil {
		return nil, formatValidationErrors(err)
	}

	return apiinfo, nil
}

// Validate validates the api info against the schema.
func (a *APIInfo) Validate() []error {

	schemaData, err := schema.Asset("rego-info.json")
	if err != nil {
		return []error{err}
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(a)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return []error{err}
	}

	if !res.Valid() {
		return makeSchemaValidationError("_api.info", res.Errors())
	}

	return nil
}
