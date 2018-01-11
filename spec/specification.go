package spec

import (
	"encoding/json"
	"fmt"
	"os"
)

type model struct {
	Aliases          []string `json:"aliases"`
	AllowsCreate     bool     `json:"create"`
	AllowsDelete     bool     `json:"delete"`
	AllowsGet        bool     `json:"get"`
	AllowsUpdate     bool     `json:"update"`
	Description      string   `json:"description"`
	EntityName       string   `json:"entity_name"`
	EntityNamePlural string   `json:"-"`
	Exposed          string   `json:"exposed"`
	Extends          []string `json:"extends"`
	InstanceName     string   `json:"instance_name"`
	IsRoot           bool     `json:"is_root"`
	Package          string   `json:"package"`
	ResourceName     string   `json:"resource_name"`
	RestName         string   `json:"rest_name"`
}

// A Specification represents the a Monolithe Specification.
type Specification struct {
	Attributes []*Attribute `json:"attributes"`
	APIs       []*API       `json:"children"`

	attributeMap            map[string]*Attribute
	apiMap                  map[string]*API
	orderingAttributes      []*Attribute
	additionalTypeProviders []string
	identifier              *Attribute

	*model `json:"model,inline"`
}

// NewSpecification returns a new specification.
func NewSpecification() *Specification {
	return &Specification{
		Attributes: []*Attribute{},
		APIs:       []*API{},
		model: &model{
			Extends: []string{},
			Aliases: []string{},
		},
	}
}

// LoadSpecification returns a new specification using the given file path.
func LoadSpecification(path string) (*Specification, error) {

	spec := NewSpecification()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	if err = json.NewDecoder(file).Decode(spec); err != nil {
		return nil, err
	}

	if err = spec.buildAttributesInfo(); err != nil {
		return nil, err
	}

	if err = spec.buildAPIsInfo(); err != nil {
		return nil, err
	}

	spec.EntityNamePlural = Pluralize(spec.EntityName)

	return spec, nil
}

// Attribute returns the Attributes with the given name.
func (s *Specification) Attribute(name string) *Attribute {
	return s.attributeMap[name]
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
			out[attr.Name] = attr.Initializer
			continue
		}

		if attr.DefaultValue != nil {
			if attr.Type == AttributeTypeString || attr.Type == AttributeTypeEnum {
				out[attr.Name] = `"` + attr.DefaultValue.(string) + `"`
				continue
			}
			out[attr.Name] = attr.DefaultValue
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
			return fmt.Errorf("Specification %s has more than one attribute named %s", s.RestName, attr.Name)
		}

		s.attributeMap[attr.Name] = attr

		if attr.Orderable {
			s.orderingAttributes = append(s.orderingAttributes, attr)
		}

		if attr.Identifier {
			if s.identifier != nil {
				return fmt.Errorf("Specification %s has more than one identifier attributes: At least %s and %s", s.RestName, s.identifier.Name, attr.Name)
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
