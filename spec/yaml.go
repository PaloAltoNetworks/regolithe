package spec

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/mitchellh/go-wordwrap"
	yaml "gopkg.in/yaml.v2"
)

func toYAMLMapSlice(s interface{}) yaml.MapSlice {

	var out yaml.MapSlice

	for _, field := range structs.Fields(s) {

		if !field.IsExported() {
			continue
		}

		yamlName, omit := splitTags(field.Tag("yaml"))
		if yamlName == "" {
			continue
		}

		if (field.IsZero() || field.Value() == nil) && omit {
			continue
		}

		var v interface{}
		if yamlName == "description" {
			v = wordwrap.WrapString(field.Value().(string), 80)
		} else {
			v = field.Value()
		}

		item := yaml.MapItem{
			Key:   yamlName,
			Value: v,
		}

		out = append(out, item)
	}

	return out
}

func splitTags(tag string) (string, bool) {

	if tag == "" || tag == "-" {
		return "", false
	}

	parts := strings.Split(tag, ",")

	if len(parts) == 1 {
		return parts[0], false
	}

	return parts[0], parts[1] == "omitempty"
}
