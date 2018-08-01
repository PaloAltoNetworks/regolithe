package spec

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecification_NewSpecification(t *testing.T) {

	Convey("Given I create a new specification", t, func() {

		spec := NewSpecification()

		Convey("Then the spec should be correctly initialized", func() {
			So(spec, ShouldNotBeNil)
		})
	})
}

func TestSpecification_Validate(t *testing.T) {

	Convey("Given I have a specification with no validation error", t, func() {

		s := &specification{
			RawModel: &Model{
				ResourceName: "things",
				RestName:     "thing",
				Description:  "desc.",
				EntityName:   "toto",
				Package:      "package",
			},
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name:        "attr1",
						Description: "desc.",
						Type:        "string",
					},
				},
			},
		}

		Convey("When I call validate", func() {

			err := s.Validate()

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a specification with validation error", t, func() {

		s := &specification{
			RawModel: &Model{
				ResourceName: "things",
				RestName:     "thing",
				Description:  "desc.",
				EntityName:   "toto",
			},
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Identifier: false,
						Name:       "not-id",
					},
					&Attribute{
						Name:        "id",
						Type:        "coucou",
						Description: "wala",
					},
				},
			},
		}

		Convey("When I call validate", func() {

			errs := s.Validate()

			Convey("Then err should not be nil", func() {
				So(errs, ShouldNotBeNil)
				So(formatValidationErrors(errs).Error(), ShouldEqual, `thing.spec: schema error: attributes.v1.1.type: attributes.v1.1.type must be one of the following: "string", "integer", "float", "boolean", "enum", "list", "object", "time", "external"
thing.spec: schema error: description: description is required
thing.spec: schema error: package: package is required
thing.spec: schema error: type: type is required`)
			})
		})
	})

	Convey("Given I have an abstract with no validation error", t, func() {

		s := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name:        "attr1",
						Description: "desc.",
						Type:        "string",
					},
				},
			},
		}

		Convey("When I call validate", func() {

			err := s.Validate()

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given I have an abstract with validation error", t, func() {

		s := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name:        "attr1",
						Description: "desc.",
						Type:        "string",
					},
				},
			},
			RawRelations: []*Relation{
				&Relation{
					RestName: "a",
				},
			},
		}

		Convey("When I call validate", func() {

			err := s.Validate()

			Convey("Then err should not be nil", func() {
				So(len(err), ShouldEqual, 1)
				So(err[0].Error(), ShouldEqual, ".: schema error: relations: Additional property relations is not allowed")
			})
		})
	})
}

func TestSpecification_Getters(t *testing.T) {

	Convey("Given I have a new API", t, func() {

		s := &specification{
			RawModel: &Model{
				EntityName:   "Test",
				ResourceName: "tests",
				RestName:     "test",
				Delete:       &RelationAction{},
				Get:          &RelationAction{},
				Update:       &RelationAction{},
			},
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Identifier: false,
						Name:       "not-id",
					},
					&Attribute{
						Identifier: true,
						Name:       "id",
					},
				},
			},
		}

		s.buildAttributesMapping() // nolint: errcheck

		Convey("Then the getters should work", func() {
			So(s.Identifier().Name, ShouldEqual, "id")
		})
	})
}

func TestSpecification_TypeProviders(t *testing.T) {

	Convey("Given I have a new API", t, func() {

		s := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
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
			},
		}

		s.buildAttributesMapping() // nolint: errcheck

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

		s := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
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
			},
		}

		s.buildAttributesMapping() // nolint: errcheck

		Convey("When I call AttributeInitializers", func() {

			inits := s.AttributeInitializers("v1")

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

		s := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{

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
			},
		}

		s.buildAttributesMapping() // nolint: errcheck

		Convey("When I call OrderingAttributes", func() {

			o := s.OrderingAttributes("v1")

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

		spec, err := LoadSpecification("./tests/task.spec", true)

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the attribute map should be correctly built", func() {
			So(len(spec.Attributes("v1")), ShouldEqual, 3)
			So(spec.Attribute("name", "v1").Name, ShouldEqual, "name")
			So(spec.Attribute("description", "v1").Name, ShouldEqual, "description")
			So(spec.Attribute("status", "v1").Name, ShouldEqual, "status")
		})
	})
}

func TestSpecification_buildAttributesMapping(t *testing.T) {

	Convey("Given I create a specification with the same attribute twice.", t, func() {

		spec := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name: "a",
					},
					&Attribute{
						Name: "a",
					},
				},
			},
		}

		Convey("When I call buildAttributesMapping", func() {

			err := spec.buildAttributesMapping()

			Convey("Then err Should Not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSpecification_APIMap(t *testing.T) {

	Convey("Given I load the root specification file", t, func() {

		spec, err := LoadSpecification("./tests/root.spec", true)

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the relation map should be correctly built", func() {
			So(len(spec.Relations()), ShouldEqual, 2)
			So(spec.Relation("list").RestName, ShouldEqual, "list")
			So(spec.Relation("user").RestName, ShouldEqual, "user")
		})
	})
}

func TestSpecification_buildRelationsMapping(t *testing.T) {

	Convey("Given I create a specification with the same relation twice.", t, func() {

		spec := &specification{
			RawRelations: []*Relation{
				&Relation{
					RestName: "a",
				},
				&Relation{
					RestName: "a",
				},
			},
		}

		Convey("When I call buildRelationsMapping", func() {

			err := spec.buildRelationsMapping()

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSpecification_LoadSpecification(t *testing.T) {

	Convey("Given I load a non existing specification file", t, func() {

		_, err := LoadSpecification("./tests/not.spec", true)

		Convey("Then err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given I load a bad formatted specification file", t, func() {

		_, err := LoadSpecification("./tests/task.spec.bad", true)

		Convey("Then err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given I load the root specification file", t, func() {

		spec, err := LoadSpecification("./tests/root.spec", true)
		rels := spec.Relations()

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the spec should be correctly initialized", func() {

			model := spec.Model()

			So(model.Get, ShouldNotBeNil)
			So(model.Delete, ShouldBeNil)
			So(model.Update, ShouldBeNil)
			So(model.Description, ShouldEqual, "Root object of the API.")
			So(model.EntityName, ShouldEqual, "Root")
			So(model.Package, ShouldEqual, "todo-list")
			So(model.ResourceName, ShouldEqual, "root")
			So(model.RestName, ShouldEqual, "root")
			So(model.Extends, ShouldBeNil)
			So(model.IsRoot, ShouldBeTrue)
			So(model.Aliases, ShouldBeNil)
		})

		Convey("Then the number of relations should be correct", func() {
			So(len(rels), ShouldEqual, 2)
		})

		Convey("Then the list of relations should be correct", func() {
			So(rels[0].Create, ShouldNotBeNil)
			So(rels[0].Delete, ShouldBeNil)
			So(rels[0].Update, ShouldBeNil)
			So(rels[0].Get, ShouldNotBeNil)
			So(rels[0].RestName, ShouldEqual, "list")
		})

		Convey("Then the user relation should be correct", func() {
			So(rels[1].Create, ShouldNotBeNil)
			So(rels[1].Delete, ShouldBeNil)
			So(rels[1].Update, ShouldBeNil)
			So(rels[1].Get, ShouldNotBeNil)
			So(rels[1].RestName, ShouldEqual, "user")
		})

	})

	Convey("Given I load the task specification file", t, func() {

		spec, err := LoadSpecification("./tests/task.spec", true)
		attrs := spec.Attributes("v1")
		model := spec.Model()

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the spec should be correctly initialized", func() {
			So(model.Get, ShouldNotBeNil)
			So(model.Delete, ShouldNotBeNil)
			So(model.Update, ShouldNotBeNil)
			So(model.Description, ShouldEqual, "Represent a task to do in a listd.")
			So(model.EntityName, ShouldEqual, "Task")
			So(model.Package, ShouldEqual, "todo-list")
			So(model.ResourceName, ShouldEqual, "tasks")
			So(model.RestName, ShouldEqual, "task")
			So(model.Extends, ShouldResemble, []string{"@base"})
			So(model.IsRoot, ShouldBeFalse)
			So(model.Aliases, ShouldResemble, []string{"tsk"})
		})

		Convey("Then the number of attributes should be correct", func() {
			So(len(spec.Attributes("v1")), ShouldEqual, 3)
		})

		Convey("Then the spec attribute description be correctly initialized", func() {

			So(attrs[0].AllowedChars, ShouldBeEmpty)
			So(attrs[0].AllowedChoices, ShouldBeEmpty)
			So(attrs[0].Autogenerated, ShouldBeFalse)
			So(attrs[0].CreationOnly, ShouldBeFalse)
			So(attrs[0].DefaultOrder, ShouldBeFalse)
			So(attrs[0].DefaultValue, ShouldBeNil)
			So(attrs[0].Deprecated, ShouldBeFalse)
			So(attrs[0].Description, ShouldEqual, "The description.")
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
		})

		Convey("Then the spec attribute name be correctly initialized", func() {

			So(attrs[1].AllowedChars, ShouldBeEmpty)
			So(attrs[1].AllowedChoices, ShouldBeEmpty)
			So(attrs[1].Autogenerated, ShouldBeFalse)
			So(attrs[1].CreationOnly, ShouldBeFalse)
			So(attrs[1].DefaultOrder, ShouldBeFalse)
			So(attrs[1].DefaultValue, ShouldBeNil)
			So(attrs[1].Deprecated, ShouldBeFalse)
			So(attrs[1].Description, ShouldEqual, "The name.")
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
		})

		Convey("Then the spec attribbute status be correctly initialized", func() {

			So(attrs[2].AllowedChars, ShouldBeEmpty)
			So(attrs[2].AllowedChoices, ShouldResemble, []string{"DONE", "PROGRESS", "TODO"})
			So(attrs[2].Autogenerated, ShouldBeFalse)
			So(attrs[2].CreationOnly, ShouldBeFalse)
			So(attrs[2].DefaultOrder, ShouldBeFalse)
			So(attrs[2].DefaultValue, ShouldResemble, "TODO")
			So(attrs[2].Deprecated, ShouldBeFalse)
			So(attrs[2].Description, ShouldEqual, "The status of the task.")
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
		})

		Convey("When I apply the base specification", func() {

			base, err := LoadSpecification("./tests/@base.abs", true)
			spec.ApplyBaseSpecifications(base) // nolint: errcheck

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the number of attributes should be correct", func() {
				So(len(spec.Attributes("v1")), ShouldEqual, 6)
			})

			Convey("Then the additional attributes should have been applied", func() {
				So(spec.Attribute("ID", "v1").Name, ShouldEqual, "ID")
				So(spec.Attribute("ID", "v1").Autogenerated, ShouldBeTrue)
				So(spec.Attribute("ID", "v1").Description, ShouldEqual, "The identifier.")
				So(spec.Attribute("ID", "v1").Identifier, ShouldBeTrue)
				So(spec.Attribute("ID", "v1").PrimaryKey, ShouldBeTrue)
				So(spec.Attribute("ID", "v1").ReadOnly, ShouldBeTrue)

				So(spec.Attribute("parentID", "v1").Name, ShouldEqual, "parentID")
				So(spec.Attribute("parentID", "v1").Autogenerated, ShouldBeTrue)
				So(spec.Attribute("parentID", "v1").Description, ShouldEqual, "The identifier of the parent of the object.")
				So(spec.Attribute("parentID", "v1").ForeignKey, ShouldBeTrue)
				So(spec.Attribute("parentID", "v1").Identifier, ShouldBeFalse)
				So(spec.Attribute("parentID", "v1").PrimaryKey, ShouldBeFalse)
				So(spec.Attribute("parentID", "v1").ReadOnly, ShouldBeTrue)
			})
		})
	})
}

func TestSpecification_ApplyBaseSpecifications(t *testing.T) {

	Convey("Given I have a spec and an abstract", t, func() {

		s1 := &specification{
			RawModel: &Model{
				ResourceName: "things",
				RestName:     "thing",
				Description:  "desc.",
				EntityName:   "toto",
				Package:      "package",
			},
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name:        "attr1",
						Description: "desc.",
						Type:        "string",
					},
				},
			},
		}

		abs := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name:        "attr1",
						Description: "desc from abs.",
						Type:        "string",
					},
					&Attribute{
						Name:        "attr2",
						Description: "desc2",
						Type:        "string",
					},
				},
			},
		}

		Convey("When I call ApplyBaseSpecifications", func() {

			err := s1.ApplyBaseSpecifications(abs)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the extends should be applied", func() {
				So(len(s1.Attributes("v1")), ShouldEqual, 2)
				So(s1.Attribute("attr1", "v1").Name, ShouldEqual, "attr1")
				So(s1.Attribute("attr2", "v1").Name, ShouldEqual, "attr2")
			})
		})
	})

	Convey("Given I have a spec with no attributes and an abstract", t, func() {

		s1 := &specification{
			RawModel: &Model{
				ResourceName: "things",
				RestName:     "thing",
				Description:  "desc.",
				EntityName:   "toto",
				Package:      "package",
			},
		}

		abs := &specification{
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name:        "attr2",
						Description: "desc2",
						Type:        "string",
					},
				},
			},
		}

		Convey("When I call ApplyBaseSpecifications", func() {

			err := s1.ApplyBaseSpecifications(abs)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the extends should be applied", func() {
				So(len(s1.Attributes("v1")), ShouldEqual, 1)
				So(s1.Attribute("attr2", "v1").Name, ShouldEqual, "attr2")
			})
		})
	})
}

func TestSpecifications_Versionning(t *testing.T) {

	Convey("Given I have a specification", t, func() {

		s := &specification{
			RawModel: &Model{
				ResourceName: "things",
				RestName:     "thing",
				Description:  "desc.",
				EntityName:   "toto",
				Package:      "package",
			},
			RawAttributes: map[string][]*Attribute{
				"v1": []*Attribute{
					&Attribute{
						Name:        "1.1",
						Description: "desc.",
						Type:        "string",
					},
					&Attribute{
						Name:        "1.2",
						Description: "desc.",
						Type:        "string",
					},
				},
				"v2": []*Attribute{
					&Attribute{
						Name:        "2.1",
						Description: "desc.",
						Type:        "string",
					},
				},
				"v3": []*Attribute{
					&Attribute{
						Name:        "3.1",
						Description: "desc.",
						Type:        "string",
					},
					&Attribute{
						Name:        "3.2",
						Description: "desc.",
						Type:        "string",
					},
					&Attribute{
						Name:        "1.1",
						Description: "desc.",
						Type:        "string",
					},
				},
			},
		}

		Convey("When I versionsFrom with v1", func() {

			versions := s.versionsFrom("v1")

			Convey("Then versions should be correct", func() {
				So(versions, ShouldResemble, []string{"v1"})
			})
		})

		Convey("When I versionsFrom with v2", func() {

			versions := s.versionsFrom("v2")

			Convey("Then versions should be correct", func() {
				So(versions, ShouldResemble, []string{"v1", "v2"})
			})
		})

		Convey("When I versionsFrom with v3", func() {

			versions := s.versionsFrom("v3")

			Convey("Then versions should be correct", func() {
				So(versions, ShouldResemble, []string{"v1", "v2", "v3"})
			})
		})

		Convey("When I versionsFrom with vNope", func() {

			Convey("Then it should panic", func() {
				So(func() { s.versionsFrom("vNope") }, ShouldPanicWith, "Invalid version 'vNope'")
			})
		})

		Convey("When I call Attributes on v1", func() {

			attrs := s.Attributes("v1")

			Convey("Then attributes should be correct", func() {
				So(len(attrs), ShouldEqual, 2)
				So(attrs[0].Name, ShouldEqual, "1.1")
				So(attrs[1].Name, ShouldEqual, "1.2")
			})
		})

		Convey("When I call Attributes on v2", func() {

			attrs := s.Attributes("v2")

			Convey("Then attributes should be correct", func() {
				So(len(attrs), ShouldEqual, 3)
				So(attrs[0].Name, ShouldEqual, "1.1")
				So(attrs[1].Name, ShouldEqual, "1.2")
				So(attrs[2].Name, ShouldEqual, "2.1")
			})
		})

		Convey("When I call Attributes on v3", func() {

			attrs := s.Attributes("v3")

			Convey("Then attributes should be correct", func() {
				So(len(attrs), ShouldEqual, 5)
				So(attrs[0].Name, ShouldEqual, "1.1")
				So(attrs[1].Name, ShouldEqual, "1.2")
				So(attrs[2].Name, ShouldEqual, "2.1")
				So(attrs[3].Name, ShouldEqual, "3.1")
				So(attrs[4].Name, ShouldEqual, "3.2")
			})
		})
	})
}

func TestSpecification_Write(t *testing.T) {

	Convey("Given I load the task specification file", t, func() {

		spec, err := LoadSpecification("./tests/list.spec", true)

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When I call write", func() {
			buf := bytes.NewBuffer(nil)

			err = spec.Write(buf)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then buff should be correct", func() {
				So(buf.String(), ShouldEqual, `# Model
model:
  rest_name: list
  resource_name: lists
  entity_name: List
  package: todo-list
  description: Represent a a list of task to do.
  aliases:
  - lst
  get:
    description: Retrieves the list with the given ID.
    global_parameters:
    - sharedParameterA
    - sharedParameterB
    parameters:
      entries:
      - name: lgp1
        description: this is lgp1.
        type: string
        example_value: lgp1

      - name: lgp2
        description: this is lgp2.
        type: boolean
        example_value: "true"
  update:
    description: Updates the list with the given ID.
    parameters:
      entries:
      - name: lup1
        description: this is lup1.
        type: string
        example_value: lup1

      - name: lup2
        description: this is lup2.
        type: boolean
        example_value: "true"
  delete:
    description: Deletes the list with the given ID.
    parameters:
      entries:
      - name: ldp1
        description: this is ldp1.
        type: string
        example_value: ldp1

      - name: ldp2
        description: this is ldp2.
        type: boolean
        example_value: "true"
  extends:
  - '@base'

# Attributes
attributes:
  v1:
  - name: creationOnly
    description: This attribute is creation only.
    type: string
    exposed: true
    stored: true
    creation_only: true
    filterable: true
    format: free
    orderable: true

  - name: date
    description: The date.
    type: time
    exposed: true
    stored: true
    filterable: true
    orderable: true

  - name: description
    description: The description.
    type: string
    exposed: true
    stored: true
    filterable: true
    format: free
    orderable: true

  - name: name
    description: The name.
    type: string
    exposed: true
    stored: true
    required: true
    example_value: the name
    filterable: true
    format: free
    getter: true
    setter: true
    orderable: true

  - name: readOnly
    description: This attribute is readonly.
    type: string
    exposed: true
    stored: true
    read_only: true
    filterable: true
    format: free
    orderable: true

  - name: slice
    description: this is a slice.
    type: list
    exposed: true
    subtype: string
    stored: true
    filterable: true
    orderable: true

  - name: unexposed
    description: This attribute is not exposed.
    type: string
    stored: true
    filterable: true
    format: free
    orderable: true

# Relations
relations:
- rest_name: task
  get:
    description: yeye.
    parameters:
      entries:
      - name: ltgp1
        description: this is ltgp1.
        type: string
        example_value: ltgp1

      - name: ltgp2
        description: this is ltgp2.
        type: boolean
        example_value: "true"
  create:
    description: yoyo.
    parameters:
      entries:
      - name: ltcp1
        description: this is ltcp1.
        type: string
        example_value: ltcp1

      - name: ltcp2
        description: this is ltcp2.
        type: boolean
        example_value: "true"

- rest_name: user
  get:
    description: yeye.
`)
			})
		})
	})
}
