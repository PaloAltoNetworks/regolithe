package spec

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

// A SpecificationSet represents a compete set of Specification
type SpecificationSet struct {
	Configuration *Config
	ExternalTypes *TypeMapping

	specs map[string]*Specification
}

// NewSpecificationSet loads and parses all specification in a folder.
func NewSpecificationSet(dirname string) (*SpecificationSet, error) {

	var loadedMonolitheINI bool

	set := &SpecificationSet{
		specs: map[string]*Specification{},
	}

	filesInfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	baseSpecs := map[string]*Specification{}

	for _, info := range filesInfo {

		switch info.Name() {

		case "monolithe.ini":

			set.Configuration, err = LoadConfig(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

			loadedMonolitheINI = true
			continue

		case "type_mapping.ini":

			set.ExternalTypes, err = LoadTypeMapping(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

			continue

		default:

			if path.Ext(info.Name()) != ".spec" {
				continue
			}

			baseName := strings.Replace(info.Name(), ".spec", "", 1)

			targetMap := set.specs
			if strings.HasPrefix(info.Name(), "@") {
				targetMap = baseSpecs
			}

			targetMap[baseName], err = LoadSpecification(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}
		}
	}

	if !loadedMonolitheINI {
		return nil, fmt.Errorf("Could not find monolithe.ini in folder %s", dirname)
	}

	// Massage the specs
	for _, spec := range set.specs {

		// Apply base specs.
		for _, ext := range spec.Extends {

			base, ok := baseSpecs[ext]
			if !ok {
				return nil, fmt.Errorf("Unable to find base spec '%s' for spec '%s'", ext, spec.RestName)
			}

			if err = spec.ApplyBaseSpecifications(base); err != nil {
				return nil, err
			}
		}

		// Link the APIs to corresponding specifications
		for _, api := range spec.APIs {

			linked, ok := set.specs[api.RestName]
			if !ok {
				return nil, fmt.Errorf("Unable to find linked spec '%s' for spec '%s'", api.RestName, spec.RestName)
			}

			api.linkedSpecification = linked
		}
	}

	return set, nil
}

// Specification returns the Specification with the given name.
func (s *SpecificationSet) Specification(name string) *Specification {
	return s.specs[name]
}

// Len returns the number of specifications in the set.
func (s *SpecificationSet) Len() int {
	return len(s.specs)
}
