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
					"test": {
						Type: "map[string]string",
					},
				},
				"int_array": map[string]*TypeMap{
					"test": {
						Type:        "[]int",
						Initializer: "[]int{}",
					},
				},
				"toto": map[string]*TypeMap{
					"test": {
						Type:        "Toto",
						Initializer: "Toto{}",
						Import:      "github.com/toto/toto",
					},
					"other": {
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
				So(m[0].Type, ShouldEqual, "Toto")
				So(m[1].Type, ShouldEqual, "[]int")
				So(m[2].Type, ShouldEqual, "map[string]string")
			})
		})
	})
}
