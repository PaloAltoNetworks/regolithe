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

func TestValidationMapping_LoadValidationMapping(t *testing.T) {

	Convey("Given I load a type mapping", t, func() {

		tm, err := LoadValidationMapping("./tests/_validation.mapping")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the file should be loaded", func() {
			So(tm, ShouldResemble, ValidationMapping{
				"$username": map[string]*ValidationMap{
					"test": {
						Name: "validate.CheckUserName",
					},
				},
				"$nospace": map[string]*ValidationMap{
					"test": {
						Name: "validate.NoSpace",
					},
				},
				"$nocap": map[string]*ValidationMap{
					"test": {
						Name:   "nocapper.NoCap",
						Import: "github.com/aporeto/nocapper",
					},
					"other": {
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
