package spec

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// A Specification represents the a Monolithe Specification.
type Specification struct {
	Attributes []*Attribute `yaml:"attributes,omitempty"`
	APIs       []*API       `yaml:"children,omitempty"`
	Model      *Model       `yaml:"model,omitempty"`

	attributeMap       map[string]*Attribute
	apiMap             map[string]*API
	orderingAttributes []*Attribute
	identifier         *Attribute
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

	if err := spec.Read(file); err != nil {
		return nil, err
	}

	return spec, nil
}

// Read loads a specifaction from the given io.Reader
func (s *Specification) Read(reader io.Reader) error {

	if err := yaml.NewDecoder(reader).Decode(s); err != nil {
		return err
	}

	if err := s.buildAttributesInfo(); err != nil {
		return err
	}

	if err := s.buildAPIsInfo(); err != nil {
		return err
	}

	if s.Model != nil {
		s.Model.EntityNamePlural = Pluralize(s.Model.EntityName)
	}

	return nil
}

// Write dumps the specification into a []byte.
func (s *Specification) Write(writer io.Writer) error {

	s.Attributes = s.OriginalSortedAttributes()

	encoder := yaml.NewEncoder(writer)
	return encoder.Encode(s)
}

// Attribute returns the Attributes with the given name.
func (s *Specification) Attribute(name string) *Attribute {
	return s.attributeMap[name]
}

// OriginalSortedAttributes returns the list of attribute sorted by names.
func (s *Specification) OriginalSortedAttributes() []*Attribute {

	attrs := append([]*Attribute{}, s.Attributes...)

	sort.Slice(attrs, func(i int, j int) bool {
		return strings.Compare(attrs[i].Name, attrs[j].Name) == -1
	})

	return attrs
}

// API returns the API with the given rest name.
func (s *Specification) API(restName string) *API {
	return s.apiMap[restName]
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

		for _, api := range spec.APIs {
			if _, ok := s.apiMap[api.RestName]; !ok {
				s.APIs = append(s.APIs, api)
			}
		}
	}

	if err := s.buildAttributesInfo(); err != nil {
		return err
	}

	if err := s.buildAPIsInfo(); err != nil {
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

// buildAPIsInfo builds the apis map.
func (s *Specification) buildAPIsInfo() error {

	s.apiMap = map[string]*API{}
	for _, api := range s.APIs {

		if _, ok := s.apiMap[api.RestName]; ok {
			return fmt.Errorf("Specification has more than one child api pointing to %s", api.RestName)
		}

		s.apiMap[api.RestName] = api
	}

	return nil
}
