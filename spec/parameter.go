package spec

import "fmt"

// ParameterType represents the various type for a parameter.
type ParameterType string

// Various values for ParameterType.
const (
	ParameterTypeString ParameterType = "string"
	ParameterTypeInt    ParameterType = "integer"
	ParameterTypeFloat  ParameterType = "float"
	ParameterTypeBool   ParameterType = "boolean"
	ParameterTypeObject ParameterType = "object"
	ParameterTypeTime   ParameterType = "time"
)

// A Parameter represent one parameter that can be
// sent with a query
type Parameter struct {
	Name        string        `yaml:"name,omitempty"            json:"name,omitempty"`
	Description string        `yaml:"description,omitempty"     json:"description,omitempty"`
	Type        ParameterType `yaml:"type,omitempty"            json:"type,omitempty"`
	Mutliple    bool          `yaml:"multiple,omitempty"        json:"multiple,omitempty"`
}

// Validate validates the parameter definition.
func (p *Parameter) Validate(relatedReSTName string) []error {

	var errs []error

	if p.Description == "" || p.Description[len(p.Description)-1] != '.' {
		errs = append(errs, fmt.Errorf("%s.spec: description of parameter '%s' must end with a period", relatedReSTName, p.Name))
	}

	return errs
}
