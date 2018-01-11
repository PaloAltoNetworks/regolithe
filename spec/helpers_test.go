package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHelpers_Pluralize(t *testing.T) {

	Convey("Given I have the word 'Test'", t, func() {

		word := "Test"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Tests")
			})
		})
	})

	Convey("Given I have the word 'Tests'", t, func() {

		word := "Tests"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Tests")
			})
		})
	})

	Convey("Given I have the word 'Testay'", t, func() {

		word := "Testay"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Testays")
			})
		})
	})

	Convey("Given I have the word 'Testey'", t, func() {

		word := "Testey"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Testeys")
			})
		})
	})

	Convey("Given I have the word 'Testiy'", t, func() {

		word := "Testiy"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Testiys")
			})
		})
	})

	Convey("Given I have the word 'Testoy'", t, func() {

		word := "Testoy"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Testoys")
			})
		})
	})

	Convey("Given I have the word 'Testuy'", t, func() {

		word := "Testuy"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Testuys")
			})
		})
	})

	Convey("Given I have the word 'Testy'", t, func() {

		word := "Testy"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "Testies")
			})
		})
	})

	Convey("Given I have the word ''", t, func() {

		word := ""

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "")
			})
		})
	})

	Convey("Given I have the word 'y'", t, func() {

		word := "y"

		Convey("When I call pluralize", func() {

			p := Pluralize(word)

			Convey("Then the word should be pluralized correctly", func() {
				So(p, ShouldEqual, "ys")
			})
		})
	})
}
