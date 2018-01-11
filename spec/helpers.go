package spec

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
		secondLastChar := word[len(word)-2:]

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
