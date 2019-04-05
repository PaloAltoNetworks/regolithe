package doc

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"go.aporeto.io/regolithe/cmd/rego/static"
	"go.aporeto.io/regolithe/spec"
)

func shortString(s string, max int) string {
	s = strings.Split(s, "\n")[0]
	if len(s) < max {
		return s
	}
	return s[:max] + "..."
}

var functions = template.FuncMap{
	"join":             strings.Join,
	"title":            strings.Title,
	"trimspace":        strings.TrimSpace,
	"toc":              toc,
	"operations":       operations,
	"characteristics":  characteristics,
	"example":          makeExample,
	"typeOf":           typeOf,
	"makeDefaultValue": makeDefaultValue,
}

// Write writes the documentation for the given spec.SpecificationSet.
func Write(set spec.SpecificationSet, format string) error {

	switch format {
	case "markdown":
		return writeMarkdownDoc(set)
	default:
		return fmt.Errorf("unsuported documentation format: %s", format)
	}
}

func writeMarkdownDoc(set spec.SpecificationSet) error {

	tocData, err := static.Asset("templates/toc-md.gotpl")
	if err != nil {
		return fmt.Errorf("cannot open toc template: %s", err)
	}

	specData, err := static.Asset("templates/spec-md.gotpl")
	if err != nil {
		return fmt.Errorf("cannot open spec template: %s", err)
	}

	relationships := set.RelationshipsByRestName()

	out := fmt.Sprintf(`<!-- markdownlint-disable MD024 MD025 -->

# %s API Documentation

> Version: %d

`, strings.Title(set.Configuration().ProductName), set.APIInfo().Version)

	r := regexp.MustCompile(`\n\n\n+`)

	groups := set.Groups()

	for _, g := range groups {

		specs := set.SpecificationGroup(g)

		if len(specs) == 0 || g == "none" {
			continue
		}

		buf := &bytes.Buffer{}
		s := struct {
			Set       spec.SpecificationSet
			GroupName string
		}{
			set,
			g,
		}

		tocTemplate, err := template.New("toc-" + g).Funcs(functions).Parse(string(tocData))
		if err != nil {
			return fmt.Errorf("cannot parse template: %s", err)
		}

		if err := tocTemplate.Execute(buf, s); err != nil {
			return fmt.Errorf("cannot execute template: %s", err)
		}

		var initializedGroup bool

		for _, s := range specs {

			model := s.Model()

			if model.Private || model.IsRoot {
				continue
			}

			if !initializedGroup {
				out = fmt.Sprintf("%s%s\n", out, buf.String())
				initializedGroup = true
			}

			temp, err := template.New(model.RestName).Funcs(functions).Parse(string(specData[:len(specData)-1]))
			if err != nil {
				return fmt.Errorf("cannot parse template: %s", err)
			}

			buf := &bytes.Buffer{}
			if err := temp.Execute(buf, struct {
				Set           spec.SpecificationSet
				Spec          spec.Specification
				Relationships map[string]*spec.Relationship
			}{
				Set:           set,
				Spec:          s,
				Relationships: relationships,
			}); err != nil {
				return fmt.Errorf("cannot execute template: %s", err)
			}

			out = out + r.ReplaceAllString(buf.String(), "\n\n")
		}
	}

	fmt.Print(out[:len(out)-1])

	return nil
}
