package spec

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec_LoadSpecificationDir(t *testing.T) {

	Convey("Given I load a spec folder", t, func() {

		set, err := NewSpecificationSet(
			"./tests",
			func(n string) string {
				return strings.ToUpper(n)
			},
			func(typ AttributeType, subtype string) (string, string) {
				if typ == AttributeTypeString {
					return "String", ""
				}
				if typ == AttributeTypeTime {
					return "time.Time", "time"
				}
				return string(typ), ""
			},
			"elemental",
		)

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the specification set should have 4 entries", func() {
			So(set.Len(), ShouldEqual, 4)
		})

		Convey("Then calling Specifications should return the sorted list", func() {
			ss := set.Specifications()
			So(ss[0].RestName, ShouldEqual, "list")
			So(ss[1].RestName, ShouldEqual, "root")
			So(ss[2].RestName, ShouldEqual, "task")
			So(ss[3].RestName, ShouldEqual, "user")
		})

		Convey("Then the relationships should be correct", func() {
			rs := set.Relationships()
			So(rs["List"].Get("get"), ShouldResemble, []string{"root"})
			So(rs["List"].Get("create"), ShouldResemble, []string{"root"})
			So(rs["List"].Get("update"), ShouldResemble, []string{"root"})
			So(rs["List"].Get("delete"), ShouldResemble, []string{"root"})
			So(rs["Task"].Get("get"), ShouldResemble, []string{"list", "root"})
			So(rs["Task"].Get("create"), ShouldResemble, []string{"list"})
			So(rs["Task"].Get("update"), ShouldResemble, []string{"root"})
			So(rs["Task"].Get("delete"), ShouldResemble, []string{"root"})
			So(rs["User"].Get("get"), ShouldResemble, []string{"list", "root"})
			So(rs["User"].Get("create"), ShouldResemble, []string{"root"})
			So(rs["User"].Get("update"), ShouldResemble, []string{"list", "root"})
			So(rs["User"].Get("delete"), ShouldResemble, []string{"root"})
		})

		Convey("Then the specification set should be correct", func() {
			So(len(set.Specification("task").Attributes), ShouldEqual, 6)
			So(len(set.Specification("root").Attributes), ShouldEqual, 0)
			So(len(set.Specification("list").Attributes), ShouldEqual, 6)
			So(len(set.Specification("user").Attributes), ShouldEqual, 6)
		})

		Convey("Then the API linking should be correct", func() {
			So(set.Specification("root").API("user").Specification(), ShouldEqual, set.Specification("user"))
			So(set.Specification("root").API("list").Specification(), ShouldEqual, set.Specification("list"))
			So(set.Specification("list").API("task").Specification(), ShouldEqual, set.Specification("task"))
		})

		Convey("Then the config should be correctly loaded", func() {
			So(set.Configuration.Name, ShouldEqual, "testmodel")
		})

		Convey("Then the type mapping should be correctly loaded", func() {
			m, _ := set.ExternalTypes.Mapping("elemental", "string_map")
			So(m.Type, ShouldEqual, "map[string]string")
		})

		Convey("Then the type conversion should have worked", func() {
			So(set.Specification("list").Attribute("name").ConvertedName, ShouldEqual, "NAME")
			So(set.Specification("list").Attribute("name").TypeProvider, ShouldEqual, "")

			So(set.Specification("list").Attribute("name").ConvertedType, ShouldEqual, "String")
			So(set.Specification("list").Attribute("name").TypeProvider, ShouldEqual, "")

			So(set.Specification("list").Attribute("date").ConvertedType, ShouldEqual, "time.Time")
			So(set.Specification("list").Attribute("date").TypeProvider, ShouldEqual, "time")
		})
	})
}
