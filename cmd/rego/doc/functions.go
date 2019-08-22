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

package doc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"go.aporeto.io/regolithe/spec"
)

const rootSpecRestName = "root"

func typeOf(attr *spec.Attribute) string {

	switch attr.Type {
	case spec.AttributeTypeExt:
		return "`" + attr.SubType + "`"
	case spec.AttributeTypeList:
		return "`[]" + attr.SubType + "`"
	case spec.AttributeTypeEnum:
		return "`emum(" + strings.Join(attr.AllowedChoices, " | ") + ")`"
	case spec.AttributeTypeRef:
		return fmt.Sprintf("[`%s`](#%s)", attr.SubType, attr.SubType)
	case spec.AttributeTypeRefList:
		return fmt.Sprintf("[`[]%s`](#%s)", attr.SubType, attr.SubType)
	case spec.AttributeTypeRefMap:
		return fmt.Sprintf("[`map[string]%s`](#%s)", attr.SubType, attr.SubType)
	default:
		return "`" + string(attr.Type) + "`"
	}
}

func toc(specs []spec.Specification) string {

	buf := &bytes.Buffer{}
	w := &tabwriter.Writer{}
	w.Init(buf, 0, 8, 0, ' ', 0)

	fmt.Fprintln(w, "| Resource \t|\t Description \t|") // nolint: errcheck
	fmt.Fprintln(w, "| - \t|\t - \t|")                  // nolint: errcheck

	for _, spec := range specs {

		model := spec.Model()

		if model.Group == "none" {
			continue
		}

		_, ok := model.Extensions["forceDocumentation"]
		if model.Private && !ok {
			continue
		}

		fmt.Fprintln(
			w,
			fmt.Sprintf(
				"| [%s](#%s) \t|\t %s \t|",
				model.EntityName,
				strings.ToLower(model.EntityName),
				shortString(model.Description, 80),
			),
		) // nolint: errcheck
	}
	w.Flush() // nolint: errcheck

	return buf.String()
}

type operation struct {
	method string
	url    string
	doc    string
	params *spec.ParameterDefinition
}

func (o operation) String() string {
	return fmt.Sprintf("| `%s` \t|\t `%s` \t|\t %s \t|", o.method, o.url, o.doc)
}

func operations(spec spec.Specification, relationships map[string]*spec.Relationship, set spec.SpecificationSet) string {

	var rootOps []operation
	var parentOps []operation // nolint
	var childOps []operation

	model := spec.Model()

	if model.Update != nil {
		rootOps = append(rootOps, operation{
			method: "PUT",
			url:    fmt.Sprintf("/%s/:id", model.ResourceName),
			doc:    model.Update.Description,
			params: model.Update.ParameterDefinition,
		})
	}

	if model.Delete != nil {
		rootOps = append(rootOps, operation{
			method: "DELETE",
			url:    fmt.Sprintf("/%s/:id", model.ResourceName),
			doc:    model.Delete.Description,
			params: model.Delete.ParameterDefinition,
		})
	}

	if model.Get != nil {
		rootOps = append(rootOps, operation{
			method: "GET",
			url:    fmt.Sprintf("/%s/:id", model.ResourceName),
			doc:    model.Get.Description,
			params: model.Get.ParameterDefinition,
		})
	}

	for k, ra := range relationships[model.RestName].GetMany {
		if k == rootSpecRestName {
			rootOps = append(rootOps, operation{
				method: "GET",
				url:    fmt.Sprintf("/%s", model.ResourceName),
				doc:    ra.Description,
				params: ra.ParameterDefinition,
			})
			continue
		}
		childSpec := set.Specification(k)
		parentOps = append(parentOps, operation{
			method: "GET",
			url:    fmt.Sprintf("/%s/:id/%s", childSpec.Model().ResourceName, model.ResourceName),
			doc:    ra.Description,
			params: ra.ParameterDefinition,
		})
	}

	for k, ra := range relationships[model.RestName].Create {
		if k == rootSpecRestName {
			rootOps = append(rootOps, operation{
				method: "POST",
				url:    fmt.Sprintf("/%s", model.ResourceName),
				doc:    ra.Description,
				params: ra.ParameterDefinition,
			})
			continue
		}
		childSpec := set.Specification(k)
		parentOps = append(parentOps, operation{
			method: "POST",
			url:    fmt.Sprintf("/%s/:id/%s", childSpec.Model().ResourceName, model.ResourceName),
			doc:    ra.Description,
			params: ra.ParameterDefinition,
		})
	}

	for _, rel := range spec.Relations() {

		childSpec := set.Specification(rel.RestName)
		childModel := childSpec.Model()

		if rel.Create != nil {
			childOps = append(childOps, operation{
				method: "POST",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Create.Description,
				params: rel.Create.ParameterDefinition,
			})
		}
		if rel.Update != nil {
			childOps = append(childOps, operation{
				method: "PUT",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Update.Description,
				params: rel.Update.ParameterDefinition,
			})
		}
		if rel.Get != nil {
			childOps = append(childOps, operation{
				method: "GET",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Get.Description,
				params: rel.Get.ParameterDefinition,
			})
		}
		if rel.Delete != nil {
			childOps = append(childOps, operation{
				method: "DELETE",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Delete.Description,
				params: rel.Delete.ParameterDefinition,
			})
		}
	}

	cmp := func(a []operation) func(i int, j int) bool {
		return func(i int, j int) bool {
			c := strings.Compare(a[i].url, a[j].url)
			if c == 0 {
				return strings.Compare(a[i].method, a[j].method) == -1
			}
			return c == -1
		}
	}

	sort.Slice(rootOps, cmp(rootOps))
	sort.Slice(parentOps, cmp(parentOps))
	sort.Slice(childOps, cmp(childOps))
	full := append(rootOps, append(parentOps, childOps...)...)
	total := len(full)

	if total == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	for i, r := range full {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("##### `%s %s`\n\n", r.method, r.url))
		buf.WriteString(r.doc)

		if r.params != nil {
			buf.WriteString("\n\nParameters:\n\n")
			for _, pd := range r.params.Entries {

				var enumValues string
				if pd.Type == "enum" {
					enumValues = "(" + strings.Join(pd.AllowedChoices, " | ") + ")"
				}
				buf.WriteString(fmt.Sprintf("- `%s` (`%s%s`): %s\n", pd.Name, pd.Type, enumValues, strings.Replace(pd.Description, "\n", "", -1)))
			}

			if r.params.Required != nil {
				buf.WriteString("\n\nMandatory Parameters\n\n")

				var out string
				for i, lvl1 := range r.params.Required {
					if len(r.params.Required) > 1 {
						out += "("
					}
					for j, lvl2 := range lvl1 {
						if len(lvl1) > 1 {
							out += "("
						}
						out += "`" + strings.Join(lvl2, "` and `") + "`"
						if len(lvl1) > 1 {
							out += ")"
						}
						if j+1 != len(lvl1) {
							out += " or "
						}
					}
					if len(r.params.Required) > 1 {
						out += ")"
					}
					if i+1 != len(r.params.Required) {
						out += " and "
					}
				}

				buf.WriteString(out)
			}
		}

		if i < len(full) {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func characteristics(attr *spec.Attribute) string {

	var out []string

	if attr.Identifier {
		out = append(out, "`identifier`")
	}

	if attr.Autogenerated {
		out = append(out, "`autogenerated`")
	}

	if attr.Required {
		out = append(out, "`required`")
	}

	if attr.ReadOnly {
		out = append(out, "`read_only`")
	}

	if attr.CreationOnly {
		out = append(out, "`creation_only`")
	}

	if attr.AllowedChars != "" {
		out = append(out, fmt.Sprintf("`format=%s`", attr.AllowedChars))
	}

	if attr.MinLength > 0 {
		out = append(out, fmt.Sprintf("`min_length=%d`", attr.MinLength))
	}

	if attr.MaxLength > 0 {
		out = append(out, fmt.Sprintf("`max_length=%d`", attr.MaxLength))
	}

	if attr.MinValue > 0 {
		out = append(out, fmt.Sprintf("`min_value=%f`", attr.MinValue))
	}

	if attr.MaxValue > 0 {
		out = append(out, fmt.Sprintf("`max_value=%f`", attr.MaxValue))
	}

	if len(out) == 0 {
		return ""
	}

	return " [" + strings.Join(out, ",") + "]"
}

func makeDefaultValue(attr *spec.Attribute) string {

	if attr.DefaultValue == nil {
		return ""
	}

	dv, err := json.MarshalIndent(attr.DefaultValue, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("\n\nDefault value:\n\n```json\n%s\n```\n", string(dv))
}

func makeExample(s spec.Specification, version string) string {

	data := map[string]interface{}{}

	for _, attr := range s.Attributes(version) {

		if !attr.Exposed {
			continue
		}

		if attr.ExampleValue != nil {
			data[attr.Name] = attr.ExampleValue
			continue
		}

		if attr.Autogenerated {
			continue
		}

		if attr.DefaultValue != nil {
			data[attr.Name] = attr.DefaultValue
			continue
		}

		if attr.Type == spec.AttributeTypeEnum {
			data[attr.Name] = attr.AllowedChoices[0]
		}

		if attr.Type == spec.AttributeTypeBool {
			data[attr.Name] = false
		}
	}

	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		return ""
	}

	return strings.Replace(string(d), "\n", `\n`, -1)
}
