package spec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/xeipuuv/gojsonschema"
	"go.aporeto.io/regolithe/schema"
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

	if err = tm.Read(file, true); err != nil {
		return nil, err
	}

	return tm, nil
}

// Read loads a type mapping from the given io.Reader
func (t TypeMapping) Read(reader io.Reader, validate bool) (err error) {

	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)

	if err = decoder.Decode(&t); err != nil {
		return err
	}

	if validate {
		if errs := t.Validate(); len(errs) != 0 {
			return formatValidationErrors(errs)
		}
	}

	return nil
}

// Write dumps the specification into a []byte.
func (t TypeMapping) Write(writer io.Writer) error {

	repr := yaml.MapSlice{}

	keys := make([]string, len(t))
	var i int
	for k := range t {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, k := range keys {
		repr = append(repr, yaml.MapItem{
			Key:   k,
			Value: t[k],
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

// Mapping returns the TypeMap for the given external type.
func (t TypeMapping) Mapping(mode string, externalType string) (mapping *TypeMap, err error) {

	m, ok := t[externalType]
	if !ok {
		return nil, fmt.Errorf("no type '%s' found in type mapping mode %s", externalType, mode)
	}

	tm, ok := m[mode]
	if !ok {
		return nil, fmt.Errorf("no mode '%s' found in type mapping", mode)
	}

	return tm, nil
}

// Validate validates the type mappings against the schema.
func (t TypeMapping) Validate() []error {

	schemaData, err := schema.Asset("rego-type-mapping.json")
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
		return makeSchemaValidationError("_type.mapping", res.Errors())
	}

	return nil
}
