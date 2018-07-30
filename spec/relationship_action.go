package spec

import "fmt"

// A RelationAction represents one the the possible action
type RelationAction struct {
	Description   string       `yaml:"description,omitempty"  json:"description,omitempty"`
	RawParameters []*Parameter `yaml:"parameters,omitempty"   json:"parameters,omitempty"`
	Deprecated    bool         `yaml:"deprecated,omitempty"     json:"deprecated,omitempty"`
}

// Validate validates the relation action.
func (ra *RelationAction) Validate(currentRestName string, remoteRestName string, k string) error {

	if ra.Description == "" {
		return fmt.Errorf("%s.spec: relation '%s' to '%s' must have a description", currentRestName, k, remoteRestName)
	}

	if ra.Description[len(ra.Description)-1] != '.' {
		return fmt.Errorf("%s.spec: relation '%s' to '%s' description must end with a period", currentRestName, k, remoteRestName)
	}

	var errs []error

	for _, p := range ra.RawParameters {
		if err := p.Validate(currentRestName); err != nil {
			errs = append(errs, err...)
		}
	}

	return nil
}

// Parameters returns the parameters list.
func (ra *RelationAction) Parameters() []*Parameter {
	return ra.RawParameters
}

// Parameter returns the Parameter with the given name.
func (ra *RelationAction) Parameter(name string) *Parameter {

	for _, p := range ra.RawParameters {
		if p.Name == name {
			return p
		}
	}

	return nil
}
