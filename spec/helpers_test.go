// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	Convey("Given I have no error", t, func() {

		errs := []error{}

		Convey("When I call formatValidationErrors", func() {

			err := formatValidationErrors(errs)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestHelpers_sortVersionString(t *testing.T) {

	Convey("Given I have a version array", t, func() {

		v := []string{"v1", "v4", "v3", "v2"}

		Convey("When I Call sortVersionString", func() {

			sorted := sortVersionStrings(v)

			Convey("Then the version should be sorted", func() {
				So(sorted, ShouldResemble, []string{"v1", "v2", "v3", "v4"})
			})
		})
	})

	Convey("Given I have an invalid version array", t, func() {

		v := []string{"v1", "a", "v3", "v2"}

		Convey("When I Call sortVersionString", func() {

			Convey("Then it should panic", func() {
				So(func() { sortVersionStrings(v) }, ShouldPanicWith, `invalid version 'a'`)
			})
		})
	})
}

func TestHelpers_sortAttributes(t *testing.T) {

	Convey("Given I have a some attributes", t, func() {

		v := []*Attribute{
			&Attribute{
				Name: "c",
			},
			&Attribute{
				Name: "a",
			},
			&Attribute{
				Name: "b",
			},
		}

		Convey("When I Call sortAttributes", func() {

			sortAttributes(v)

			Convey("Then the attributes should be sorted", func() {
				So(len(v), ShouldEqual, 3)
				So(v[0].Name, ShouldEqual, "a")
				So(v[1].Name, ShouldEqual, "b")
				So(v[2].Name, ShouldEqual, "c")
			})
		})
	})
}

func TestHelpers_sortParameters(t *testing.T) {

	Convey("Given I have a some parameters", t, func() {

		p := []*Parameter{
			&Parameter{
				Name: "c",
			},
			&Parameter{
				Name: "a",
			},
			&Parameter{
				Name: "b",
			},
		}

		Convey("When I Call sortParameters", func() {

			sortParameters(p)

			Convey("Then the paremeters should be sorted", func() {
				So(len(p), ShouldEqual, 3)
				So(p[0].Name, ShouldEqual, "a")
				So(p[1].Name, ShouldEqual, "b")
				So(p[2].Name, ShouldEqual, "c")
			})
		})
	})
}
