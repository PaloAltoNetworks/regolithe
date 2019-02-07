package jsonschema

import (
	"go.aporeto.io/regolithe/spec"
)

// Generate generates the json schema
func Generate(set spec.SpecificationSet, outFolder string) error {

	if err := writeModel(set, outFolder); err != nil {
		return err
	}

	if err := writeGlobal(set, outFolder); err != nil {
		return err
	}

	return nil
}
