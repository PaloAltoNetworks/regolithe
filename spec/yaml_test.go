package spec

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	yaml "gopkg.in/yaml.v2"
)

func TestModel_orderedKeys(t *testing.T) {

	Convey("Given I have a mode", t, func() {

		m := &Model{
			Aliases:      []string{"alias1", "alias2"},
			AllowsCreate: false,
			AllowsDelete: true,
			AllowsGet:    true,
			AllowsUpdate: true,
			Description:  "desc",
			EntityName:   "entity name",
			Extends:      []string{"ext1", "ext2"},
			IsRoot:       true,
			Package:      "package",
			ResourceName: "resource name",
			RestName:     "rest name",
			Private:      true,

			EntityNamePlural: "plural",
		}

		Convey("When I call toYAMLMapSlice ", func() {

			ms := toYAMLMapSlice(m)

			Convey("Then the result should be correct", func() {
				So(len(ms), ShouldEqual, 12)

				// ShouldResemble is weird and doesn't work here.
				So(fmt.Sprintf("%#v", ms[0]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "rest_name", Value: "rest name"}))
				So(fmt.Sprintf("%#v", ms[1]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "resource_name", Value: "resource name"}))
				So(fmt.Sprintf("%#v", ms[2]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "entity_name", Value: "entity name"}))
				So(fmt.Sprintf("%#v", ms[3]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "package", Value: "package"}))
				So(fmt.Sprintf("%#v", ms[4]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "description", Value: "desc"}))
				So(fmt.Sprintf("%#v", ms[5]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "aliases", Value: []string{"alias1", "alias2"}}))
				So(fmt.Sprintf("%#v", ms[6]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "private", Value: true}))
				So(fmt.Sprintf("%#v", ms[7]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "get", Value: true}))
				So(fmt.Sprintf("%#v", ms[8]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "update", Value: true}))
				So(fmt.Sprintf("%#v", ms[9]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "delete", Value: true}))
				So(fmt.Sprintf("%#v", ms[10]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "extends", Value: []string{"ext1", "ext2"}}))
				So(fmt.Sprintf("%#v", ms[11]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "root", Value: true}))
			})
		})
	})
}
