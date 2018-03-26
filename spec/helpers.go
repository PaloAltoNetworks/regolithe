package spec

import (
	"fmt"
	"sort"
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

func makeSchemaValidationError(message string, res []gojsonschema.ResultError) error {

	var out []string
	for _, r := range res {
		out = append(out, fmt.Sprintf("- %s", r.String()))
	}

	sort.Strings(out)

	return fmt.Errorf("%s:\n%s", message, strings.Join(out, "\n"))
}
