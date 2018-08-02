package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidationMapping_LoadValidationMapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadValidationMapping("./tests/_validation.mapping")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the file should be loaded", func() {
			So(tm, ShouldResemble, ValidationMapping{
				"$username": map[string]*ValidationMap{
					"test": &ValidationMap{
						Name: "validate.CheckUserName",
					},
				},
				"$nospace": map[string]*ValidationMap{
					"test": &ValidationMap{
						Name: "validate.NoSpace",
					},
				},
				"$nocap": map[string]*ValidationMap{
					"test": &ValidationMap{
						Name:   "nocapper.NoCap",
						Import: "github.com/aporeto/nocapper",
					},
					"other": &ValidationMap{
						Name: "noCap",
					},
				},
			})
		})
	})
}

func TestValidationMapping_Mapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadValidationMapping("./tests/_validation.mapping")

		Convey("Then err should bbe nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When I call Mapping on $username for mode test", func() {

			m, err := tm.Mapping("test", "$username")

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the mapping be correct", func() {
				So(m.Name, ShouldEqual, "validate.CheckUserName")
				So(m.Import, ShouldEqual, "")
			})
		})

		Convey("When I call Mapping on $nocap for mode other", func() {

			m, err := tm.Mapping("other", "$nocap")

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the mapping be correct", func() {
				So(m.Name, ShouldEqual, "noCap")
				So(m.Import, ShouldEqual, "")
			})
		})
	})
}
