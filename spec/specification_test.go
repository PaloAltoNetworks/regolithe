package spec

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecification_NewSpecification(t *testing.T) {

	Convey("Given I create a new specification", t, func() {

		spec := NewSpecification()

		Convey("Then the spec should be correctly initialized", func() {
			So(spec.APIs, ShouldNotBeNil)
			So(spec.Attributes, ShouldNotBeNil)
			So(spec.Extends, ShouldNotBeNil)
		})
	})
}

func TestSpecification_Getters(t *testing.T) {

	Convey("Given I have a new API", t, func() {

		s := &Specification{
			model: &model{
				EntityName:   "Test",
				ResourceName: "tests",
				RestName:     "test",
				AllowsCreate: true,
				AllowsDelete: true,
				AllowsGet:    true,
				AllowsUpdate: true,
			},
			Attributes: []*Attribute{
				&Attribute{
					Identifier: false,
					Name:       "not-id",
				},
				&Attribute{
					Identifier: true,
					Name:       "id",
				},
			},
		}

		s.buildAttributesInfo() // nolint: errcheck

		Convey("Then the getters should work", func() {
			So(s.GetRestName(), ShouldEqual, "test")
			So(s.GetEntityName(), ShouldEqual, "Test")
			So(s.GetAllowsGet(), ShouldBeTrue)
			So(s.GetAllowsUpdate(), ShouldBeTrue)
			So(s.GetAllowsCreate(), ShouldBeTrue)
			So(s.GetAllowsDelete(), ShouldBeTrue)
			So(s.Identifier().Name, ShouldEqual, "id")
		})
	})
}

func TestSpecification_TypeProviders(t *testing.T) {

	Convey("Given I have a new API", t, func() {

		s := &Specification{
			model: &model{},
			Attributes: []*Attribute{
				&Attribute{
					Name:          "not-id",
					ConvertedName: "not-id",
					TypeProvider:  "toto",
				},
				&Attribute{
					ConvertedName: "id",
					TypeProvider:  "titi",
				},
				&Attribute{
					ConvertedName: "id2",
					TypeProvider:  "titi",
				},
				&Attribute{},
			},
		}

		s.buildAttributesInfo() // nolint: errcheck

		Convey("When I call TypeProviders", func() {

			providers := s.TypeProviders()

			Convey("Then the providers should be correct", func() {
				So(providers, ShouldResemble, []string{"toto", "titi"})
			})
		})
	})
}

func TestSpecification_AttributeInitializers(t *testing.T) {

	Convey("Given I have a new API", t, func() {

		s := &Specification{
			Attributes: []*Attribute{
				&Attribute{
					Identifier:    false,
					Name:          "not-id",
					ConvertedName: "not-id",
					DefaultValue:  "default1",
					Type:          AttributeTypeString,
				},
				&Attribute{
					Identifier:    true,
					Name:          "id",
					ConvertedName: "id",
					Initializer:   "init",
				},
			},
		}

		s.buildAttributesInfo() // nolint: errcheck

		Convey("When I call AttributeInitializers", func() {

			inits := s.AttributeInitializers()

			Convey("Then the initializers should be correct", func() {
				So(inits["id"], ShouldEqual, "init")
				So(inits["not-id"], ShouldEqual, `"default1"`)
				So(len(inits), ShouldEqual, 2)
			})
		})
	})
}

func TestSpecification_OrderingAttributes(t *testing.T) {

	Convey("Given I have a new API", t, func() {

		s := &Specification{
			Attributes: []*Attribute{
				&Attribute{
					DefaultOrder: true,
					Name:         "a1",
				},
				&Attribute{
					DefaultOrder: true,
					Name:         "a2",
				},
				&Attribute{
					Name: "a3",
				},
			},
		}

		s.buildAttributesInfo() // nolint: errcheck

		Convey("When I call OrderingAttributes", func() {

			o := s.OrderingAttributes()

			Convey("Then the orderingAttributes should be correct", func() {
				So(len(o), ShouldEqual, 2)
				So(o[0].Name, ShouldEqual, "a1")
				So(o[1].Name, ShouldEqual, "a2")
			})
		})
	})
}

func TestSpecification_AttributeMap(t *testing.T) {

	Convey("Given I load the task specification file", t, func() {

		spec, err := LoadSpecification("./tests/task.spec")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the attribute map should be correctly built", func() {
			So(len(spec.Attributes), ShouldEqual, 3)
			So(spec.Attribute("name").Name, ShouldEqual, "name")
			So(spec.Attribute("description").Name, ShouldEqual, "description")
			So(spec.Attribute("status").Name, ShouldEqual, "status")
		})
	})
}

func TestSpecification_BuildAttributeNames(t *testing.T) {

	Convey("Given I create a specification with the same attribute twice.", t, func() {

		spec := &Specification{
			model: &model{},
			Attributes: []*Attribute{
				&Attribute{
					Name: "a",
				},
				&Attribute{
					Name: "a",
				},
			},
		}

		Convey("When I call BuildAttributeNames", func() {

			err := spec.buildAttributesInfo()

			Convey("Then err Should Not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSpecification_APIMap(t *testing.T) {

	Convey("Given I load the root specification file", t, func() {

		spec, err := LoadSpecification("./tests/root.spec")

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the child api map should be correctly built", func() {
			So(len(spec.APIs), ShouldEqual, 2)
			So(spec.API("list").RestName, ShouldEqual, "list")
			So(spec.API("user").RestName, ShouldEqual, "user")
		})
	})
}

func TestSpecification_BuildAPINames(t *testing.T) {

	Convey("Given I create a specification with the same child api twice.", t, func() {

		spec := &Specification{
			APIs: []*API{
				&API{
					RestName: "a",
				},
				&API{
					RestName: "a",
				},
			},
		}

		Convey("When I call BuildAPINames", func() {

			err := spec.buildAPIsInfo()

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSpecification_LoadSpecification(t *testing.T) {

	Convey("Given I load a non existing specification file", t, func() {

		_, err := LoadSpecification("./tests/not.spec")

		Convey("Then err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given I load a bad formatted specification file", t, func() {

		_, err := LoadSpecification("./tests/task.spec.bad")

		Convey("Then err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given I load the root specification file", t, func() {

		spec, err := LoadSpecification("./tests/root.spec")
		apis := spec.APIs

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the spec should be correctly initialized", func() {
			So(spec.AllowsGet, ShouldBeTrue)
			So(spec.AllowsCreate, ShouldBeFalse)
			So(spec.AllowsDelete, ShouldBeFalse)
			So(spec.AllowsUpdate, ShouldBeFalse)
			So(spec.Description, ShouldEqual, "Root object of the API")
			So(spec.EntityName, ShouldEqual, "Root")
			So(spec.Package, ShouldEqual, "todo-list")
			So(spec.ResourceName, ShouldEqual, "root")
			So(spec.RestName, ShouldEqual, "root")
			So(spec.Extends, ShouldResemble, []string{})
			So(spec.IsRoot, ShouldBeTrue)
			So(spec.Aliases, ShouldResemble, []string{})
		})

		Convey("Then the number of api should be correct", func() {
			So(len(spec.APIs), ShouldEqual, 2)
		})

		Convey("Then the list child API should be correct", func() {
			So(apis[0].AllowsBulkCreate, ShouldBeFalse)
			So(apis[0].AllowsBulkDelete, ShouldBeFalse)
			So(apis[0].AllowsBulkUpdate, ShouldBeFalse)
			So(apis[0].AllowsCreate, ShouldBeTrue)
			So(apis[0].AllowsDelete, ShouldBeFalse)
			So(apis[0].Deprecated, ShouldBeFalse)
			So(apis[0].AllowsGet, ShouldBeTrue)
			So(apis[0].Relationship, ShouldEqual, APIRelationshipRoot)
			So(apis[0].RestName, ShouldEqual, "list")
			So(apis[0].AllowsUpdate, ShouldBeFalse)
		})

		Convey("Then the user child API should be correct", func() {
			So(apis[1].AllowsBulkCreate, ShouldBeFalse)
			So(apis[1].AllowsBulkDelete, ShouldBeFalse)
			So(apis[1].AllowsBulkUpdate, ShouldBeFalse)
			So(apis[1].AllowsCreate, ShouldBeTrue)
			So(apis[1].AllowsDelete, ShouldBeFalse)
			So(apis[1].Deprecated, ShouldBeFalse)
			So(apis[1].AllowsGet, ShouldBeTrue)
			So(apis[1].Relationship, ShouldEqual, APIRelationshipRoot)
			So(apis[1].RestName, ShouldEqual, "user")
			So(apis[1].AllowsUpdate, ShouldBeFalse)
		})

	})

	Convey("Given I load the task specification file", t, func() {

		spec, err := LoadSpecification("./tests/task.spec")
		attrs := spec.Attributes

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the spec should be correctly initialized", func() {
			So(spec.AllowsGet, ShouldBeTrue)
			So(spec.AllowsCreate, ShouldBeFalse)
			So(spec.AllowsDelete, ShouldBeTrue)
			So(spec.AllowsUpdate, ShouldBeTrue)
			So(spec.Description, ShouldEqual, "Represent a task to do in a listd")
			So(spec.EntityName, ShouldEqual, "Task")
			So(spec.Package, ShouldEqual, "todo-list")
			So(spec.ResourceName, ShouldEqual, "tasks")
			So(spec.RestName, ShouldEqual, "task")
			So(spec.Extends, ShouldResemble, []string{"@base"})
			So(spec.IsRoot, ShouldBeFalse)
			So(spec.Aliases, ShouldResemble, []string{"tsk"})
		})

		Convey("Then the number of attributes should be correct", func() {
			So(len(spec.Attributes), ShouldEqual, 3)
		})

		Convey("Then the spec attribute description be correctly initialized", func() {

			So(attrs[0].AllowedChars, ShouldBeEmpty)
			So(attrs[0].AllowedChoices, ShouldBeEmpty)
			So(attrs[0].Autogenerated, ShouldBeFalse)
			So(attrs[0].CreationOnly, ShouldBeFalse)
			So(attrs[0].DefaultOrder, ShouldBeFalse)
			So(attrs[0].DefaultValue, ShouldBeNil)
			So(attrs[0].Deprecated, ShouldBeFalse)
			So(attrs[0].Description, ShouldEqual, "The description")
			So(attrs[0].Exposed, ShouldBeTrue)
			So(attrs[0].Filterable, ShouldBeTrue)
			So(attrs[0].ForeignKey, ShouldBeFalse)
			So(attrs[0].Format, ShouldEqual, AttributeFormatFree)
			So(attrs[0].Getter, ShouldBeFalse)
			So(attrs[0].Identifier, ShouldBeFalse)
			So(attrs[0].Index, ShouldBeFalse)
			So(attrs[0].MaxLength, ShouldEqual, 0)
			So(attrs[0].MaxValue, ShouldEqual, 0.0)
			So(attrs[0].MinLength, ShouldEqual, 0)
			So(attrs[0].MinValue, ShouldEqual, 0.0)
			So(attrs[0].Name, ShouldEqual, "description")
			So(attrs[0].Orderable, ShouldBeTrue)
			So(attrs[0].PrimaryKey, ShouldBeFalse)
			So(attrs[0].ReadOnly, ShouldBeFalse)
			So(attrs[0].Required, ShouldBeFalse)
			So(attrs[0].Secret, ShouldBeFalse)
			So(attrs[0].Setter, ShouldBeFalse)
			So(attrs[0].Stored, ShouldBeTrue)
			So(attrs[0].SubType, ShouldEqual, "")
			So(attrs[0].Transient, ShouldBeFalse)
			So(attrs[0].Type, ShouldEqual, AttributeTypeString)
			So(attrs[0].Unique, ShouldBeFalse)
		})

		Convey("Then the spec attribute name be correctly initialized", func() {

			So(attrs[1].AllowedChars, ShouldBeEmpty)
			So(attrs[1].AllowedChoices, ShouldBeEmpty)
			So(attrs[1].Autogenerated, ShouldBeFalse)
			So(attrs[1].CreationOnly, ShouldBeFalse)
			So(attrs[1].DefaultOrder, ShouldBeFalse)
			So(attrs[1].DefaultValue, ShouldBeNil)
			So(attrs[1].Deprecated, ShouldBeFalse)
			So(attrs[1].Description, ShouldEqual, "The name")
			So(attrs[1].Exposed, ShouldBeTrue)
			So(attrs[1].Filterable, ShouldBeTrue)
			So(attrs[1].ForeignKey, ShouldBeFalse)
			So(attrs[1].Format, ShouldEqual, AttributeFormatFree)
			So(attrs[1].Getter, ShouldBeTrue)
			So(attrs[1].Identifier, ShouldBeFalse)
			So(attrs[1].Index, ShouldBeFalse)
			So(attrs[1].MaxLength, ShouldEqual, 0)
			So(attrs[1].MaxValue, ShouldEqual, 0.0)
			So(attrs[1].MinLength, ShouldEqual, 0)
			So(attrs[1].MinValue, ShouldEqual, 0.0)
			So(attrs[1].Name, ShouldEqual, "name")
			So(attrs[1].Orderable, ShouldBeTrue)
			So(attrs[1].PrimaryKey, ShouldBeFalse)
			So(attrs[1].ReadOnly, ShouldBeFalse)
			So(attrs[1].Required, ShouldBeTrue)
			So(attrs[1].Secret, ShouldBeFalse)
			So(attrs[1].Setter, ShouldBeTrue)
			So(attrs[1].Stored, ShouldBeTrue)
			So(attrs[1].SubType, ShouldEqual, "")
			So(attrs[1].Transient, ShouldBeFalse)
			So(attrs[1].Type, ShouldEqual, AttributeTypeString)
			So(attrs[1].Unique, ShouldBeFalse)
		})

		Convey("Then the spec attribbute status be correctly initialized", func() {

			So(attrs[2].AllowedChars, ShouldBeEmpty)
			So(attrs[2].AllowedChoices, ShouldResemble, []string{"DONE", "PROGRESS", "TODO"})
			So(attrs[2].Autogenerated, ShouldBeFalse)
			So(attrs[2].CreationOnly, ShouldBeFalse)
			So(attrs[2].DefaultOrder, ShouldBeFalse)
			So(attrs[2].DefaultValue, ShouldResemble, "TODO")
			So(attrs[2].Deprecated, ShouldBeFalse)
			So(attrs[2].Description, ShouldEqual, "The status of the task")
			So(attrs[2].Exposed, ShouldBeTrue)
			So(attrs[2].Filterable, ShouldBeTrue)
			So(attrs[2].ForeignKey, ShouldBeFalse)
			So(attrs[2].Format, ShouldEqual, "")
			So(attrs[2].Getter, ShouldBeFalse)
			So(attrs[2].Identifier, ShouldBeFalse)
			So(attrs[2].Index, ShouldBeFalse)
			So(attrs[2].MaxLength, ShouldEqual, 0)
			So(attrs[2].MaxValue, ShouldEqual, 0.0)
			So(attrs[2].MinLength, ShouldEqual, 0)
			So(attrs[2].MinValue, ShouldEqual, 0.0)
			So(attrs[2].Name, ShouldEqual, "status")
			So(attrs[2].Orderable, ShouldBeTrue)
			So(attrs[2].PrimaryKey, ShouldBeFalse)
			So(attrs[2].ReadOnly, ShouldBeFalse)
			So(attrs[2].Required, ShouldBeFalse)
			So(attrs[2].Secret, ShouldBeFalse)
			So(attrs[2].Setter, ShouldBeFalse)
			So(attrs[2].Stored, ShouldBeTrue)
			So(attrs[2].SubType, ShouldEqual, "")
			So(attrs[2].Transient, ShouldBeFalse)
			So(attrs[2].Type, ShouldEqual, AttributeTypeEnum)
			So(attrs[2].Unique, ShouldBeFalse)
		})

		Convey("When I apply the base specification", func() {

			base, err := LoadSpecification("./tests/@base.spec")
			spec.ApplyBaseSpecifications(base) // nolint: errcheck

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the number of attributes should be correct", func() {
				So(len(spec.Attributes), ShouldEqual, 6)
			})

			Convey("Then the additional attributes should have been applied", func() {
				So(spec.Attribute("ID").Name, ShouldEqual, "ID")
				So(spec.Attribute("ID").Autogenerated, ShouldBeTrue)
				So(spec.Attribute("ID").Description, ShouldEqual, "The identifier")
				So(spec.Attribute("ID").Identifier, ShouldBeTrue)
				So(spec.Attribute("ID").PrimaryKey, ShouldBeTrue)
				So(spec.Attribute("ID").ReadOnly, ShouldBeTrue)

				So(spec.Attribute("parentID").Name, ShouldEqual, "parentID")
				So(spec.Attribute("parentID").Autogenerated, ShouldBeTrue)
				So(spec.Attribute("parentID").Description, ShouldEqual, "The identifier of the parent of the object")
				So(spec.Attribute("parentID").ForeignKey, ShouldBeTrue)
				So(spec.Attribute("parentID").Identifier, ShouldBeFalse)
				So(spec.Attribute("parentID").PrimaryKey, ShouldBeFalse)
				So(spec.Attribute("parentID").ReadOnly, ShouldBeTrue)
			})
		})
	})
}
