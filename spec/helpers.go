package spec

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Pluralize pluralizes the given word.
func Pluralize(word string) string {

	if len(word) == 0 {
		return word
	}

	lastChar := word[len(word)-1:]
	if lastChar == "s" {
		return word
	}

	if len(word) >= 2 {
		secondLastChar := word[len(word)-2 : len(word)-1]

		if lastChar == "y" &&
			secondLastChar != "a" &&
			secondLastChar != "e" &&
			secondLastChar != "i" &&
			secondLastChar != "o" &&
			secondLastChar != "u" &&
			secondLastChar != "y" {
			return word[:len(word)-1] + "ies"
		}
	}

	return word + "s"
}

func makeSchemaValidationError(message string, res []gojsonschema.ResultError) []error {

	var out []error
	for _, r := range res {
		out = append(out, fmt.Errorf("%s: schema error: %s", message, r.String()))
	}

	return out
}

func formatValidationErrors(errs []error) error {

	var out []string
	for _, err := range errs {
		out = append(out, err.Error())
	}

	if len(out) == 0 {
		return nil
	}

	sort.Strings(out)

	return errors.New(strings.Join(out, "\n"))
}

func sortVersionStrings(versions []string) []string {

	var vs []int

	for _, v := range versions {

		currentVersion, err := versionToInt(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid version '%s'", v))
		}

		vs = append(vs, currentVersion)
	}

	sort.Ints(vs)

	var out []string
	for _, v := range vs {
		out = append(out, fmt.Sprintf("v%d", v))
	}

	return out
}

func sortAttributes(attrs []*Attribute) {

	sort.Slice(attrs, func(i int, j int) bool {
		return strings.Compare(attrs[i].Name, attrs[j].Name) == -1
	})
}

func versionToInt(version string) (int, error) {

	return strconv.Atoi(strings.TrimPrefix(version, "v"))
}
