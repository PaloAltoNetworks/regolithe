package spec

// A Model holds generic information about a specification.
type Model struct {
	Aliases          []string `json:"aliases,omitempty"`
	AllowsCreate     bool     `json:"create,omitempty"`
	AllowsDelete     bool     `json:"delete,omitempty"`
	AllowsGet        bool     `json:"get,omitempty"`
	AllowsUpdate     bool     `json:"update,omitempty"`
	Description      string   `json:"description,omitempty"`
	EntityName       string   `json:"entity_name,omitempty"`
	EntityNamePlural string   `json:"-"`
	Extends          []string `json:"extends,omitempty"`
	InstanceName     string   `json:"instance_name,omitempty"`
	IsRoot           bool     `json:"root,omitempty"`
	Package          string   `json:"package,omitempty"`
	ResourceName     string   `json:"resource_name,omitempty"`
	RestName         string   `json:"rest_name,omitempty"`
	Private          bool     `json:"private,omitempty"`
}
