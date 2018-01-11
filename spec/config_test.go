package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfog_LoadConfig(t *testing.T) {

	Convey("Given I load a monolithe.ini", t, func() {

		cfg, err := LoadConfig("./tests/monolithe.ini")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then cfg should be correctly initialized", func() {
			So(cfg.Author, ShouldEqual, "aporeto")
			So(cfg.Copyright, ShouldEqual, "aporeto")
			So(cfg.Description, ShouldEqual, "basic test model")
			So(cfg.Email, ShouldEqual, "contact@aporeto.com")
			So(cfg.Name, ShouldEqual, "testmodel")
			So(cfg.Output, ShouldEqual, "codegen")
			So(cfg.ProductName, ShouldEqual, "Fixture")
			So(cfg.URL, ShouldEqual, "aporeto.com")
			So(cfg.Version, ShouldEqual, "1.0")
		})
	})
}
