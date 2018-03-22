package spec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/aporeto-inc/regolithe/schema"
	"github.com/xeipuuv/gojsonschema"

	wordwrap "github.com/mitchellh/go-wordwrap"
	yaml "gopkg.in/yaml.v2"
)

const (
	rootModelKey      = "model"
	rootAttributesKey = "attributes"
	rootRelationsKey  = "relations"
)

// A Specification represents the a Monolithe Specification.
type Specification struct {
	Attributes []*Attribute `yaml:"attributes,omitempty"    json:"attributes,omitempty"`
	Relations  []*Relation  `yaml:"relations,omitempty"     json:"relations,omitempty"`
	Model      *Model       `yaml:"model,omitempty"         json:"model,omitempty"`

	attributeMap       map[string]*Attribute
	relationsMap       map[string]*Relation
	orderingAttributes []*Attribute
	identifier         *Attribute
	path               string
}

// NewSpecification returns a new specification.
func NewSpecification() *Specification {
	return &Specification{}
}

// LoadSpecification returns a new specification using the given file path.
func LoadSpecification(specPath string) (*Specification, error) {

	file, err := os.Open(specPath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	spec := NewSpecification()
	spec.path = specPath

	if err = spec.Read(file); err != nil {
		return nil, err
	}

	return spec, nil
}

// Read loads a specifaction from the given io.Reader
func (s *Specification) Read(reader io.Reader) error {

	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)

	if err := decoder.Decode(s); err != nil {
		return err
	}

	if err := s.buildAttributesInfo(); err != nil {
		return err
	}

	if err := s.buildRelationssInfo(); err != nil {
		return err
	}

	if s.Model != nil {
		s.Model.EntityNamePlural = Pluralize(s.Model.EntityName)
	}

	if res, err := s.Validate(); err != nil {
		writeValidationErrors(fmt.Sprintf("validation error in %s", s.path), res)
		return err
	}

	return nil
}

// Write dumps the specification into a []byte.
func (s *Specification) Write(writer io.Writer) error {

	repr := yaml.MapSlice{}

	if s.Model != nil {
		repr = append(repr, yaml.MapItem{Key: rootModelKey, Value: toYAMLMapSlice(s.Model)})
	}

	if len(s.Attributes) != 0 {

		attrs := make([]yaml.MapSlice, len(s.Attributes))

		for i, attr := range s.SortedAttributes() {
			attrs[i] = toYAMLMapSlice(attr)
		}

		repr = append(repr, yaml.MapItem{Key: rootAttributesKey, Value: attrs})
	}

	if len(s.Relations) != 0 {

		relations := make([]yaml.MapSlice, len(s.Relations))

		for i, rel := range s.Relations {
			for k, v := range rel.Descriptions {
				rel.Descriptions[k] = wordwrap.WrapString(v, 80)
			}
			relations[i] = toYAMLMapSlice(rel)
		}

		repr = append(repr, yaml.MapItem{Key: rootRelationsKey, Value: relations})
	}

	data, err := yaml.Marshal(repr)
	if err != nil {
		return err
	}

	var previousLine []byte
	buf := &bytes.Buffer{}
	prfx1 := []byte("- ")
	yamlModelKey := []byte(rootModelKey + ":")
	yamlAttrKey := []byte(rootAttributesKey + ":")
	yamlAttrRelation := []byte(rootRelationsKey + ":")

	lines := bytes.Split(data, []byte("\n"))
	lineN := len(lines)

	for i, line := range lines {

		condFirstLine := i == 0
		condFirstIn := bytes.Equal(previousLine, yamlAttrKey) || bytes.Equal(previousLine, yamlAttrRelation)
		condPrefixed := bytes.HasPrefix(line, prfx1)

		if !condFirstLine && !condFirstIn && condPrefixed {
			buf.WriteRune('\n')
		}

		if bytes.Equal(line, yamlModelKey) {
			if !condFirstLine {
				buf.WriteRune('\n')
			}
			buf.WriteString("# Model\n")
		}
		if bytes.Equal(line, yamlAttrKey) {
			if !condFirstLine {
				buf.WriteRune('\n')
			}
			buf.WriteString("# Attributes\n")
		}
		if bytes.Equal(line, yamlAttrRelation) {
			if !condFirstLine {
				buf.WriteRune('\n')
			}
			buf.WriteString("# Relations\n")
		}

		buf.Write(line)
		if i+1 < lineN {
			buf.WriteRune('\n')
		}

		previousLine = line
	}

	_, err = writer.Write(buf.Bytes())
	return err
}

// Validate validates the spec agains the schema.
func (s *Specification) Validate() ([]gojsonschema.ResultError, error) {

	var schemaData []byte
	var err error

	if s.Model == nil {
		schemaData, err = schema.Asset("rego-abstract.json")
	} else {
		schemaData, err = schema.Asset("rego-spec.json")
	}

	if err != nil {
		return nil, err
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(s)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return nil, err
	}

	if !res.Valid() {
		return res.Errors(), fmt.Errorf("Invalid specification")
	}

	return nil, nil
}

// Attribute returns the Attributes with the given name.
func (s *Specification) Attribute(name string) *Attribute {
	return s.attributeMap[name]
}

// SortedAttributes returns the list of attribute sorted by names.
func (s *Specification) SortedAttributes() []*Attribute {

	attrs := append([]*Attribute{}, s.Attributes...)

	sort.Slice(attrs, func(i int, j int) bool {
		return strings.Compare(attrs[i].Name, attrs[j].Name) == -1
	})

	return attrs
}

// ExposedAttributes returns the exposed attributes.
func (s *Specification) ExposedAttributes() []*Attribute {

	var attrs []*Attribute

	for _, attr := range s.SortedAttributes() {
		if !attr.Exposed {
			continue
		}

		attrs = append(attrs, attr)
	}

	return attrs
}

// Relation returns the Relation with the given rest name.
func (s *Specification) Relation(restName string) *Relation {
	return s.relationsMap[restName]
}

// Identifier returns all the identifier attribute.
func (s *Specification) Identifier() *Attribute {
	return s.identifier
}

// OrderingAttributes returns all the ordering attribute.
func (s *Specification) OrderingAttributes() []*Attribute {
	return append([]*Attribute{}, s.orderingAttributes...)
}

// ApplyBaseSpecifications applies attributes of the given *Specifications to the receiver.
func (s *Specification) ApplyBaseSpecifications(specs ...*Specification) error {

	for _, spec := range specs {

		for _, attr := range spec.Attributes {
			if _, ok := s.attributeMap[attr.Name]; !ok {
				s.Attributes = append(s.Attributes, attr)
			}
		}

		for _, rel := range spec.Relations {
			if _, ok := s.relationsMap[rel.RestName]; !ok {
				s.Relations = append(s.Relations, rel)
			}
		}
	}

	if err := s.buildAttributesInfo(); err != nil {
		return err
	}

	if err := s.buildRelationssInfo(); err != nil {
		return err
	}

	return nil
}

// TypeProviders returns the unique list of all attributes type providers.
func (s *Specification) TypeProviders() []string {

	yes := &struct{}{}
	cache := map[string]*struct{}{}
	var providers []string

	for _, attr := range s.Attributes {
		if attr.TypeProvider == "" {
			continue
		}

		if _, ok := cache[attr.TypeProvider]; ok {
			continue
		}

		cache[attr.TypeProvider] = yes
		providers = append(providers, attr.TypeProvider)
	}

	return providers
}

// AttributeInitializers returns all initializers of the Attributes.
func (s *Specification) AttributeInitializers() map[string]interface{} {

	out := map[string]interface{}{}

	for _, attr := range s.Attributes {

		if attr.Initializer != "" {
			out[attr.ConvertedName] = attr.Initializer
			continue
		}

		if attr.DefaultValue != nil {
			if attr.Type == AttributeTypeString || attr.Type == AttributeTypeEnum {
				out[attr.ConvertedName] = `"` + attr.DefaultValue.(string) + `"`
				continue
			}
			out[attr.ConvertedName] = attr.DefaultValue
		}
	}

	return out
}

// buildAttributesInfo builds the attributes map.
func (s *Specification) buildAttributesInfo() error {

	s.attributeMap = map[string]*Attribute{}
	s.orderingAttributes = []*Attribute{}
	s.identifier = nil

	for _, attr := range s.Attributes {

		if _, ok := s.attributeMap[attr.Name]; ok {
			if s.Model != nil {
				return fmt.Errorf("Specification %s has more than one attribute named %s", s.Model.RestName, attr.Name)
			}

			return fmt.Errorf("One abstract has more than one attribute named %s", attr.Name)
		}

		s.attributeMap[attr.Name] = attr

		if attr.DefaultOrder {
			s.orderingAttributes = append(s.orderingAttributes, attr)
		}

		if attr.Identifier {
			if s.identifier != nil {
				return fmt.Errorf("Specification %s has more than one identifier attributes: At least %s and %s", s.Model.RestName, s.identifier.Name, attr.Name)
			}
			s.identifier = attr
		}
	}

	return nil
}

// buildRelationssInfo builds the relations map.
func (s *Specification) buildRelationssInfo() error {

	s.relationsMap = map[string]*Relation{}
	for _, rel := range s.Relations {

		if _, ok := s.relationsMap[rel.RestName]; ok {
			return fmt.Errorf("Specification has more than one child relation pointing to %s", rel.RestName)
		}

		s.relationsMap[rel.RestName] = rel
	}

	return nil
}
