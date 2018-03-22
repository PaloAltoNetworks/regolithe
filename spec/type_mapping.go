package spec

import (
	"fmt"
	"os"

	"github.com/aporeto-inc/regolithe/schema"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

// A TypeMap represent a single Type Map.
type TypeMap struct {
	Type        string `yaml:"type,omitempty"           json:"type,omitempty"`
	Initializer string `yaml:"init,omitempty"           json:"init,omitempty"`
	Import      string `yaml:"import,omitempty"         json:"import,omitempty"`
	Description string `yaml:"description,omitempty"    json:"description,omitempty"`
}

// TypeMapping holds the mapping of the external types.
type TypeMapping map[string]map[string]*TypeMap

// NewTypeMapping returns a new TypeMapping.
func NewTypeMapping() TypeMapping {
	return TypeMapping{}
}

// LoadTypeMapping loads a TypeMapping from the given ini file.
func LoadTypeMapping(path string) (TypeMapping, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	tm := NewTypeMapping()

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)

	if err := decoder.Decode(tm); err != nil {
		return nil, err
	}

	if res, err := tm.Validate(); err != nil {
		writeValidationErrors("validation error in _type.mapping", res)
		return nil, err
	}

	return tm, nil
}

// Mapping returns the TypeMap for the given external type.
func (t TypeMapping) Mapping(mode string, externalType string) (mapping *TypeMap, err error) {

	m, ok := t[mode]
	if !ok {
		return nil, fmt.Errorf("no mode '%s' found in type mapping", mode)
	}

	tm, ok := m[externalType]
	if !ok {
		return nil, fmt.Errorf("no type '%s' found in type mapping mode %s", externalType, mode)
	}

	return tm, nil
}

// Validate validates the type mappings against the schema.
func (t TypeMapping) Validate() ([]gojsonschema.ResultError, error) {

	schemaData, err := schema.Asset("rego-type-mapping.json")
	if err != nil {
		return nil, err
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(t)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return nil, err
	}

	if !res.Valid() {
		return res.Errors(), fmt.Errorf("Invalid _type.mapping")
	}

	return nil, nil
}
