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
