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

func TestConfig_LoadConfig(t *testing.T) {

	Convey("Given I load a regolithe.ini", t, func() {

		cfg, err := LoadConfig("./tests/regolithe.ini")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then cfg should be correctly initialized", func() {
			So(cfg.Author, ShouldEqual, "aporeto")
			So(cfg.Copyright, ShouldEqual, "aporeto")
			So(cfg.Description, ShouldEqual, "basic test model")
			So(cfg.Email, ShouldEqual, "contact@aporeto.com")
			So(cfg.Name, ShouldEqual, "testmodel")
			So(cfg.ProductName, ShouldEqual, "Fixture")
			So(cfg.URL, ShouldEqual, "aporeto.com")
			So(cfg.Version, ShouldEqual, "1.0")
		})

		Convey("Then I call key on valid section and key", func() {
			So(cfg.Key("test", "key"), ShouldEqual, "value")
		})

		Convey("Then I call key on valid section and invalid key", func() {
			So(cfg.Key("test", "notkey"), ShouldEqual, "")
		})

		Convey("Then I call key on invalid section", func() {
			So(cfg.Key("nottest", "key"), ShouldEqual, "")
		})
	})
}
