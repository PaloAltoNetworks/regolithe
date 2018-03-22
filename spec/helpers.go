package spec

import (
	"fmt"
	"os"

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

func writeValidationErrors(message string, res []gojsonschema.ResultError) {

	fmt.Fprintf(os.Stderr, "\n%s:\n", message)

	for _, r := range res {
		fmt.Fprintf(os.Stderr, " - %s\n", r.String())
	}

	fmt.Fprintln(os.Stderr)
}
