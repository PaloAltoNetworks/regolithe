package spec

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

// An APIInfo holds general information about the API.
type APIInfo struct {
	Prefix  string `yaml:"prefix,omitempty"`
	Root    string `yaml:"root,omitempty"`
	Version string `yaml:"version,omitempty"`
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

	return apiinfo, nil
}
