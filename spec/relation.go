package spec

import "fmt"

// An Relation represents a specification Relation.
type Relation struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

	RestName     string            `yaml:"rest_name,omitempty"      json:"rest_name,omitempty"`
	Descriptions map[string]string `yaml:"descriptions,omitempty"   json:"descriptions,omitempty"`
	AllowsGet    bool              `yaml:"get,omitempty"            json:"get,omitempty"`
	AllowsCreate bool              `yaml:"create,omitempty"         json:"create,omitempty"`
	AllowsUpdate bool              `yaml:"update,omitempty"         json:"update,omitempty"`
	AllowsDelete bool              `yaml:"delete,omitempty"         json:"delete,omitempty"`
	Deprecated   bool              `yaml:"deprecated,omitempty"     json:"deprecated,omitempty"`

	currentSpecification Specification
	remoteSpecification  Specification
}

// Specification returns the Specification the API links to.
func (a *Relation) Specification() Specification {
	return a.remoteSpecification
}

// Validate validates the relationship
func (a *Relation) Validate() []error {

	var errs []error

	check := func(k string) error {

		d := a.Descriptions[k]

		if d == "" {
			return fmt.Errorf("%s.spec: relation '%s' to '%s' must have a description", a.currentSpecification.Model().RestName, k, a.RestName)
		}

		if d[len(d)-1] != '.' {
			return fmt.Errorf("%s.spec: relation '%s' to '%s' description must end with a period", a.currentSpecification.Model().RestName, k, a.RestName)
		}

		return nil
	}

	if a.AllowsGet {
		if err := check("get"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.AllowsCreate {
		if err := check("create"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.AllowsUpdate {
		if err := check("update"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.AllowsDelete {
		if err := check("delete"); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
