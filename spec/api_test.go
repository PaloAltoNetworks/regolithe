package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI_Getters(t *testing.T) {

	Convey("Given I have a new API", t, func() {
		api := &API{
			RestName:         "test",
			AllowsGet:        true,
			AllowsCreate:     true,
			AllowsUpdate:     true,
			AllowsDelete:     true,
			AllowsBulkCreate: true,
			AllowsBulkUpdate: true,
			AllowsBulkDelete: true,
			linkedSpecification: &Specification{
				model: &model{
					EntityName: "test2",
				},
			},
		}

		Convey("Then the getters should work", func() {
			So(api.Specification().EntityName, ShouldEqual, api.GetEntityName())
			So(api.GetRestName(), ShouldEqual, "test")
			So(api.GetAllowsGet(), ShouldBeTrue)
			So(api.GetAllowsUpdate(), ShouldBeTrue)
			So(api.GetAllowsCreate(), ShouldBeTrue)
			So(api.GetAllowsDelete(), ShouldBeTrue)
		})
	})
}
