package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRelation_Getters(t *testing.T) {

	Convey("Given I have a new Relation", t, func() {
		rel := &Relation{
			RestName:     "test",
			AllowsGet:    true,
			AllowsCreate: true,
			AllowsUpdate: true,
			AllowsDelete: true,
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
			AllowsGet:    true,
			AllowsCreate: true,
			AllowsUpdate: true,
			AllowsDelete: true,
			RestName:     "remote",
			Descriptions: map[string]string{
				"get":    "hello.",
				"create": "hello.",
				"update": "hello.",
				"delete": "hello.",
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
			AllowsGet:    true,
			AllowsCreate: true,
			AllowsUpdate: true,
			AllowsDelete: true,
			RestName:     "remote",
			currentSpecification: &specification{
				RawModel: &Model{
					RestName: "currentSpec",
				},
			},
			remoteSpecification: &specification{
				RawModel: &Model{
					RestName: "remoteSpec",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := r.Validate()

			Convey("Then there should be some validation errors", func() {
				So(len(errs), ShouldEqual, 4)
				So(errs[0].Error(), ShouldEqual, "currentSpec.spec: relation 'get' to 'remote' must have a description")
				So(errs[1].Error(), ShouldEqual, "currentSpec.spec: relation 'create' to 'remote' must have a description")
				So(errs[2].Error(), ShouldEqual, "currentSpec.spec: relation 'update' to 'remote' must have a description")
				So(errs[3].Error(), ShouldEqual, "currentSpec.spec: relation 'delete' to 'remote' must have a description")
			})
		})
	})

	Convey("Given I have a relation with all missing final period", t, func() {

		r := &Relation{
			AllowsGet:    true,
			AllowsCreate: true,
			AllowsUpdate: true,
			AllowsDelete: true,
			RestName:     "remote",
			Descriptions: map[string]string{
				"get":    "hello",
				"create": "hello",
				"update": "hello",
				"delete": "hello",
			},
			currentSpecification: &specification{
				RawModel: &Model{
					RestName: "currentSpec",
				},
			},
			remoteSpecification: &specification{
				RawModel: &Model{
					RestName: "remoteSpec",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := r.Validate()

			Convey("Then there should be some validation errors", func() {
				So(len(errs), ShouldEqual, 4)
				So(errs[0].Error(), ShouldEqual, "currentSpec.spec: relation 'get' to 'remote' description must end with a period")
				So(errs[1].Error(), ShouldEqual, "currentSpec.spec: relation 'create' to 'remote' description must end with a period")
				So(errs[2].Error(), ShouldEqual, "currentSpec.spec: relation 'update' to 'remote' description must end with a period")
				So(errs[3].Error(), ShouldEqual, "currentSpec.spec: relation 'delete' to 'remote' description must end with a period")
			})
		})
	})
}
