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

func toc(specs []spec.Specification) string {

	buf := &bytes.Buffer{}
	w := &tabwriter.Writer{}
	w.Init(buf, 0, 8, 0, ' ', 0)

	fmt.Fprintln(w, "| Object \t|\t Description \t|") // nolint: errcheck
	fmt.Fprintln(w, "| - \t|\t - \t|")                // nolint: errcheck

	for _, spec := range specs {

		model := spec.Model()

		if model.Private {
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
		buf.WriteString(fmt.Sprintf("#### `%s %s`\n\n", r.method, r.url))
		buf.WriteString(fmt.Sprintf(`%s`, r.doc))

		if r.params != nil {
			buf.WriteString("\n\n##### Parameters\n\n")
			for _, pd := range r.params.Entries {
				buf.WriteString(fmt.Sprintf("- `%s` (%s): %s\n", pd.Name, pd.Type, strings.Replace(pd.Description, "\n", "", -1)))
			}

			if r.params.Required != nil {
				buf.WriteString("\n\n##### Mandatory Parameters\n\n")

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

				buf.WriteString(fmt.Sprintf("%s", out))
			}
		}

		if i < len(full) {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func characteristics(attr *spec.Attribute) string {

	buf := &bytes.Buffer{}
	w := &tabwriter.Writer{}
	w.Init(buf, 0, 8, 0, ' ', 0)

	fmt.Fprintln(w, "")                                  // nolint: errcheck
	fmt.Fprintln(w, "| Characteristics \t|\t Value \t|") // nolint: errcheck
	fmt.Fprintln(w, "| - \t|\t -: \t|")                  // nolint: errcheck

	if attr.Identifier {
		fmt.Fprintln(w, "| Identifier \t|\t `true` \t|") // nolint: errcheck
	}

	if attr.AllowedChars != "" {
		fmt.Fprintln(w, fmt.Sprintf("| Format \t|\t `/%s/` \t|", attr.AllowedChars)) // nolint: errcheck
	}

	if len(attr.AllowedChoices) > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Allowed Value \t|\t `%s` \t|", strings.Join(attr.AllowedChoices, ", "))) // nolint: errcheck
	}

	if attr.DefaultValue != nil {
		fmt.Fprintln(w, fmt.Sprintf("| Default \t|\t `%#v` \t|", attr.DefaultValue)) // nolint: errcheck
	}

	if attr.MinLength > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Min length \t|\t `%d` \t|", attr.MinLength)) // nolint: errcheck
	}

	if attr.MaxLength > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Max length \t|\t `%d` \t|", attr.MaxLength)) // nolint: errcheck
	}

	if attr.MinValue > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Min length \t|\t `%v` \t|", attr.MinValue)) // nolint: errcheck
	}

	if attr.MaxValue > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Max length \t|\t `%v` \t|", attr.MaxValue)) // nolint: errcheck
	}

	if attr.Autogenerated {
		fmt.Fprintln(w, "| Autogenerated \t|\t `true` \t|") // nolint: errcheck
	}

	if attr.Required {
		fmt.Fprintln(w, "| Required \t|\t `true` \t|") // nolint: errcheck
	}

	if attr.ReadOnly {
		fmt.Fprintln(w, "| Read only \t|\t `true` \t|") // nolint: errcheck
	}

	if attr.CreationOnly {
		fmt.Fprintln(w, "| Creation only \t|\t `true` \t|") // nolint: errcheck
	}

	if attr.Orderable {
		fmt.Fprintln(w, "| Orderable \t|\t `true` \t|") // nolint: errcheck
	}

	if attr.Filterable {
		fmt.Fprintln(w, "| Filterable \t|\t `true` \t|") // nolint: errcheck
	}

	if attr.OmitEmpty {
		fmt.Fprintln(w, "| Omit if empty \t|\t `true` \t|") // nolint: errcheck
	}

	w.Flush() // nolint: errcheck

	str := buf.String()

	if len(strings.Split(str, "\n")) == 4 {
		return "\n"
	}

	return str + "\n"
}

func makeExample(s spec.Specification, version string) string {

	data := map[string]interface{}{}

	for _, attr := range s.Attributes(version) {

		if attr.DefaultValue != nil {
			continue
		}

		if attr.ExampleValue != nil {
			data[attr.Name] = attr.ExampleValue
		}
	}

	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	if len(data) == 0 {
		return ""
	}

	return strings.Replace(string(d), "\n", "\\n", -1)
}
