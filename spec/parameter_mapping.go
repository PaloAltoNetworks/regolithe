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
	"bytes"
	"io"
	"os"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

// ParameterType represents the various type for a parameter.
type ParameterType string

// Various values for ParameterType.
const (
	ParameterTypeString   ParameterType = "string"
	ParameterTypeInt      ParameterType = "integer"
	ParameterTypeFloat    ParameterType = "float"
	ParameterTypeBool     ParameterType = "boolean"
	ParameterTypeTime     ParameterType = "time"
	ParameterTypeEnum     ParameterType = "enum"
	ParameterTypeDuration ParameterType = "duration"
)

// A ParameterMapping is a list parameter mapping
type ParameterMapping map[string]*ParameterDefinition

// NewParameterMapping returns a new ParameterMapping.
func NewParameterMapping() ParameterMapping {
	return ParameterMapping{}
}

// LoadGlobalParameters loads the global parameters file.
func LoadGlobalParameters(path string) (ParameterMapping, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	pm := ParameterMapping{}

	if err = pm.Read(file, true); err != nil {
		return nil, err
	}

	return pm, nil
}

// Read loads a validation mapping from the given io.Reader
func (p ParameterMapping) Read(reader io.Reader, validate bool) (err error) {

	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)

	if err = decoder.Decode(&p); err != nil {
		return err
	}

	if validate {
		if errs := p.Validate(); len(errs) != 0 {
			return formatValidationErrors(errs)
		}
	}

	return nil
}

// Write dumps the specification into a []byte.
func (p ParameterMapping) Write(writer io.Writer) error {

	repr := yaml.MapSlice{}

	keys := make([]string, len(p))
	var i int
	for k := range p {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, k := range keys {
		repr = append(repr, yaml.MapItem{
			Key:   k,
			Value: p[k],
		})
	}

	data, err := yaml.Marshal(repr)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	prfx1 := []byte("  - name: ")
	prfx2 := []byte("  entries:")
	lines := bytes.Split(data, []byte("\n"))
	var previousLine []byte

	for i, line := range lines {
		condFirstLine := i == 0

		if !condFirstLine &&
			(len(line) > 0 && line[0] != ' ') ||
			(bytes.HasPrefix(line, prfx1) && !bytes.HasPrefix(previousLine, prfx2)) {
			buf.WriteRune('\n')
		}

		buf.Write(line)

		if i+1 < len(lines) {
			buf.WriteRune('\n')
		}

		previousLine = line
	}

	_, err = writer.Write(buf.Bytes())
	return err
}

// Validate the ParameterMapping
func (p ParameterMapping) Validate() []error {

	var errs []error

	for _, v := range p {
		if err := v.Validate("_parameter.mapping"); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}
