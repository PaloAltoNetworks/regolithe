package spec

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// A TypeMap represent a single Type Map.
type TypeMap struct {
	Type        string `yaml:"type,omitempty"`
	Initializer string `yaml:"init,omitempty"`
	Import      string `yaml:"import,omitempty"`
}

// TypeMapping holds the mapping of the external types.
type TypeMapping map[string]map[string]*TypeMap

// NewTypeMapping returns a new TypeMapping.
func NewTypeMapping() TypeMapping {
	return TypeMapping{}
}

// LoadTypeMapping loads a TypeMapping from the given ini file.
func LoadTypeMapping(path string) (TypeMapping, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tm := NewTypeMapping()
	if err := yaml.Unmarshal(data, tm); err != nil {
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
