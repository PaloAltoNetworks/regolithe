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
