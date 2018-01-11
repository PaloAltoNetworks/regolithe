package spec

import (
	"errors"

	ini "gopkg.in/ini.v1"
)

// A TypeMap represent a single Type Map.
type TypeMap struct {
	Type        string
	Initializer string
	Import      string
}

// TypeMapping holds the mapping of the external types.
type TypeMapping struct {
	data *ini.File
}

// NewTypeMapping returns a new TypeMapping.
func NewTypeMapping() *TypeMapping {
	return &TypeMapping{}
}

// LoadTypeMapping loads a TypeMapping from the given ini file.
func LoadTypeMapping(path string) (*TypeMapping, error) {

	tm := NewTypeMapping()

	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, path)
	if err != nil {
		return nil, err
	}

	tm.data = cfg

	return tm, nil
}

// Mapping returns the TypeMap for the given external type.
func (t *TypeMapping) Mapping(mode string, externalType string) (mapping TypeMap, err error) {

	section, err := t.data.GetSection(mode)
	if err != nil {
		return
	}

	key, err := section.GetKey(externalType)
	if err != nil {
		return
	}

	parts := key.Strings(";")
	if len(parts) == 0 {
		err = errors.New("Invalid type mapping")
		return
	}

	if len(parts) >= 1 {
		mapping.Type = parts[0]
	}
	if len(parts) >= 2 {
		mapping.Initializer = parts[1]
	}
	if len(parts) == 3 {
		mapping.Import = parts[2]
	}

	return
}
