package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/aporeto-inc/regolithe/spec"
)

func toc(specs []*spec.Specification) string {

	buf := &bytes.Buffer{}
	w := &tabwriter.Writer{}
	w.Init(buf, 0, 8, 0, ' ', 0)

	fmt.Fprintln(w, "| Object \t|\t Description \t|")
	fmt.Fprintln(w, "| - \t|\t - \t|")

	for _, spec := range specs {
		if spec.Model.Private {
			continue
		}

		fmt.Fprintln(
			w,
			fmt.Sprintf(
				"| [%s](#%s) \t|\t %s \t|",
				spec.Model.EntityName,
				strings.ToLower(spec.Model.EntityName),
				shortString(spec.Model.Description, 80),
			),
		)
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

func operations(spec *spec.Specification, relationships map[string]*spec.Relationship, set *spec.SpecificationSet) string {

	var rootOps []operation
	var parentOps []operation
	var childOps []operation

	if spec.Model.AllowsUpdate {
		rootOps = append(rootOps, operation{
			method: "PUT",
			url:    fmt.Sprintf("/%s/:id", spec.Model.ResourceName),
			doc:    fmt.Sprintf("Updates the `%s` with the given `:id`.", spec.Model.RestName),
		})
	}

	if spec.Model.AllowsDelete {
		rootOps = append(rootOps, operation{
			method: "DELETE",
			url:    fmt.Sprintf("/%s/:id", spec.Model.ResourceName),
			doc:    fmt.Sprintf("Deletes the `%s` with the given `:id`.", spec.Model.RestName),
		})
	}

	if spec.Model.AllowsGet {
		rootOps = append(rootOps, operation{
			method: "GET",
			url:    fmt.Sprintf("/%s/:id", spec.Model.ResourceName),
			doc:    fmt.Sprintf("Retrieve the `%s` with the given `:id`.", spec.Model.RestName),
		})
	}

	for k := range relationships[spec.Model.RestName].AllowsGetMany {
		if k == "root" {
			if k == "root" {
				childSpec := set.Specification("root")
				rootOps = append(rootOps, operation{
					method: "GET",
					url:    fmt.Sprintf("/%s", spec.Model.ResourceName),
					doc:    childSpec.Relation(spec.Model.RestName).Descriptions["get"],
				})
				continue
			}
		}
		childSpec := set.Specification(k)
		parentOps = append(parentOps, operation{
			method: "GET",
			url:    fmt.Sprintf("/%s/:id/%s", childSpec.Model.ResourceName, spec.Model.ResourceName),
			doc:    childSpec.Relation(spec.Model.RestName).Descriptions["get"],
		})
	}

	for k := range relationships[spec.Model.RestName].AllowsCreate {
		if k == "root" {
			if k == "root" {
				childSpec := set.Specification("root")
				rootOps = append(rootOps, operation{
					method: "POST",
					url:    fmt.Sprintf("/%s", spec.Model.ResourceName),
					doc:    childSpec.Relation(spec.Model.RestName).Descriptions["create"],
				})
				continue
			}
		}
		childSpec := set.Specification(k)
		parentOps = append(parentOps, operation{
			method: "POST",
			url:    fmt.Sprintf("/%s/:id/%s", childSpec.Model.ResourceName, spec.Model.ResourceName),
			doc:    childSpec.Relation(spec.Model.RestName).Descriptions["create"],
		})
	}

	for _, rel := range spec.Relations {
		childSpec := set.Specification(rel.RestName)

		if rel.AllowsCreate {
			childOps = append(childOps, operation{
				method: "POST",
				url:    fmt.Sprintf("/%s/:id/%s", spec.Model.ResourceName, childSpec.Model.ResourceName),
				doc:    rel.Descriptions["create"],
			})
		}
		if rel.AllowsUpdate {
			childOps = append(childOps, operation{
				method: "PUT",
				url:    fmt.Sprintf("/%s/:id/%s", spec.Model.ResourceName, childSpec.Model.ResourceName),
				doc:    rel.Descriptions["update"],
			})
		}
		if rel.AllowsGet {
			childOps = append(childOps, operation{
				method: "GET",
				url:    fmt.Sprintf("/%s/:id/%s", spec.Model.ResourceName, childSpec.Model.ResourceName),
				doc:    rel.Descriptions["get"],
			})
		}
		if rel.AllowsDelete {
			childOps = append(childOps, operation{
				method: "DELETE",
				url:    fmt.Sprintf("/%s/:id/%s", spec.Model.ResourceName, childSpec.Model.ResourceName),
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

	fmt.Fprintln(w, "| Method \t|\t URL \t|\t Description \t|")
	fmt.Fprintln(w, "| -: \t|\t - \t|\t - \t|")

	for i := 0; i < total; i++ {
		fmt.Fprintln(w, full[i].String())
	}
	w.Flush() // nolint: errcheck

	return buf.String()
}

func characteristics(attr *spec.Attribute) string {

	buf := &bytes.Buffer{}
	w := &tabwriter.Writer{}
	w.Init(buf, 0, 8, 0, ' ', 0)

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "| Characteristics \t|\t Value \t|")
	fmt.Fprintln(w, "| - \t|\t -: \t|")

	if attr.Identifier {
		fmt.Fprintln(w, "| Identifier \t|\t `true` \t|")
	}

	if attr.AllowedChars != "" {
		fmt.Fprintln(w, fmt.Sprintf("| Format \t|\t `/%s/` \t|", attr.AllowedChars))
	}

	if len(attr.AllowedChoices) > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Allowed Value \t|\t `%s` \t|", strings.Join(attr.AllowedChoices, ", ")))
	}

	if attr.DefaultValue != nil {
		fmt.Fprintln(w, fmt.Sprintf("| Default \t|\t `%s` \t|", attr.DefaultValue))
	}

	if attr.MinLength > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Min length \t|\t `%d` \t|", attr.MinLength))
	}

	if attr.MaxLength > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Max length \t|\t `%d` \t|", attr.MaxLength))
	}

	if attr.MinValue > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Min length \t|\t `%v` \t|", attr.MinValue))
	}

	if attr.MaxValue > 0 {
		fmt.Fprintln(w, fmt.Sprintf("| Max length \t|\t `%v` \t|", attr.MaxValue))
	}

	if attr.Autogenerated {
		fmt.Fprintln(w, "| Autogenerated \t|\t `true` \t|")
	}

	if attr.Required {
		fmt.Fprintln(w, "| Required \t|\t `true` \t|")
	}

	if attr.ReadOnly {
		fmt.Fprintln(w, "| Read only \t|\t `true` \t|")
	}

	if attr.CreationOnly {
		fmt.Fprintln(w, "| Creation only \t|\t `true` \t|")
	}

	if attr.Orderable {
		fmt.Fprintln(w, "| Orderable \t|\t `true` \t|")
	}

	if attr.Filterable {
		fmt.Fprintln(w, "| Filterable \t|\t `true` \t|")
	}

	w.Flush() // nolint: errcheck

	str := buf.String()

	if len(strings.Split(str, "\n")) == 4 {
		return "\n"
	}

	return str + "\n"
}
