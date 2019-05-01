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
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec_LoadSpecificationDir(t *testing.T) {

	Convey("Given I load a spec folder", t, func() {

		set, err := LoadSpecificationSet(
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
			"test",
		)

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the specification set should have 4 entries", func() {
			So(set.Len(), ShouldEqual, 4)
		})

		Convey("Then calling Specifications should return the sorted list", func() {
			ss := set.Specifications()
			So(ss[0].Model().RestName, ShouldEqual, "list")
			So(ss[1].Model().RestName, ShouldEqual, "root")
			So(ss[2].Model().RestName, ShouldEqual, "task")
			So(ss[3].Model().RestName, ShouldEqual, "user")
		})

		Convey("Then the relationships should be correct", func() {
			rs := set.Relationships()
			So(rs["List"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieves the list with the given ID.",
					ParameterReferences: []string{
						"sharedParameterA",
						"sharedParameterB",
					},
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "lgp1",
								Description:  "this is lgp1.",
								Type:         ParameterTypeString,
								ExampleValue: "lgp1",
							},
							&Parameter{
								Name:         "lgp2",
								Description:  "this is lgp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
							&Parameter{
								Name:         "sAp1",
								Description:  "this is sAp1.",
								Type:         ParameterTypeString,
								ExampleValue: "sAp1",
							},
							&Parameter{
								Name:         "sAp2",
								Description:  "this is sAp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
							&Parameter{
								Name:         "sBp1",
								Description:  "this is sBp1.",
								Type:         ParameterTypeString,
								ExampleValue: "sBp1",
							},
							&Parameter{
								Name:         "sBp2",
								Description:  "this is sBp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["List"].Create, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "you.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rlcp1",
								Description:  "this is rlcp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rlcp1",
							},
							&Parameter{
								Name:         "rlcp2",
								Description:  "this is rlcp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["List"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the list with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "lup1",
								Description:  "this is lup1.",
								Type:         ParameterTypeString,
								ExampleValue: "lup1",
							},
							&Parameter{
								Name:         "lup2",
								Description:  "this is lup2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["List"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the list with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ldp1",
								Description:  "this is ldp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ldp1",
							},
							&Parameter{
								Name:         "ldp2",
								Description:  "this is ldp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["Task"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieve the task with the given ID.",
				},
			})
			So(rs["Task"].GetMany, ShouldResemble, map[string]*RelationAction{
				"list": &RelationAction{
					Description: "yeye.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ltgp1",
								Description:  "this is ltgp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ltgp1",
							},
							&Parameter{
								Name:         "ltgp2",
								Description:  "this is ltgp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["Task"].Create, ShouldResemble, map[string]*RelationAction{
				"list": &RelationAction{
					Description: "yoyo.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ltcp1",
								Description:  "this is ltcp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ltcp1",
							},
							&Parameter{
								Name:         "ltcp2",
								Description:  "this is ltcp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["Task"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the task with the given ID.",
				},
			})
			So(rs["Task"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the task with the given ID.",
				},
			})
			So(rs["User"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieves the user with the given ID.",
				},
			})
			So(rs["User"].GetMany, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "yey.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rugmp1",
								Description:  "this is rugmp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rugmp1",
							},
							&Parameter{
								Name:         "rugmp2",
								Description:  "this is rugmp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
				"list": &RelationAction{
					Description: "yeye.",
				},
			})
			So(rs["User"].Create, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "you.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rucp1",
								Description:  "this is rucp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rucp1",
							},
							&Parameter{
								Name:         "rucp2",
								Description:  "this is rucp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["User"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the user with the given ID.",
				},
			})
			So(rs["User"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the user with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Required: [][][]string{[][]string{[]string{"confirm"}}},
						Entries: []*Parameter{
							&Parameter{
								Name:        "confirm",
								Description: "this is required.",
								Type:        ParameterTypeBool,
							},
						},
					},
				},
			})
		})

		Convey("Then the relationships by rest name should be correct", func() {
			rs := set.RelationshipsByRestName()
			So(rs["list"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieves the list with the given ID.",
					ParameterReferences: []string{
						"sharedParameterA",
						"sharedParameterB",
					},
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "lgp1",
								Description:  "this is lgp1.",
								Type:         ParameterTypeString,
								ExampleValue: "lgp1",
							},
							&Parameter{
								Name:         "lgp2",
								Description:  "this is lgp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
							&Parameter{
								Name:         "sAp1",
								Description:  "this is sAp1.",
								Type:         ParameterTypeString,
								ExampleValue: "sAp1",
							},
							&Parameter{
								Name:         "sAp2",
								Description:  "this is sAp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
							&Parameter{
								Name:         "sBp1",
								Description:  "this is sBp1.",
								Type:         ParameterTypeString,
								ExampleValue: "sBp1",
							},
							&Parameter{
								Name:         "sBp2",
								Description:  "this is sBp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["list"].Create, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "you.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rlcp1",
								Description:  "this is rlcp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rlcp1",
							},
							&Parameter{
								Name:         "rlcp2",
								Description:  "this is rlcp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["list"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the list with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "lup1",
								Description:  "this is lup1.",
								Type:         ParameterTypeString,
								ExampleValue: "lup1",
							},
							&Parameter{
								Name:         "lup2",
								Description:  "this is lup2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["list"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the list with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ldp1",
								Description:  "this is ldp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ldp1",
							},
							&Parameter{
								Name:         "ldp2",
								Description:  "this is ldp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["task"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieve the task with the given ID.",
				},
			})
			So(rs["task"].GetMany, ShouldResemble, map[string]*RelationAction{
				"list": &RelationAction{
					Description: "yeye.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ltgp1",
								Description:  "this is ltgp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ltgp1",
							},
							&Parameter{
								Name:         "ltgp2",
								Description:  "this is ltgp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["task"].Create, ShouldResemble, map[string]*RelationAction{
				"list": &RelationAction{
					Description: "yoyo.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ltcp1",
								Description:  "this is ltcp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ltcp1",
							},
							&Parameter{
								Name:         "ltcp2",
								Description:  "this is ltcp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["task"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the task with the given ID.",
				},
			})
			So(rs["task"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the task with the given ID.",
				},
			})
			So(rs["user"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieves the user with the given ID.",
				},
			})
			So(rs["user"].GetMany, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "yey.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rugmp1",
								Description:  "this is rugmp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rugmp1",
							},
							&Parameter{
								Name:         "rugmp2",
								Description:  "this is rugmp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
				"list": &RelationAction{
					Description: "yeye.",
				},
			})
			So(rs["user"].Create, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "you.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rucp1",
								Description:  "this is rucp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rucp1",
							},
							&Parameter{
								Name:         "rucp2",
								Description:  "this is rucp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["user"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the user with the given ID.",
				},
			})
			So(rs["user"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the user with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Required: [][][]string{[][]string{[]string{"confirm"}}},
						Entries: []*Parameter{
							&Parameter{
								Name:        "confirm",
								Description: "this is required.",
								Type:        ParameterTypeBool,
							},
						},
					},
				},
			})
		})

		Convey("Then the relationships by resource name should be correct", func() {
			rs := set.RelationshipsByResourceName()
			So(rs["lists"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieves the list with the given ID.",
					ParameterReferences: []string{
						"sharedParameterA",
						"sharedParameterB",
					},
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "lgp1",
								Description:  "this is lgp1.",
								Type:         ParameterTypeString,
								ExampleValue: "lgp1",
							},
							&Parameter{
								Name:         "lgp2",
								Description:  "this is lgp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
							&Parameter{
								Name:         "sAp1",
								Description:  "this is sAp1.",
								Type:         ParameterTypeString,
								ExampleValue: "sAp1",
							},
							&Parameter{
								Name:         "sAp2",
								Description:  "this is sAp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
							&Parameter{
								Name:         "sBp1",
								Description:  "this is sBp1.",
								Type:         ParameterTypeString,
								ExampleValue: "sBp1",
							},
							&Parameter{
								Name:         "sBp2",
								Description:  "this is sBp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["lists"].Create, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "you.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rlcp1",
								Description:  "this is rlcp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rlcp1",
							},
							&Parameter{
								Name:         "rlcp2",
								Description:  "this is rlcp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["lists"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the list with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "lup1",
								Description:  "this is lup1.",
								Type:         ParameterTypeString,
								ExampleValue: "lup1",
							},
							&Parameter{
								Name:         "lup2",
								Description:  "this is lup2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["lists"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the list with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ldp1",
								Description:  "this is ldp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ldp1",
							},
							&Parameter{
								Name:         "ldp2",
								Description:  "this is ldp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["tasks"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieve the task with the given ID.",
				},
			})
			So(rs["tasks"].GetMany, ShouldResemble, map[string]*RelationAction{
				"list": &RelationAction{
					Description: "yeye.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ltgp1",
								Description:  "this is ltgp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ltgp1",
							},
							&Parameter{
								Name:         "ltgp2",
								Description:  "this is ltgp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["tasks"].Create, ShouldResemble, map[string]*RelationAction{
				"list": &RelationAction{
					Description: "yoyo.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "ltcp1",
								Description:  "this is ltcp1.",
								Type:         ParameterTypeString,
								ExampleValue: "ltcp1",
							},
							&Parameter{
								Name:         "ltcp2",
								Description:  "this is ltcp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["tasks"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the task with the given ID.",
				},
			})
			So(rs["tasks"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the task with the given ID.",
				},
			})
			So(rs["users"].Get, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Retrieves the user with the given ID.",
				},
			})
			So(rs["users"].GetMany, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "yey.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rugmp1",
								Description:  "this is rugmp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rugmp1",
							},
							&Parameter{
								Name:         "rugmp2",
								Description:  "this is rugmp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
				"list": &RelationAction{
					Description: "yeye.",
				},
			})
			So(rs["users"].Create, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "you.",
					ParameterDefinition: &ParameterDefinition{
						Entries: []*Parameter{
							&Parameter{
								Name:         "rucp1",
								Description:  "this is rucp1.",
								Type:         ParameterTypeString,
								ExampleValue: "rucp1",
							},
							&Parameter{
								Name:         "rucp2",
								Description:  "this is rucp2.",
								Type:         ParameterTypeBool,
								ExampleValue: "true",
							},
						},
					},
				},
			})
			So(rs["users"].Update, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Updates the user with the given ID.",
				},
			})
			So(rs["users"].Delete, ShouldResemble, map[string]*RelationAction{
				"root": &RelationAction{
					Description: "Deletes the user with the given ID.",
					ParameterDefinition: &ParameterDefinition{
						Required: [][][]string{[][]string{[]string{"confirm"}}},
						Entries: []*Parameter{
							&Parameter{
								Name:        "confirm",
								Description: "this is required.",
								Type:        ParameterTypeBool,
							},
						},
					},
				},
			})
		})

		Convey("Then the specification set should be correct", func() {
			So(len(set.Specification("task").Attributes("v1")), ShouldEqual, 6)
			So(len(set.Specification("root").Attributes("v1")), ShouldEqual, 0)
			So(len(set.Specification("list").Attributes("v1")), ShouldEqual, 11)
			So(len(set.Specification("user").Attributes("v1")), ShouldEqual, 6)
		})

		Convey("Then the API linking should be correct", func() {
			So(set.Specification("root").Relation("user").Specification(), ShouldEqual, set.Specification("user"))
			So(set.Specification("root").Relation("list").Specification(), ShouldEqual, set.Specification("list"))
			So(set.Specification("list").Relation("task").Specification(), ShouldEqual, set.Specification("task"))
		})

		Convey("Then the config should be correctly loaded", func() {
			So(set.Configuration().Name, ShouldEqual, "testmodel")
		})

		Convey("Then the type mapping should be correctly loaded", func() {
			m, _ := set.TypeMapping().Mapping("test", "string_map")
			So(m.Type, ShouldEqual, "map[string]string")
		})

		Convey("Then the validation mapping should be correctly loaded", func() {
			m, _ := set.ValidationMapping().Mapping("test", "$username")
			So(m.Name, ShouldEqual, "validate.CheckUserName")
		})

		Convey("Then the api info should be correctly loaded", func() {

			So(set.APIInfo().Version, ShouldEqual, 1)
		})
	})
}
