package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInfo_LoadAPIInfo(t *testing.T) {

	Convey("Given I load an api info", t, func() {

		info, err := LoadAPIInfo("./tests/_api.info")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then apiinfo should be correctly initialized", func() {
			So(info.Prefix, ShouldEqual, "api")
			So(info.Version, ShouldEqual, 1)
			So(info.Root, ShouldEqual, "root")
		})
	})

	Convey("Given I a file that does not exist", t, func() {

		_, err := LoadAPIInfo("./tests/not-api.info")

		Convey("Then err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given I load an bad formatted api info", t, func() {

		_, err := LoadAPIInfo("./tests/broken/_api.info.bad")

		Convey("Then err should be nil", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "yaml: line 4: did not find expected ',' or '}'")
		})
	})
}

func TestInfo_Validate(t *testing.T) {

	Convey("Given I have an api info with no validation error", t, func() {

		a := &APIInfo{
			Prefix:  "/api",
			Root:    "root",
			Version: 1,
		}

		Convey("When I call validate", func() {

			errs := a.Validate()

			Convey("Then there should be no validation error", func() {
				So(len(errs), ShouldEqual, 0)
			})
		})
	})

	Convey("Given I have an api info with validation error", t, func() {

		a := &APIInfo{
			Root:    "root",
			Version: 1,
		}

		Convey("When I call validate", func() {

			errs := a.Validate()

			Convey("Then there should be validation errors", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "_api.info: schema error: prefix: prefix is required")
			})
		})
	})
}
