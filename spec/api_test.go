package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI_Getters(t *testing.T) {

	Convey("Given I have a new API", t, func() {
		api := &API{
			RestName:     "test",
			AllowsGet:    true,
			AllowsCreate: true,
			AllowsUpdate: true,
			AllowsDelete: true,
			linkedSpecification: &Specification{
				Model: &Model{
					EntityName: "test2",
				},
			},
		}

		Convey("Then the getters should work", func() {
			So(api.Specification().Model.EntityName, ShouldEqual, api.linkedSpecification.Model.EntityName)
		})
	})
}
