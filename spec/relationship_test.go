package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRelationshio_NewRelationship(t *testing.T) {

	Convey("Given I call NewRelationship", t, func() {

		r := NewRelationship()

		Convey("Then r should be correctlt initialized", func() {
			So(r.Create, ShouldResemble, map[string]*RelationAction{})
			So(r.Delete, ShouldResemble, map[string]*RelationAction{})
			So(r.Get, ShouldResemble, map[string]*RelationAction{})
			So(r.Update, ShouldResemble, map[string]*RelationAction{})
		})
	})
}

// func TestRelationshio_GetSet(t *testing.T) {

// 	Convey("Given I call NewRelationship", t, func() {

// 		r := NewRelationship()
// 		r.Set("get", "a", &RelationAction{})
// 		r.Set("update", "a", &RelationAction{})
// 		r.Set("delete", "toto", &RelationAction{})
// 		r.Set("create")

// 		Convey("When I call Get on 'get'", func() {

// 			names := r.Get("get")

// 			Convey("Then it should be correct", func() {
// 				So(names, ShouldResemble, []string{"a", "b", "c"})
// 			})
// 		})

// 		Convey("When I call Get on 'update'", func() {

// 			names := r.Get("update")

// 			Convey("Then it should be correct", func() {
// 				So(names, ShouldResemble, []string{"a", "b", "d"})
// 			})
// 		})

// 		Convey("When I call Get on 'delete'", func() {

// 			names := r.Get("delete")

// 			Convey("Then it should be correct", func() {
// 				So(names, ShouldResemble, []string{"toto"})
// 			})
// 		})

// 		Convey("When I call Get on 'create'", func() {

// 			names := r.Get("create")

// 			Convey("Then it should be correct", func() {
// 				So(names, ShouldBeNil)
// 			})
// 		})

// 		Convey("When I call Get on 'nope'", func() {

// 			Convey("Then it should panic", func() {
// 				So(func() { r.Get("nope") }, ShouldPanic)
// 			})
// 		})

// 		Convey("When I call Set on 'nope'", func() {

// 			Convey("Then it should panic", func() {
// 				So(func() { r.Set("nope") }, ShouldPanic)
// 			})
// 		})
// 	})

// }
