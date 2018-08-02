package spec

import (
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
	"go.aporeto.io/regolithe/schema"
	"gopkg.in/yaml.v2"
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

	tm := NewValidationMapping()

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)

	if err := decoder.Decode(tm); err != nil {
		return nil, err
	}

	if err := tm.Validate(); err != nil {
		return nil, formatValidationErrors(err)
	}

	return tm, nil
}

// Mapping returns the ValidationMap for the given external type.
func (t ValidationMapping) Mapping(mode string, functionName string) (mapping *ValidationMap, err error) {

	m, ok := t[functionName]
	if !ok {
		return nil, fmt.Errorf("no function '%s' found in type mapping mode %s", functionName, mode)
	}

	return m[mode], nil
}

// Validate validates the type mappings against the schema.
func (t ValidationMapping) Validate() []error {

	schemaData, err := schema.Asset("rego-validation-mapping.json")
	if err != nil {
		return []error{err}
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(t)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return []error{err}
	}

	if !res.Valid() {
		return makeSchemaValidationError("_validation.mapping", res.Errors())
	}

	return nil
}
