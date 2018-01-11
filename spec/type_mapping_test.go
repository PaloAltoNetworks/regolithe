package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeMapping_LoadTypeMapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadTypeMapping("./tests/type_mapping.ini")

		Convey("Then err should bbe nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the file should be loaded", func() {
			So(tm.data, ShouldNotBeNil)
		})
	})
}

func TestTypeMapping_Mapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadTypeMapping("./tests/type_mapping.ini")

		Convey("Then err should bbe nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When I call Mapping on string_map", func() {

			m, err := tm.Mapping("elemental", "string_map")

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the mapping be correct", func() {
				So(m.Type, ShouldEqual, "map[string]string")
				So(m.Initializer, ShouldEqual, "make(map[string]string)")
				So(m.Import, ShouldBeEmpty)
			})
		})
	})
}
