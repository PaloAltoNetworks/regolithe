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
				"string_map": map[string]*TypeMap{
					"test": &TypeMap{
						Type: "map[string]string",
					},
				},
				"int_array": map[string]*TypeMap{
					"test": &TypeMap{
						Type:        "[]int",
						Initializer: "[]int{}",
					},
				},
				"toto": map[string]*TypeMap{
					"test": &TypeMap{
						Type:        "Toto",
						Initializer: "Toto{}",
						Import:      "github.com/toto/toto",
					},
					"other": &TypeMap{
						Type:        "Object",
						Initializer: "new Object()",
					},
				},
			})
		})
	})
}

func TestTypeMapping_Mapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadTypeMapping("./tests/_type.mapping")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When I call Mapping on string_map for mode test", func() {

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

		Convey("When I call Mapping on string_map for mode other", func() {

			m, err := tm.Mapping("other", "toto")

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the mapping be correct", func() {
				So(m.Type, ShouldEqual, "Object")
				So(m.Initializer, ShouldEqual, "new Object()")
				So(m.Import, ShouldEqual, "")
			})
		})
	})
}

func TestTypeMapping_All(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadTypeMapping("./tests/_type.mapping")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When I call All", func() {

			m := tm.All("test")

			Convey("Then the mapping be correct", func() {
				So(len(m), ShouldEqual, 3)
				So(m[0].Type, ShouldEqual, "[]int")
				So(m[1].Type, ShouldEqual, "map[string]string")
				So(m[2].Type, ShouldEqual, "Toto")
			})
		})
	})
}
