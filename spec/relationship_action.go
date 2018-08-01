package spec

import "fmt"

// A RelationAction represents one the the possible action
type RelationAction struct {
	Description         string               `yaml:"description,omitempty"           json:"description,omitempty"`
	Deprecated          bool                 `yaml:"deprecated,omitempty"            json:"deprecated,omitempty"`
	ParameterReferences []string             `yaml:"globalParameters,omitempty"      json:"globalParameters,omitempty"`
	ParameterDefinition *ParameterDefinition `yaml:"parameters,omitempty"            json:"parameters,omitempty"`
}

// Validate validates the relation action.
func (ra *RelationAction) Validate(currentRestName string, remoteRestName string, k string) []error {

	var errs []error

	if ra.Description == "" {
		errs = append(errs, fmt.Errorf("%s.spec: relation '%s' to '%s' must have a description", currentRestName, k, remoteRestName))
	}

	if ra.Description != "" && ra.Description[len(ra.Description)-1] != '.' {
		errs = append(errs, fmt.Errorf("%s.spec: relation '%s' to '%s' description must end with a period", currentRestName, k, remoteRestName))
	}

	if ra.ParameterDefinition != nil {
		if err := ra.ParameterDefinition.Validate(currentRestName); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}
