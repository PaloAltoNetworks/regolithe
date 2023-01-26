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
)

func TestParameter_Validate(t *testing.T) {
	type fields struct {
		Name           string
		Description    string
		Type           ParameterType
		Multiple       bool
		AllowedChoices []string
		DefaultValue   any
		ExampleValue   any
	}
	type args struct {
		relatedReSTName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []error
	}{
		{
			"missing description",
			fields{
				Name:        "p",
				Description: "",
				Type:        ParameterTypeInt,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: description of parameter 'p' must end with a period"),
			},
		},
		{
			"description with not period",
			fields{
				Name:        "p",
				Description: "desc",
				Type:        ParameterTypeInt,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: description of parameter 'p' must end with a period"),
			},
		},
		{
			"parameter with no type",
			fields{
				Name:        "p",
				Description: "desc.",
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: type of parameter 'p' must be set"),
			},
		},
		{
			"parameter with unknown type",
			fields{
				Name:        "p",
				Description: "desc.",
				Type:        ParameterType("nope"),
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: type of parameter 'p' must be 'string', 'integer', 'float', 'boolean', 'enum', 'time' or 'duration'"),
			},
		},
		{
			"enum with no allowed_choices",
			fields{
				Name:        "p",
				Description: "desc.",
				Type:        ParameterTypeEnum,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: enum parameter 'p' must define allowed_choices"),
			},
		},
		{
			"not enum with allowed_choices",
			fields{
				Name:           "p",
				Description:    "desc.",
				Type:           ParameterTypeString,
				AllowedChoices: []string{"hello"},
				ExampleValue:   1,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is not an enum but defines allowed_choices"),
			},
		},
		{
			"string with no default_value, no example_value",
			fields{
				Name:        "p",
				Description: "desc.",
				Type:        ParameterTypeString,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' must provide an example value as it doesn't have a default"),
			},
		},
		{
			"string non string default value",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeString,
				DefaultValue: 1,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an string, but the default value is not"),
			},
		},
		{
			"enum non string default value",
			fields{
				Name:           "p",
				Description:    "desc.",
				Type:           ParameterTypeEnum,
				AllowedChoices: []string{"a"},
				DefaultValue:   1,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an enum, but the default value is not"),
			},
		},
		{
			"int non int default value",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeInt,
				DefaultValue: "1",
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an integer, but the default value is not"),
			},
		},
		{
			"float non float default value",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeFloat,
				DefaultValue: "1",
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an float, but the default value is not"),
			},
		},
		{
			"bool non bool default value",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeBool,
				DefaultValue: "1",
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an boolean, but the default value is not"),
			},
		},
		{
			"duration non duration default value",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeDuration,
				DefaultValue: "1",
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an duration, but the default value is not"),
			},
		},
		{
			"duration non duration default value 2",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeDuration,
				DefaultValue: true,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an duration, but the default value is not"),
			},
		},
		{
			"time non time default value",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeTime,
				DefaultValue: "1",
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an time, but the default value is not"),
			},
		},
		{
			"time non time default value 2",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeTime,
				DefaultValue: true,
			},
			args{
				"spec",
			},
			[]error{
				fmt.Errorf("spec.spec: parameter 'p' is defined as an time, but the default value is not"),
			},
		},
		{
			"valid string param",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeString,
				DefaultValue: "hello",
			},
			args{
				"spec",
			},
			nil,
		},
		{
			"valid enum param",
			fields{
				Name:           "p",
				Description:    "desc.",
				Type:           ParameterTypeEnum,
				AllowedChoices: []string{"hello", "bye"},
				DefaultValue:   "hello",
			},
			args{
				"spec",
			},
			nil,
		},
		{
			"valid int param",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeInt,
				DefaultValue: 42,
			},
			args{
				"spec",
			},
			nil,
		},
		{
			"valid float param",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeFloat,
				DefaultValue: 42.42,
			},
			args{
				"spec",
			},
			nil,
		},
		{
			"valid bool param",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeBool,
				DefaultValue: true,
			},
			args{
				"spec",
			},
			nil,
		},
		{
			"valid duration param",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeDuration,
				DefaultValue: "3h",
			},
			args{
				"spec",
			},
			nil,
		},
		{
			"valid time param",
			fields{
				Name:         "p",
				Description:  "desc.",
				Type:         ParameterTypeTime,
				DefaultValue: "Mon, 02 Jan 2006 15:04:05 -0700",
			},
			args{
				"spec",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameter{
				Name:           tt.fields.Name,
				Description:    tt.fields.Description,
				Type:           tt.fields.Type,
				Multiple:       tt.fields.Multiple,
				AllowedChoices: tt.fields.AllowedChoices,
				DefaultValue:   tt.fields.DefaultValue,
				ExampleValue:   tt.fields.ExampleValue,
			}
			got := p.Validate(tt.args.relatedReSTName)
			switch {
			case (len(tt.want) == 0 && len(got) != 0) || (len(tt.want) != 0 && len(got) == 0):
				t.Errorf("Parameter.Validate() = %v, want %v", got, tt.want)
			case len(tt.want) == 0 && len(got) == 0:
			case got[0].Error() != tt.want[0].Error():
				t.Errorf("Parameter.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
