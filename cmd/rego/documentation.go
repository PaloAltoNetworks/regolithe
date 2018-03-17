package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/aporeto-inc/regolithe/cmd/rego/static"
	"github.com/aporeto-inc/regolithe/spec"
)

func shortString(s string, max int) string {
	if len(s) < max {
		return s
	}
	return s[:max] + "..."
}

var functions = template.FuncMap{
	"join":      strings.Join,
	"title":     strings.Title,
	"trimspace": strings.TrimSpace,
	"toc": func(specs []*spec.Specification) string {

		buf := &bytes.Buffer{}
		w := &tabwriter.Writer{}
		w.Init(buf, 0, 8, 0, ' ', 0)

		fmt.Fprintln(w, "| \t|\t \t|")
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
	},
	"characteristics": func(attr *spec.Attribute) string {

		buf := &bytes.Buffer{}
		w := &tabwriter.Writer{}
		w.Init(buf, 0, 8, 0, ' ', 0)

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

		return buf.String()
	},
}

func writeDoc(set *spec.SpecificationSet, format string, outFolder string) error {

	switch format {
	case "markdown":
		return writeMarkdownDoc(set, outFolder)
	default:
		return fmt.Errorf("Unsuported documentation format: %s", format)
	}
}

func writeMarkdownDoc(set *spec.SpecificationSet, outFolder string) error {

	data, err := static.Asset("templates/toc-md.gotpl")
	if err != nil {
		return fmt.Errorf("cannot open toc template: %s", err)
	}

	temp, err := template.New("toc").Funcs(functions).Parse(string(data))
	if err != nil {
		return fmt.Errorf("cannot parse template: %s", err)
	}

	buf := &bytes.Buffer{}
	if err := temp.Execute(buf, set); err != nil {
		return fmt.Errorf("cannot execute template: %s", err)
	}

	fmt.Println(buf.String())

	var skipped int
	specs := set.Specifications()
	for i, s := range specs {

		if s.Model.Private {
			skipped++
			continue
		}

		data, err := static.Asset("templates/spec-md.gotpl")
		if err != nil {
			return fmt.Errorf("cannot open spec template: %s", err)
		}

		temp, err := template.New(s.Model.RestName).Funcs(functions).Parse(string(data))
		if err != nil {
			return fmt.Errorf("cannot parse template: %s", err)
		}

		buf := &bytes.Buffer{}
		if err := temp.Execute(buf, struct {
			Set  *spec.SpecificationSet
			Spec *spec.Specification
		}{
			Set:  set,
			Spec: s,
		}); err != nil {
			return fmt.Errorf("cannot execute template: %s", err)
		}

		fmt.Print(buf.String())

		if i <= len(specs)-skipped {
			fmt.Println()
		}
	}

	return nil
}
