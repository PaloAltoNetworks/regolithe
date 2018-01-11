package spec

import (
	"encoding/json"
	"fmt"
	"os"
)

type model struct {
	Aliases      []string `json:"aliases"`
	AllowsCreate bool     `json:"create"`
	AllowsDelete bool     `json:"delete"`
	AllowsGet    bool     `json:"get"`
	AllowsUpdate bool     `json:"update"`
	Description  string   `json:"description"`
	EntityName   string   `json:"entity_name"`
	Exposed      string   `json:"exposed"`
	Extends      []string `json:"extends"`
	InstanceName string   `json:"instance_name"`
	IsRoot       bool     `json:"is_root"`
	Package      string   `json:"package"`
	ResourceName string   `json:"resource_name"`
	RestName     string   `json:"rest_name"`
}

// A Specification represents the a Monolithe Specification.
type Specification struct {
	Attributes []*Attribute `json:"attributes"`
	APIs       []*API       `json:"children"`

	attributeMap map[string]*Attribute
	apiMap       map[string]*API

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

	if err = spec.BuildAttributeNames(); err != nil {
		return nil, err
	}

	if err = spec.BuildAPINames(); err != nil {
		return nil, err
	}

	return spec, nil
}

// BuildAttributeNames builds the attributes map.
func (s *Specification) BuildAttributeNames() error {

	s.attributeMap = map[string]*Attribute{}

	for _, attr := range s.Attributes {

		if _, ok := s.attributeMap[attr.Name]; ok {
			return fmt.Errorf("Specification has more than one attribute named %s", attr.Name)
		}

		s.attributeMap[attr.Name] = attr
	}

	return nil
}

// BuildAPINames builds the apis map.
func (s *Specification) BuildAPINames() error {

	s.apiMap = map[string]*API{}
	for _, api := range s.APIs {

		if _, ok := s.apiMap[api.RestName]; ok {
			return fmt.Errorf("Specification has more than one child api pointing to %s", api.RestName)
		}

		s.apiMap[api.RestName] = api
	}

	return nil
}

// Attribute returns the Attributes with the given name.
func (s *Specification) Attribute(name string) *Attribute {
	return s.attributeMap[name]
}

// API returns the API with the given rest name.
func (s *Specification) API(restName string) *API {
	return s.apiMap[restName]
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

	if err := s.BuildAttributeNames(); err != nil {
		return err
	}

	if err := s.BuildAPINames(); err != nil {
		return err
	}

	return nil
}
