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
}

func (o operation) String() string {
	return fmt.Sprintf("| `%s` \t|\t `%s` \t|\t %s \t|", o.method, o.url, o.doc)
}

func operations(spec spec.Specification, relationships map[string]*spec.Relationship, set spec.SpecificationSet) string {

	var rootOps []operation
	var parentOps []operation
	var childOps []operation

	model := spec.Model()

	if model.AllowsUpdate {
		rootOps = append(rootOps, operation{
			method: "PUT",
			url:    fmt.Sprintf("/%s/:id", model.ResourceName),
			doc:    fmt.Sprintf("Updates the `%s` with the given `:id`.", model.RestName),
		})
	}

	if model.AllowsDelete {
		rootOps = append(rootOps, operation{
			method: "DELETE",
			url:    fmt.Sprintf("/%s/:id", model.ResourceName),
			doc:    fmt.Sprintf("Deletes the `%s` with the given `:id`.", model.RestName),
		})
	}

	if model.AllowsGet {
		rootOps = append(rootOps, operation{
			method: "GET",
			url:    fmt.Sprintf("/%s/:id", model.ResourceName),
			doc:    fmt.Sprintf("Retrieve the `%s` with the given `:id`.", model.RestName),
		})
	}

	for k := range relationships[model.RestName].AllowsGetMany {
		if k == rootSpecRestName {
			if k == rootSpecRestName {
				childSpec := set.Specification(rootSpecRestName)
				rootOps = append(rootOps, operation{
					method: "GET",
					url:    fmt.Sprintf("/%s", model.ResourceName),
					doc:    childSpec.Relation(model.RestName).Descriptions["get"],
				})
				continue
			}
		}
		childSpec := set.Specification(k)
		parentOps = append(parentOps, operation{
			method: "GET",
			url:    fmt.Sprintf("/%s/:id/%s", childSpec.Model().ResourceName, model.ResourceName),
			doc:    childSpec.Relation(model.RestName).Descriptions["get"],
		})
	}

	for k := range relationships[model.RestName].AllowsCreate {
		if k == rootSpecRestName {
			if k == rootSpecRestName {
				childSpec := set.Specification(rootSpecRestName)
				rootOps = append(rootOps, operation{
					method: "POST",
					url:    fmt.Sprintf("/%s", model.ResourceName),
					doc:    childSpec.Relation(model.RestName).Descriptions["create"],
				})
				continue
			}
		}
		childSpec := set.Specification(k)
		parentOps = append(parentOps, operation{
			method: "POST",
			url:    fmt.Sprintf("/%s/:id/%s", childSpec.Model().ResourceName, model.ResourceName),
			doc:    childSpec.Relation(model.RestName).Descriptions["create"],
		})
	}

	for _, rel := range spec.Relations() {

		childSpec := set.Specification(rel.RestName)
		childModel := childSpec.Model()

		if rel.AllowsCreate {
			childOps = append(childOps, operation{
				method: "POST",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Descriptions["create"],
			})
		}
		if rel.AllowsUpdate {
			childOps = append(childOps, operation{
				method: "PUT",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Descriptions["update"],
			})
		}
		if rel.AllowsGet {
			childOps = append(childOps, operation{
				method: "GET",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Descriptions["get"],
			})
		}
		if rel.AllowsDelete {
			childOps = append(childOps, operation{
				method: "DELETE",
				url:    fmt.Sprintf("/%s/:id/%s", model.ResourceName, childModel.ResourceName),
				doc:    rel.Descriptions["delete"],
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
	w := &tabwriter.Writer{}
	w.Init(buf, 0, 8, 0, ' ', 0)

	fmt.Fprintln(w, "| Method \t|\t URL \t|\t Description \t|") // nolint: errcheck
	fmt.Fprintln(w, "| -: \t|\t - \t|\t - \t|")                 // nolint: errcheck

	for i := 0; i < total; i++ {
		fmt.Fprintln(w, full[i].String()) // nolint: errcheck
	}
	w.Flush() // nolint: errcheck

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
		fmt.Fprintln(w, fmt.Sprintf("| Default \t|\t `%s` \t|", attr.DefaultValue)) // nolint: errcheck
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
