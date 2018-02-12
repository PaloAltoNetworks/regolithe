package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRelationshio_NewRelationship(t *testing.T) {

	Convey("Given I call NewRelationship", t, func() {

		r := NewRelationship()

		Convey("Then r should be correctlt initialized", func() {
			So(r.AllowsCreate, ShouldResemble, map[string]struct{}{})
			So(r.AllowsDelete, ShouldResemble, map[string]struct{}{})
			So(r.AllowsGet, ShouldResemble, map[string]struct{}{})
			So(r.AllowsUpdate, ShouldResemble, map[string]struct{}{})
		})
	})
}

func TestRelationshio_GetSet(t *testing.T) {

	Convey("Given I call NewRelationship", t, func() {

		r := NewRelationship()
		r.Set("get", "a", "b", "c")
		r.Set("update", "a", "b", "d")
		r.Set("delete", "toto")
		r.Set("create")

		Convey("When I call Get on 'get'", func() {

			names := r.Get("get")

			Convey("Then it should be correct", func() {
				So(names, ShouldResemble, []string{"a", "b", "c"})
			})
		})

		Convey("When I call Get on 'update'", func() {

			names := r.Get("update")

			Convey("Then it should be correct", func() {
				So(names, ShouldResemble, []string{"a", "b", "d"})
			})
		})

		Convey("When I call Get on 'delete'", func() {

			names := r.Get("delete")

			Convey("Then it should be correct", func() {
				So(names, ShouldResemble, []string{"toto"})
			})
		})

		Convey("When I call Get on 'create'", func() {

			names := r.Get("create")

			Convey("Then it should be correct", func() {
				So(names, ShouldBeNil)
			})
		})

		Convey("When I call Get on 'nope'", func() {

			Convey("Then it should panic", func() {
				So(func() { r.Get("nope") }, ShouldPanic)
			})
		})

		Convey("When I call Set on 'nope'", func() {

			Convey("Then it should panic", func() {
				So(func() { r.Set("nope") }, ShouldPanic)
			})
		})
	})

}
