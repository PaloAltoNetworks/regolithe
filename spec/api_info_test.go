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
			So(info.Version, ShouldEqual, "1")
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
