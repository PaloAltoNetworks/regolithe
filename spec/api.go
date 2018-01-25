package spec

// APIRelationship represents the relationhip of an API.
type APIRelationship string

// Various values for API relationship
const (
	APIRelationshipChild APIRelationship = "child"
	APIRelationshipRoot  APIRelationship = "root"
)

// An API represents a specification API.
type API struct {
	RestName         string          `json:"rest_name"`
	AllowsGet        bool            `json:"get"`
	AllowsCreate     bool            `json:"create"`
	AllowsUpdate     bool            `json:"update"`
	AllowsDelete     bool            `json:"delete"`
	AllowsBulkCreate bool            `json:"bulk_create"`
	AllowsBulkUpdate bool            `json:"bulk+update"`
	AllowsBulkDelete bool            `json:"bulk_delete"`
	Deprecated       bool            `json:"deprecated"`
	Relationship     APIRelationship `json:"relationship"`

	linkedSpecification *Specification
}

// Specification returns the Specification the API links to.
func (a *API) Specification() *Specification {
	return a.linkedSpecification
}

// GetRestName returns the rest name.
func (a *API) GetRestName() string {
	return a.RestName
}

// GetEntityName returns the rest name.
func (a *API) GetEntityName() string {
	return a.linkedSpecification.EntityName
}

// GetAllowsGet returns if get is allowed.
func (a *API) GetAllowsGet() bool {
	return a.AllowsGet
}

// GetAllowsUpdate returns if update is allowed.
func (a *API) GetAllowsUpdate() bool {
	return a.AllowsUpdate
}

// GetAllowsCreate returns if create is allowed.
func (a *API) GetAllowsCreate() bool {
	return a.AllowsCreate
}

// GetAllowsDelete returns if delete is allowed.
func (a *API) GetAllowsDelete() bool {
	return a.AllowsDelete
}
