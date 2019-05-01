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
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	yaml "gopkg.in/yaml.v2"
)

func TestModel_orderedKeys(t *testing.T) {

	Convey("Given I have a mode", t, func() {

		ra := &RelationAction{}
		m := &Model{
			Aliases:      []string{"alias1", "alias2"},
			Delete:       ra,
			Get:          ra,
			Update:       ra,
			Description:  "desc",
			EntityName:   "entity name",
			Extends:      []string{"ext1", "ext2"},
			IsRoot:       true,
			ResourceName: "resource name",
			RestName:     "rest name",
			Private:      true,

			EntityNamePlural: "plural",
		}

		Convey("When I call toYAMLMapSlice ", func() {

			ms := toYAMLMapSlice(m)

			Convey("Then the result should be correct", func() {
				So(len(ms), ShouldEqual, 11)

				// ShouldResemble is weird and doesn't work here.
				So(fmt.Sprintf("%#v", ms[0]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "rest_name", Value: "rest name"}))
				So(fmt.Sprintf("%#v", ms[1]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "resource_name", Value: "resource name"}))
				So(fmt.Sprintf("%#v", ms[2]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "entity_name", Value: "entity name"}))
				So(fmt.Sprintf("%#v", ms[3]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "description", Value: "desc"}))
				So(fmt.Sprintf("%#v", ms[4]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "aliases", Value: []string{"alias1", "alias2"}}))
				So(fmt.Sprintf("%#v", ms[5]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "private", Value: true}))
				So(fmt.Sprintf("%#v", ms[6]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "get", Value: ra}))
				So(fmt.Sprintf("%#v", ms[7]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "update", Value: ra}))
				So(fmt.Sprintf("%#v", ms[8]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "delete", Value: ra}))
				So(fmt.Sprintf("%#v", ms[9]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "extends", Value: []string{"ext1", "ext2"}}))
				So(fmt.Sprintf("%#v", ms[10]), ShouldResemble, fmt.Sprintf("%#v", yaml.MapItem{Key: "root", Value: true}))
			})
		})
	})
}

func Test_splitTags(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			"empty string",
			args{
				"",
			},
			"",
			false,
		},
		{
			"skip string",
			args{
				"-",
			},
			"",
			false,
		},
		{
			"normal string",
			args{
				"aa",
			},
			"aa",
			false,
		},
		{
			"omit empty string",
			args{
				"aa,omitempty",
			},
			"aa",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitTags(tt.args.tag)
			if got != tt.want {
				t.Errorf("splitTags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
