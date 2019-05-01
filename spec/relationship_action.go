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

import "fmt"

// A RelationAction represents one the the possible action
type RelationAction struct {
	Description         string               `yaml:"description,omitempty"           json:"description,omitempty"`
	Deprecated          bool                 `yaml:"deprecated,omitempty"            json:"deprecated,omitempty"`
	ParameterReferences []string             `yaml:"global_parameters,omitempty"     json:"global_parameters,omitempty"`
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
