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
			remoteSpecification: &Specification{
				Model: &Model{
					EntityName: "test2",
				},
			},
		}

		Convey("Then the getters should work", func() {
			So(rel.Specification().Model.EntityName, ShouldEqual, rel.remoteSpecification.Model.EntityName)
		})
	})
}
