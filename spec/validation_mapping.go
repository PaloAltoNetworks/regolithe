package spec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/xeipuuv/gojsonschema"
	"go.aporeto.io/regolithe/schema"
	yaml "gopkg.in/yaml.v2"
)

// A ValidationMap represent a single ValidationMap.
type ValidationMap struct {
	Name   string `yaml:"name,omitempty"           json:"name,omitempty"`
	Import string `yaml:"import,omitempty"         json:"import,omitempty"`
}

// ValidationMapping holds the mapping of the validation function.
type ValidationMapping map[string]map[string]*ValidationMap

// NewValidationMapping returns a new ValidationMapping.
func NewValidationMapping() ValidationMapping {
	return ValidationMapping{}
}

// LoadValidationMapping loads a ValidationMapping from the given ini file.
func LoadValidationMapping(path string) (ValidationMapping, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	vm := NewValidationMapping()

	if err = vm.Read(file, true); err != nil {
		return nil, err
	}

	return vm, nil
}

// Read loads a validation mapping from the given io.Reader
func (v ValidationMapping) Read(reader io.Reader, validate bool) (err error) {

	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)

	if err = decoder.Decode(&v); err != nil {
		return err
	}

	if validate {
		if errs := v.Validate(); len(errs) != 0 {
			return formatValidationErrors(errs)
		}
	}

	return nil
}

// Write dumps the specification into a []byte.
func (v ValidationMapping) Write(writer io.Writer) error {

	repr := yaml.MapSlice{}

	keys := make([]string, len(v))
	var i int
	for k := range v {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, k := range keys {
		repr = append(repr, yaml.MapItem{
			Key:   k,
			Value: v[k],
		})
	}

	data, err := yaml.Marshal(repr)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	lines := bytes.Split(data, []byte("\n"))

	for i, line := range lines {
		condFirstLine := i == 0

		if !condFirstLine && len(line) > 0 && line[0] != ' ' {
			buf.WriteRune('\n')
		}

		buf.Write(line)

		if i+1 < len(lines) {
			buf.WriteRune('\n')
		}
	}

	_, err = writer.Write(buf.Bytes())
	return err
}

// Mapping returns the ValidationMap for the given external type.
func (v ValidationMapping) Mapping(mode string, functionName string) (mapping *ValidationMap, err error) {

	m, ok := v[functionName]
	if !ok {
		return nil, fmt.Errorf("no function '%s' found in type mapping mode %s", functionName, mode)
	}

	return m[mode], nil
}

// Validate validates the type mappings against the schema.
func (v ValidationMapping) Validate() []error {

	schemaData, err := schema.Asset("rego-validation-mapping.json")
	if err != nil {
		return []error{err}
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(v)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return []error{err}
	}

	if !res.Valid() {
		return makeSchemaValidationError("_validation.mapping", res.Errors())
	}

	return nil
}
