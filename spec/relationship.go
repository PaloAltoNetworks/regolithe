package spec

import (
	"fmt"
	"sort"
)

// A Relationship describes the hierarchical relationship of the models.
type Relationship struct {
	AllowsCreate  map[string]struct{}
	AllowsDelete  map[string]struct{}
	AllowsGet     map[string]struct{}
	AllowsGetMany map[string]struct{}
	AllowsUpdate  map[string]struct{}
}

// NewRelationship returns a new Relationship.
func NewRelationship() *Relationship {
	return &Relationship{
		AllowsCreate:  map[string]struct{}{},
		AllowsDelete:  map[string]struct{}{},
		AllowsGet:     map[string]struct{}{},
		AllowsGetMany: map[string]struct{}{},
		AllowsUpdate:  map[string]struct{}{},
	}
}

// Set sets the names that are allows to do the given action.
func (r *Relationship) Set(action string, names ...string) {

	var prop map[string]struct{}

	switch action {
	case "create":
		prop = r.AllowsCreate
	case "delete":
		prop = r.AllowsDelete
	case "get":
		prop = r.AllowsGet
	case "getmany":
		prop = r.AllowsGetMany
	case "update":
		prop = r.AllowsUpdate
	default:
		panic(fmt.Sprintf("action '%s' is not valid. Must be 'create', 'delete', 'get', 'getmany' or 'update'", action))
	}

	for _, n := range names {
		prop[n] = struct{}{}
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
	case "getmany":
		prop = r.AllowsGetMany
	case "update":
		prop = r.AllowsUpdate
	default:
		panic(fmt.Sprintf("action '%s' is not valid. Must be 'create', 'delete', 'get', 'getmany' or 'update'", action))
	}

	for k := range prop {
		names = append(names, k)
	}

	sort.Strings(names)

	return
}
