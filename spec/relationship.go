package spec

import (
	"fmt"
	"sort"
)

// A RelationshipHolder implements some common method about relationship.
type RelationshipHolder interface {
	GetRestName() string
	GetEntityName() string
	GetAllowsGet() bool
	GetAllowsUpdate() bool
	GetAllowsCreate() bool
	GetAllowsDelete() bool
}

// A Relationship describes the hierarchical relationship of the models.
type Relationship struct {
	AllowsCreate map[string]struct{}
	AllowsDelete map[string]struct{}
	AllowsGet    map[string]struct{}
	AllowsUpdate map[string]struct{}
	Mode         APIRelationship
}

// NewRelationship returns a new Relationship.
func NewRelationship(mode APIRelationship) *Relationship {
	return &Relationship{
		AllowsCreate: map[string]struct{}{},
		AllowsDelete: map[string]struct{}{},
		AllowsGet:    map[string]struct{}{},
		AllowsUpdate: map[string]struct{}{},
		Mode:         mode,
	}
}

// Get returns the sorted list of rest name for the given action.
func (r *Relationship) Get(action string) (names []string) {

	var prop map[string]struct{}

	switch action {
	case "create":
		prop = r.AllowsCreate
	case "delete":
		prop = r.AllowsDelete
	case "get":
		prop = r.AllowsGet
	case "update":
		prop = r.AllowsUpdate
	default:
		panic(fmt.Sprintf("action '%s' is not valid. Must be 'create', 'delete', 'get' or 'update'", action))
	}

	for k := range prop {
		names = append(names, k)
	}

	sort.Strings(names)

	return
}
