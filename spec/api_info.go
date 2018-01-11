package spec

import (
	"encoding/json"
	"os"
)

// An APIInfo holds general information about the API.
type APIInfo struct {
	Prefix  string `json:"prefix"`
	Root    string `json:"root"`
	Version string `json:"version"`
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

	if err = json.NewDecoder(file).Decode(apiinfo); err != nil {
		return nil, err
	}

	return apiinfo, nil
}
