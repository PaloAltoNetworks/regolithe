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

import (
	"fmt"
)

// A Relationship describes the hierarchical relationship of the models.
type Relationship struct {
	Create  map[string]*RelationAction
	Delete  map[string]*RelationAction
	Get     map[string]*RelationAction
	GetMany map[string]*RelationAction
	Update  map[string]*RelationAction
}

// NewRelationship returns a new Relationship.
func NewRelationship() *Relationship {
	return &Relationship{
		Create:  map[string]*RelationAction{},
		Delete:  map[string]*RelationAction{},
		Get:     map[string]*RelationAction{},
		GetMany: map[string]*RelationAction{},
		Update:  map[string]*RelationAction{},
	}
}

// Set sets the names that are allows to do the given action.
func (r *Relationship) Set(action string, name string, ra *RelationAction) {

	var prop map[string]*RelationAction

	switch action {
	case "create":
		prop = r.Create
	case "delete":
		prop = r.Delete
	case "get":
		prop = r.Get
	case "getmany":
		prop = r.GetMany
	case "update":
		prop = r.Update
	default:
		panic(fmt.Sprintf("action '%s' is not valid. Must be 'create', 'delete', 'get', 'getmany' or 'update'", action))
	}

	prop[name] = ra
}
