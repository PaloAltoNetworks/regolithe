package spec

import "fmt"

// A Model holds generic information about a specification.
type Model struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

	RestName     string                 `yaml:"rest_name,omitempty"       json:"rest_name,omitempty"`
	ResourceName string                 `yaml:"resource_name,omitempty"   json:"resource_name,omitempty"`
	EntityName   string                 `yaml:"entity_name,omitempty"     json:"entity_name,omitempty"`
	Package      string                 `yaml:"package,omitempty"         json:"package,omitempty"`
	Group        string                 `yaml:"group,omitempty"           json:"group,omitempty"`
	Description  string                 `yaml:"description,omitempty"     json:"description,omitempty"`
	Aliases      []string               `yaml:"aliases,omitempty"         json:"aliases,omitempty"`
	Private      bool                   `yaml:"private,omitempty"         json:"private,omitempty"`
	Indexes      [][]string             `yaml:"indexes,omitempty"         json:"indexes,omitempty"`
	Get          *RelationAction        `yaml:"get,omitempty"             json:"get,omitempty"`
	Update       *RelationAction        `yaml:"update,omitempty"          json:"update,omitempty"`
	Delete       *RelationAction        `yaml:"delete,omitempty"          json:"delete,omitempty"`
	Extends      []string               `yaml:"extends,omitempty"         json:"extends,omitempty"`
	IsRoot       bool                   `yaml:"root,omitempty"            json:"root,omitempty"`
	Detached     bool                   `yaml:"detached,omitempty"        json:"detached,omitempty"`
	Validations  []string               `yaml:"validations,omitempty"     json:"validations,omitempty"`
	Extensions   map[string]interface{} `yaml:"extensions,omitempty"      json:"extensions,omitempty"`

	EntityNamePlural string `yaml:"-" json:"-"`
}

// Validate validates the Model.
func (m *Model) Validate() []error {

	var errs []error

	if m.Description != "" && m.Description[len(m.Description)-1] != '.' {
		errs = append(errs, fmt.Errorf("%s.spec: model description must end with a period", m.RestName))
	}

	if m.Get != nil {
		if err := m.Get.Validate(m.RestName, m.RestName, "get"); err != nil {
			errs = append(errs, err...)
		}
	}

	if m.Update != nil {
		if err := m.Update.Validate(m.RestName, m.RestName, "update"); err != nil {
			errs = append(errs, err...)
		}
	}

	if m.Delete != nil {
		if err := m.Delete.Validate(m.RestName, m.RestName, "delete"); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}
