package spec

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

type model struct {
	Aliases          []string `json:"aliases,omitempty"`
	AllowsCreate     bool     `json:"create,omitempty"`
	AllowsDelete     bool     `json:"delete,omitempty"`
	AllowsGet        bool     `json:"get,omitempty"`
	AllowsUpdate     bool     `json:"update,omitempty"`
	Description      string   `json:"description,omitempty"`
	EntityName       string   `json:"entity_name,omitempty"`
	EntityNamePlural string   `json:"-"`
	Extends          []string `json:"extends,omitempty"`
	InstanceName     string   `json:"instance_name,omitempty"`
	IsRoot           bool     `json:"root,omitempty"`
	Package          string   `json:"package,omitempty"`
	ResourceName     string   `json:"resource_name,omitempty"`
	RestName         string   `json:"rest_name,omitempty"`
	Private          bool     `json:"private,omitempty"`
}

// A Specification represents the a Monolithe Specification.
type Specification struct {
	Attributes []*Attribute `json:"attributes,omitempty"`
	APIs       []*API       `json:"children,omitempty"`

	attributeMap       map[string]*Attribute
	apiMap             map[string]*API
	orderingAttributes []*Attribute
	identifier         *Attribute

	*model `json:"model,inline,omitempty"`
}

// NewSpecification returns a new specification.
func NewSpecification() *Specification {
	return &Specification{
		model: &model{},
	}
}

// LoadSpecificationFrom loads a specifaction from the given io.Reader
func LoadSpecificationFrom(reader io.Reader) (*Specification, error) {

	spec := NewSpecification()

	if err := json.NewDecoder(reader).Decode(spec); err != nil {
		return nil, err
	}

	if err := spec.buildAttributesInfo(); err != nil {
		return nil, err
	}

	if err := spec.buildAPIsInfo(); err != nil {
		return nil, err
	}

	spec.EntityNamePlural = Pluralize(spec.EntityName)

	return spec, nil
}

// LoadSpecification returns a new specification using the given file path.
func LoadSpecification(specPath string) (*Specification, error) {

	file, err := os.Open(specPath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	spec, err := LoadSpecificationFrom(file)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// Write writes the specification in the given directory.
func (s *Specification) Write(dir string) error {

	s.Attributes = s.OriginalSortedAttributes()

	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(dir, s.RestName+".spec"), append(data, '\n'), 0644)
}

// GetRestName returns the rest name.
func (s *Specification) GetRestName() string {
	return s.RestName
}

// GetEntityName returns the rest name.
func (s *Specification) GetEntityName() string {
	return s.EntityName
}

// GetAllowsGet returns if get is allowed.
func (s *Specification) GetAllowsGet() bool {
	return s.AllowsGet
}

// GetAllowsUpdate returns if update is allowed.
func (s *Specification) GetAllowsUpdate() bool {
	return s.AllowsUpdate
}

// GetAllowsCreate returns if create is allowed.
func (s *Specification) GetAllowsCreate() bool {
	return s.AllowsCreate
}

// GetAllowsDelete returns if delete is allowed.
func (s *Specification) GetAllowsDelete() bool {
	return s.AllowsDelete
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
			return fmt.Errorf("Specification %s has more than one attribute named %s", s.RestName, attr.Name)
		}

		s.attributeMap[attr.Name] = attr

		if attr.DefaultOrder {
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
