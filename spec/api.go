package spec

// An API represents a specification API.
type API struct {
	RestName     string `yaml:"rest_name,omitempty"`
	AllowsGet    bool   `yaml:"get,omitempty"`
	AllowsCreate bool   `yaml:"create,omitempty"`
	AllowsUpdate bool   `yaml:"update,omitempty"`
	AllowsDelete bool   `yaml:"delete,omitempty"`
	Deprecated   bool   `yaml:"deprecated,omitempty"`

	linkedSpecification *Specification
}

// Specification returns the Specification the API links to.
func (a *API) Specification() *Specification {
	return a.linkedSpecification
}
