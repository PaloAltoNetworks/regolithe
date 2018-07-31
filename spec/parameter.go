package spec

import "fmt"

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

// A Parameter represent one parameter that can be
// sent with a query
type Parameter struct {
	Name           string        `yaml:"name,omitempty"              json:"name,omitempty"`
	Description    string        `yaml:"description,omitempty"       json:"description,omitempty"`
	Type           ParameterType `yaml:"type,omitempty"              json:"type,omitempty"`
	AllowedChoices []string      `yaml:"allowed_choices,omitempty"   json:"allowed_choices,omitempty"`
	DefaultValue   interface{}   `yaml:"default_value,omitempty"     json:"default_value,omitempty"`
	ExampleValue   interface{}   `yaml:"example_value,omitempty"     json:"example_value,omitempty"`
	Multiple       bool          `yaml:"multiple,omitempty"          json:"multiple,omitempty"`
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

	if p.DefaultValue == nil && p.ExampleValue == nil {
		errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' must provide an example value as it doesn't have a default", relatedReSTName, p.Name))
	}

	return errs
}
