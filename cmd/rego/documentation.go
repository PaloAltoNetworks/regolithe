package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/aporeto-inc/regolithe/cmd/rego/static"
	"github.com/aporeto-inc/regolithe/spec"
)

func shortString(s string, max int) string {
	s = strings.Split(s, "\n")[0]
	if len(s) < max {
		return s
	}
	return s[:max] + "..."
}

var functions = template.FuncMap{
	"join":            strings.Join,
	"title":           strings.Title,
	"trimspace":       strings.TrimSpace,
	"toc":             toc,
	"operations":      operations,
	"characteristics": characteristics,
}

func writeDoc(set *spec.SpecificationSet, format string) error {

	switch format {
	case "markdown":
		return writeMarkdownDoc(set)
	default:
		return fmt.Errorf("Unsuported documentation format: %s", format)
	}
}

func writeMarkdownDoc(set *spec.SpecificationSet) error {

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

	relationships := set.RelationshipsByRestName()

	fmt.Println(buf.String())

	var out string
	r := regexp.MustCompile(`\n\n\n+`)

	for _, s := range set.Specifications() {

		if s.Model.Private || s.Model.IsRoot {
			continue
		}

		data, err := static.Asset("templates/spec-md.gotpl")
		if err != nil {
			return fmt.Errorf("cannot open spec template: %s", err)
		}

		temp, err := template.New(s.Model.RestName).Funcs(functions).Parse(string(data[:len(data)-1]))
		if err != nil {
			return fmt.Errorf("cannot parse template: %s", err)
		}

		buf := &bytes.Buffer{}
		if err := temp.Execute(buf, struct {
			Set           *spec.SpecificationSet
			Spec          *spec.Specification
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

	fmt.Print(out[:len(out)-1])

	return nil
}
