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

// An Relation represents a specification Relation.
type Relation struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

	RestName string          `yaml:"rest_name,omitempty"    json:"rest_name,omitempty"`
	Get      *RelationAction `yaml:"get,omitempty"          json:"get,omitempty"`
	Create   *RelationAction `yaml:"create,omitempty"       json:"create,omitempty"`
	Update   *RelationAction `yaml:"update,omitempty"       json:"update,omitempty"`
	Delete   *RelationAction `yaml:"delete,omitempty"       json:"delete,omitempty"`

	currentSpecification Specification
	remoteSpecification  Specification
}

// Specification returns the Specification the API links to.
func (r *Relation) Specification() Specification {
	return r.remoteSpecification
}

// Validate validates the relationship
func (r *Relation) Validate() []error {

	var errs []error

	if r.Get != nil {
		if err := r.Get.Validate(r.currentSpecification.Model().RestName, r.RestName, "get"); err != nil {
			errs = append(errs, err...)
		}
	}

	if r.Create != nil {
		if err := r.Create.Validate(r.currentSpecification.Model().RestName, r.RestName, "create"); err != nil {
			errs = append(errs, err...)
		}
	}

	if r.Update != nil {
		if err := r.Update.Validate(r.currentSpecification.Model().RestName, r.RestName, "update"); err != nil {
			errs = append(errs, err...)
		}
	}

	if r.Delete != nil {
		if err := r.Delete.Validate(r.currentSpecification.Model().RestName, r.RestName, "delete"); err != nil {
			errs = append(errs, err...)
		}
	}

	return errs
}
