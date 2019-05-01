// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
)

// ParameterDefinition represents a parameter definition.
type ParameterDefinition struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

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
	} else {

		if p.Type != ParameterTypeString &&
			p.Type != ParameterTypeInt &&
			p.Type != ParameterTypeFloat &&
			p.Type != ParameterTypeBool &&
			p.Type != ParameterTypeTime &&
			p.Type != ParameterTypeDuration &&
			p.Type != ParameterTypeEnum {
			errs = append(errs, fmt.Errorf("%s.spec: type of parameter '%s' must be 'string', 'integer', 'float', 'boolean', 'enum', 'time' or 'duration'", relatedReSTName, p.Name))
		}
	}

	if p.Type == ParameterTypeEnum && len(p.AllowedChoices) == 0 {
		errs = append(errs, fmt.Errorf("%s.spec: enum parameter '%s' must define allowed_choices", relatedReSTName, p.Name))
	}

	if p.Type != ParameterTypeEnum && len(p.AllowedChoices) > 0 {
		errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is not an enum but defines allowed_choices", relatedReSTName, p.Name))
	}

	if p.DefaultValue == nil && p.ExampleValue == nil && p.Type == ParameterTypeString {
		errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' must provide an example value as it doesn't have a default", relatedReSTName, p.Name))
	}

	if p.DefaultValue != nil {
		switch p.Type {
		case ParameterTypeString:
			if _, ok := p.DefaultValue.(string); !ok {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an string, but the default value is not", relatedReSTName, p.Name))
			}
		case ParameterTypeEnum:
			if _, ok := p.DefaultValue.(string); !ok {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an enum, but the default value is not", relatedReSTName, p.Name))
			}
		case ParameterTypeInt:
			if _, ok := p.DefaultValue.(int); !ok {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an integer, but the default value is not", relatedReSTName, p.Name))
			}
		case ParameterTypeFloat:
			if _, ok := p.DefaultValue.(float64); !ok {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an float, but the default value is not", relatedReSTName, p.Name))
			}
		case ParameterTypeBool:
			if _, ok := p.DefaultValue.(bool); !ok {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an boolean, but the default value is not", relatedReSTName, p.Name))
			}
		case ParameterTypeDuration:
			if _, ok := p.DefaultValue.(string); !ok {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an duration, but the default value is not", relatedReSTName, p.Name))
				break
			}
			if _, err := time.ParseDuration(p.DefaultValue.(string)); err != nil {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an duration, but the default value is not", relatedReSTName, p.Name))
			}
		case ParameterTypeTime:
			if _, ok := p.DefaultValue.(string); !ok {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an time, but the default value is not", relatedReSTName, p.Name))
				break
			}
			if _, err := dateparse.ParseAny(p.DefaultValue.(string)); err != nil {
				errs = append(errs, fmt.Errorf("%s.spec: parameter '%s' is defined as an time, but the default value is not", relatedReSTName, p.Name))
			}
		}
	}

	return errs
}
