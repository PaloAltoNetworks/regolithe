package spec

// APIRelationship represents the relationhip of an API.
type APIRelationship string

// Various values for API relationship
const (
	APIRelationshipChild   APIRelationship = "child"
	APIRelationshipMembebr APIRelationship = "member"
	APIRelationshipAlias   APIRelationship = "alias"
	APIRelationshipRoot    APIRelationship = "root"
)

// An API represents a specifcation API.
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
