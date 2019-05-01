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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRelation_Getters(t *testing.T) {

	Convey("Given I have a new Relation", t, func() {
		rel := &Relation{
			RestName: "test",
			Get:      &RelationAction{Description: "hello get."},
			Create:   &RelationAction{Description: "hello create."},
			Update:   &RelationAction{Description: "hello update."},
			Delete:   &RelationAction{Description: "hello delete."},
			remoteSpecification: &specification{
				RawModel: &Model{
					EntityName: "test2",
				},
			},
		}

		Convey("Then the getters should work", func() {
			So(rel.Specification().Model().EntityName, ShouldEqual, rel.remoteSpecification.Model().EntityName)
		})
	})
}

func TestRelation_Validate(t *testing.T) {

	Convey("Given I have a relation with no validation error", t, func() {

		r := &Relation{
			Get:      &RelationAction{Description: "hello get."},
			Create:   &RelationAction{Description: "hello create."},
			Update:   &RelationAction{Description: "hello update."},
			Delete:   &RelationAction{Description: "hello delete."},
			RestName: "remote",
			currentSpecification: &specification{
				RawModel: &Model{
					EntityName: "current",
				},
			},
			remoteSpecification: &specification{
				RawModel: &Model{
					EntityName: "remote",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := r.Validate()

			Convey("Then there should be no validation error", func() {
				So(len(errs), ShouldEqual, 0)
			})
		})
	})

	Convey("Given I have a relation with all missing description", t, func() {

		r := &Relation{
			Get:      &RelationAction{},
			Create:   &RelationAction{},
			Update:   &RelationAction{},
			Delete:   &RelationAction{},
			RestName: "remote",
			currentSpecification: &specification{
				RawModel: &Model{
					RestName: "current",
				},
			},
			remoteSpecification: &specification{
				RawModel: &Model{
					RestName: "remote",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := r.Validate()

			Convey("Then there should be some validation errors", func() {
				So(len(errs), ShouldEqual, 4)
				So(errs[0].Error(), ShouldEqual, "current.spec: relation 'get' to 'remote' must have a description")
				So(errs[1].Error(), ShouldEqual, "current.spec: relation 'create' to 'remote' must have a description")
				So(errs[2].Error(), ShouldEqual, "current.spec: relation 'update' to 'remote' must have a description")
				So(errs[3].Error(), ShouldEqual, "current.spec: relation 'delete' to 'remote' must have a description")
			})
		})
	})

	Convey("Given I have a relation with all missing final period", t, func() {

		r := &Relation{
			Get:      &RelationAction{Description: "hello get"},
			Create:   &RelationAction{Description: "hello create"},
			Update:   &RelationAction{Description: "hello update"},
			Delete:   &RelationAction{Description: "hello delete"},
			RestName: "remote",
			currentSpecification: &specification{
				RawModel: &Model{
					RestName: "current",
				},
			},
			remoteSpecification: &specification{
				RawModel: &Model{
					RestName: "remote",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := r.Validate()

			Convey("Then there should be some validation errors", func() {
				So(len(errs), ShouldEqual, 4)
				So(errs[0].Error(), ShouldEqual, "current.spec: relation 'get' to 'remote' description must end with a period")
				So(errs[1].Error(), ShouldEqual, "current.spec: relation 'create' to 'remote' description must end with a period")
				So(errs[2].Error(), ShouldEqual, "current.spec: relation 'update' to 'remote' description must end with a period")
				So(errs[3].Error(), ShouldEqual, "current.spec: relation 'delete' to 'remote' description must end with a period")
			})
		})
	})
}
