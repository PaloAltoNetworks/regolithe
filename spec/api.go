package spec

// An API represents a specification API.
type API struct {
	RestName     string `json:"rest_name,omitempty"`
	AllowsGet    bool   `json:"get,omitempty"`
	AllowsCreate bool   `json:"create,omitempty"`
	AllowsUpdate bool   `json:"update,omitempty"`
	AllowsDelete bool   `json:"delete,omitempty"`
	Deprecated   bool   `json:"deprecated,omitempty"`

	linkedSpecification *Specification
}

// Specification returns the Specification the API links to.
func (a *API) Specification() *Specification {
	return a.linkedSpecification
}
