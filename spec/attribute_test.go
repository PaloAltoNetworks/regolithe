package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAttribute_Validate(t *testing.T) {

	Convey("Given I have an attribute with not validation error", t, func() {

		a := &Attribute{
			Name:         "name",
			Required:     true,
			ExampleValue: "hello",
			Type:         "string",
		}

		Convey("When I call validate", func() {

			errs := a.Validate()

			Convey("Then there should be no validation error", func() {
				So(len(errs), ShouldEqual, 0)
			})
		})
	})

	Convey("Given I have an attribute with required but no DefaultValue and no ExampleValue", t, func() {

		a := &Attribute{
			Name:     "name",
			Required: true,
			Type:     "string",
			linkedSpecification: &specification{
				RawModel: &Model{
					RestName: "spec",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := a.Validate()

			Convey("Then there should be validation error", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "spec.spec: 'name' is required but has no default_value or example_value")
			})
		})
	})

	Convey("Given I have an attribute with description with period at the end", t, func() {

		a := &Attribute{
			Name:        "name",
			Type:        "string",
			Description: "coucou",
			linkedSpecification: &specification{
				RawModel: &Model{
					RestName: "spec",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := a.Validate()

			Convey("Then there should be validation error", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "spec.spec: description of attribute 'name' must end with a period")
			})
		})
	})

	Convey("Given I have an attribute required with no default value", t, func() {

		a := &Attribute{
			Name:         "name",
			Type:         "enum",
			Description:  "coucou.",
			ExampleValue: "coucou",
			linkedSpecification: &specification{
				RawModel: &Model{
					RestName: "spec",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := a.Validate()

			Convey("Then there should be validation error", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "spec.spec: enum attribute 'name' must define allowed_choices")
			})
		})
	})

	Convey("Given I have an attribute with allowed_chars and no allowed_chars_message", t, func() {

		a := &Attribute{
			Name:         "name",
			Type:         "string",
			Description:  "coucou.",
			ExampleValue: "coucou",
			AllowedChars: "abc",
			linkedSpecification: &specification{
				RawModel: &Model{
					RestName: "spec",
				},
			},
		}

		Convey("When I call validate", func() {

			errs := a.Validate()

			Convey("Then there should be validation error", func() {
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "spec.spec: attribute 'name' must define allowed_chars_message")
			})
		})
	})
}
