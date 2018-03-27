package spec

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/xeipuuv/gojsonschema"
)

type mockResultError struct{}

func (r *mockResultError) Field() string                        { return "field" }
func (r *mockResultError) SetType(string)                       {}
func (r *mockResultError) Type() string                         { return "type" }
func (r *mockResultError) SetContext(*gojsonschema.JsonContext) {}
func (r *mockResultError) Context() *gojsonschema.JsonContext   { return nil }
func (r *mockResultError) SetDescription(string)                {}
func (r *mockResultError) Description() string                  { return "description" }
func (r *mockResultError) SetDescriptionFormat(string)          {}
func (r *mockResultError) DescriptionFormat() string            { return "format" }
func (r *mockResultError) SetValue(interface{})                 {}
func (r *mockResultError) Value() interface{}                   { return nil }
func (r *mockResultError) SetDetails(gojsonschema.ErrorDetails) {}
func (r *mockResultError) Details() gojsonschema.ErrorDetails   { return gojsonschema.ErrorDetails{} }
func (r *mockResultError) String() string                       { return "woops" }

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

func TestHelpers_makeSchemaValidationErrors(t *testing.T) {

	Convey("Given I have some validation error", t, func() {

		errs := []gojsonschema.ResultError{
			&mockResultError{},
		}

		Convey("When I call makeSchemaValidationErrors", func() {

			errs := makeSchemaValidationError("hello.spec", errs)

			Convey("Then err should be correct", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "hello.spec: schema error: woops")
			})
		})
	})
}

func TestHelpers_formatValidationErrors(t *testing.T) {

	Convey("Given I have some errors", t, func() {

		errs := []error{
			fmt.Errorf("1-err1"),
			fmt.Errorf("0-err2"),
			fmt.Errorf("2-err2"),
		}

		Convey("When I call formatValidationErrors", func() {

			err := formatValidationErrors(errs)

			Convey("Then err should be correct", func() {
				So(err.Error(), ShouldEqual, "0-err2\n1-err1\n2-err2")
			})
		})
	})
}
