package spec

import (
	"fmt"
	"os"

	"github.com/aporeto-inc/regolithe/schema"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

// An APIInfo holds general information about the API.
type APIInfo struct {
	Prefix  string `yaml:"prefix,omitempty"     json:"prefix,omitempty"`
	Root    string `yaml:"root,omitempty"       json:"root,omitempty"`
	Version int    `yaml:"version,omitempty"    json:"version,omitempty"`
}

// NewAPIInfo returns a new APIInfo.
func NewAPIInfo() *APIInfo {
	return &APIInfo{}
}

// LoadAPIInfo loads an APIInfo from the given file.
func LoadAPIInfo(path string) (*APIInfo, error) {

	apiinfo := NewAPIInfo()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	if err = yaml.NewDecoder(file).Decode(apiinfo); err != nil {
		return nil, err
	}

	if res, err := apiinfo.Validate(); err != nil {
		writeValidationErrors("validation error in _api.info", res)
		return nil, err
	}

	return apiinfo, nil
}

// Validate validates the api info against the schema.
func (a *APIInfo) Validate() ([]gojsonschema.ResultError, error) {

	schemaData, err := schema.Asset("rego-info.json")
	if err != nil {
		return nil, err
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(a)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return nil, err
	}

	if !res.Valid() {
		return res.Errors(), fmt.Errorf("Invalid _api.info")
	}

	return nil, nil
}
