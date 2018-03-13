package main

import (
	"fmt"
	"io/ioutil"

	"github.com/aporeto-inc/regolithe/cmd/rego/static"
	"github.com/xeipuuv/gojsonschema"
)

func validate(specFile string) error {

	schemaData, err := static.Asset("schema/regolithe.json")
	if err != nil {
		return err
	}

	docData, err := ioutil.ReadFile(specFile)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewStringLoader(string(schemaData))
	documentLoader := gojsonschema.NewStringLoader(string(docData))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	for _, err := range result.Errors() {
		fmt.Println(err.String())
	}

	return nil
}
