package spec

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ParameterType represents the various type for a parameter.
type ParameterType string

// Various values for ParameterType.
const (
	ParameterTypeString   ParameterType = "string"
	ParameterTypeInt      ParameterType = "integer"
	ParameterTypeFloat    ParameterType = "float"
	ParameterTypeBool     ParameterType = "boolean"
	ParameterTypeTime     ParameterType = "time"
	ParameterTypeEnum     ParameterType = "enum"
	ParameterTypeDuration ParameterType = "duration"
)

// LoadGlobalParameters loads the global parameters file.
func LoadGlobalParameters(path string) (map[string]*ParameterDefinition, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	pm := map[string]*ParameterDefinition{}

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)

	if err := decoder.Decode(pm); err != nil {
		return nil, err
	}

	for _, v := range pm {
		if err := v.Validate("_parameters"); err != nil {
			return nil, formatValidationErrors(err)
		}
	}

	return pm, nil
}

// ParameterDefinition represents a parameter definition.
type ParameterDefinition struct {
	Required [][][]string `yaml:"required,omitempty"    json:"required,omitempty"`
	Entries  []*Parameter `yaml:"entries,omitempty"     json:"entries,omitempty"`
}

func (p *ParameterDefinition) extend(additional *ParameterDefinition, key string) error {

	if additional == nil {
		return fmt.Errorf("unable to find global parameter key '%s'", key)
	}

	p.Required = append(p.Required, additional.Required...)
	p.Entries = append(p.Entries, additional.Entries...)

	return nil
}

// Validate validates the parameter definition.
func (p *ParameterDefinition) Validate(relatedReSTName string) []error {

	var errs []error
	for _, p := range p.Entries {
		if err := p.Validate(relatedReSTName); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}

// A Parameter represent one parameter that can be
// sent with a query
type Parameter struct {
	Name           string        `yaml:"name,omitempty"              json:"name,omitempty"`
	Description    string        `yaml:"description,omitempty"       json:"description,omitempty"`
	Type           ParameterType `yaml:"type,omitempty"              json:"type,omitempty"`
	Multiple       bool          `yaml:"multiple,omitempty"          json:"multiple,omitempty"`
	AllowedChoices []string      `yaml:"allowed_choices,omitempty"   json:"allowed_choices,omitempty"`
	DefaultValue   interface{}   `yaml:"default_value,omitempty"     json:"default_value,omitempty"`
	ExampleValue   interface{}   `yaml:"example_value,omitempty"     json:"example_value,omitempty"`
}

// Validate validates the parameter definition.
func (p *Parameter) Validate(relatedReSTName string) []error {

	var errs []error

	if p.Description == "" || p.Description[len(p.Description)-1] != '.' {
		errs = append(errs, fmt.Errorf("%s.spec: description of parameter '%s' must end with a period", relatedReSTName, p.Name))
	}

	if p.Type == "" {
		errs = append(errs, fmt.Errorf("%s.spec: type of parameter '%s' must be set", relatedReSTName, p.Name))
	}

	if p.Type != ParameterTypeString &&
		p.Type != ParameterTypeInt &&
		p.Type != ParameterTypeFloat &&
		p.Type != ParameterTypeBool &&
		p.Type != ParameterTypeTime &&
		p.Type != ParameterTypeDuration &&
		p.Type != ParameterTypeEnum {
		errs = append(errs, fmt.Errorf("%s.spec: type of parameter '%s' must be 'string', 'integer', 'float', 'boolean', 'enum', 'time' or 'duration'", relatedReSTName, p.Name))
	}

	if p.Type == ParameterTypeEnum && len(p.AllowedChoices) == 0 {
		errs = append(errs, fmt.Errorf("%s.spec: enum parameter '%s' must define allowed_choices", relatedReSTName, p.Name))
	}

	if p.DefaultValue == nil && p.ExampleValue == nil && p.Type == ParameterTypeString {
		errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' must provide an example value as it doesn't have a default", relatedReSTName, p.Name))
	}

	return errs
}
