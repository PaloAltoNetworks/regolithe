package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeMapping_LoadTypeMapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadTypeMapping("./tests/_type.mapping")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the file should be loaded", func() {
			So(tm, ShouldResemble, TypeMapping{
				"test": map[string]*TypeMap{
					"string_map": &TypeMap{
						Type: "map[string]string",
					},
					"int_array": &TypeMap{
						Type:        "[]int",
						Initializer: "[]int{}",
					},
					"toto": &TypeMap{
						Type:        "Toto",
						Initializer: "Toto{}",
						Import:      "github.com/toto/toto",
					},
				},
			})
		})
	})
}

func TestTypeMapping_Mapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadTypeMapping("./tests/_type.mapping")

		Convey("Then err should bbe nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When I call Mapping on string_map", func() {

			m, err := tm.Mapping("test", "toto")

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the mapping be correct", func() {
				So(m.Type, ShouldEqual, "Toto")
				So(m.Initializer, ShouldEqual, "Toto{}")
				So(m.Import, ShouldEqual, "github.com/toto/toto")
			})
		})
	})
}
