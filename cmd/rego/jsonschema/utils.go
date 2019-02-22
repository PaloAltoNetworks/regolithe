package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"go.aporeto.io/regolithe/cmd/rego/static"
	"go.aporeto.io/regolithe/spec"
)

func convertType(t spec.AttributeType) string {
	switch t {
	case spec.AttributeTypeString:
		return "string"
	case spec.AttributeTypeInt:
		return "integer"
	case spec.AttributeTypeFloat:
		return "float"
	case spec.AttributeTypeBool:
		return "boolean"
	case spec.AttributeTypeTime:
		return "time"
	case spec.AttributeTypeObject:
		return "object"
	case spec.AttributeTypeRefList, spec.AttributeTypeList:
		return "$list"
	case spec.AttributeTypeRefMap:
		return "$map"
	case spec.AttributeTypeRef:
		return "$ref"
	case spec.AttributeTypeExt:
		return "$external"
	default:
	}
	return "UNRECOGNIZED_TYPE/" + string(t)
}

func convertRegexp(str string, required bool) string {
	escaped := strings.Replace(str, "\\", "\\\\", -1)
	if required {
		return escaped
	}
	return fmt.Sprintf("(%s)?", escaped)
}

func isNil(target interface{}) bool {
	return target == nil
}

func jsonStringify(target interface{}) string {
	data, err := json.Marshal(target)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(data)
}

func stripFirstLevelBrackets(str string) string {
	return strings.TrimSuffix(strings.TrimPrefix(str, "{"), "}")
}

func makeTemplate(p string) (*template.Template, error) {

	data, err := static.Asset(p)
	if err != nil {
		return nil, err
	}

	return template.New(path.Base(p)).Funcs(functions).Parse(string(data))
}

func writeFile(path string, data []byte) error {

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Unable to write file: %s", f.Name())
	}

	defer f.Close() // nolint: errcheck
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("Unable to write file: %s", f.Name())
	}

	return nil
}
