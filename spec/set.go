package spec

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"
)

// A SpecificationSet represents a compete set of Specification
type SpecificationSet struct {
	Configuration *Config
	ExternalTypes TypeMapping
	APIInfo       *APIInfo

	specs map[string]*Specification
}

// NewSpecificationSet loads and parses all specification in a folder.
func NewSpecificationSet(
	dirname string,
	nameConvertFunc AttributeNameConverterFunc,
	typeConvertFunc AttributeTypeConverterFunc,
	typeMappingName string,
) (*SpecificationSet, error) {

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

		case "_type.mapping":

			set.ExternalTypes, err = LoadTypeMapping(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

		case "_api.info":
			set.APIInfo, err = LoadAPIInfo(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

		default:

			if path.Ext(info.Name()) != ".spec" && path.Ext(info.Name()) != ".abs" {
				continue
			}

			targetMap := set.specs

			if path.Ext(info.Name()) == ".abs" {
				targetMap = baseSpecs
			}

			baseName := strings.Replace(strings.Replace(info.Name(), ".spec", "", 1), ".abs", "", 1)

			targetMap[baseName], err = LoadSpecification(path.Join(dirname, info.Name()), false)
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
		for _, ext := range spec.Model.Extends {

			base, ok := baseSpecs[ext]
			if !ok {
				return nil, fmt.Errorf("Unable to find base spec '%s' for spec '%s'", ext, spec.Model.RestName)
			}

			if err = spec.ApplyBaseSpecifications(base); err != nil {
				return nil, err
			}
		}

		// Link the APIs to corresponding specifications
		for _, rel := range spec.Relations {

			linked, ok := set.specs[rel.RestName]
			if !ok {
				return nil, fmt.Errorf("Unable to find related spec '%s' for spec '%s'", rel.RestName, spec.Model.RestName)
			}

			rel.remoteSpecification = linked
		}

		if set.ExternalTypes != nil {

			for _, attr := range spec.Attributes {

				if nameConvertFunc != nil {
					attr.ConvertedName = nameConvertFunc(attr.Name)
				}
				if typeConvertFunc != nil {
					attr.ConvertedType, attr.TypeProvider = typeConvertFunc(attr.Type, attr.SubType)
				}

				if attr.Type != AttributeTypeExt {
					continue
				}

				if typeMappingName != "" {
					m, err := set.ExternalTypes.Mapping(typeMappingName, attr.SubType)
					if err != nil {
						return nil, fmt.Errorf("Cannot apply type mapping for attribute '%s' for subtype '%s'", attr.Name, attr.SubType)
					}

					if m != nil {
						attr.ConvertedType = m.Type
						attr.Initializer = m.Initializer
						attr.TypeProvider = m.Import
					} else {
						attr.ConvertedType = string(attr.Type)
					}
				}
			}
		}
	}

	var errs []error
	for _, spec := range set.Specifications() {
		if err := spec.Validate(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, formatValidationErrors(errs)
	}

	return set, nil
}

// Specification returns the Specification with the given name.
func (s *SpecificationSet) Specification(name string) *Specification {
	return s.specs[name]
}

// Specifications returns all Specifications.
func (s *SpecificationSet) Specifications() (specs []*Specification) {

	for _, s := range s.specs {
		specs = append(specs, s)
	}

	sort.Slice(specs, func(i int, j int) bool {
		return strings.Compare(specs[i].Model.RestName, specs[j].Model.RestName) == -1
	})
	return
}

// Len returns the number of specifications in the set.
func (s *SpecificationSet) Len() int {
	return len(s.specs)
}

// Relationships is better
func (s *SpecificationSet) Relationships() map[string]*Relationship {

	relationships := map[string]*Relationship{}

	for _, spec := range s.Specifications() {
		relationships[spec.Model.EntityName] = NewRelationship()
	}

	for _, spec := range s.Specifications() {

		if !spec.Model.IsRoot {
			if spec.Model.AllowsUpdate {
				relationships[spec.Model.EntityName].Set("update", "root")
			}
			if spec.Model.AllowsDelete {
				relationships[spec.Model.EntityName].Set("delete", "root")
			}
			if spec.Model.AllowsGet {
				relationships[spec.Model.EntityName].Set("get", "root")
			}
		}

		for _, rel := range spec.Relations {

			childrenSpec := s.specs[rel.RestName]

			if rel.AllowsGet {
				relationships[childrenSpec.Model.EntityName].Set("getmany", spec.Model.RestName)
			}
			if rel.AllowsCreate {
				relationships[childrenSpec.Model.EntityName].Set("create", spec.Model.RestName)
			}

		}
	}

	return relationships
}

// RelationshipsByRestName returns the relationships indexed by rest name.
func (s *SpecificationSet) RelationshipsByRestName() map[string]*Relationship {

	relationships := map[string]*Relationship{}

	for _, spec := range s.Specifications() {
		relationships[spec.Model.RestName] = NewRelationship()
	}

	for _, spec := range s.Specifications() {

		if !spec.Model.IsRoot {
			if spec.Model.AllowsUpdate {
				relationships[spec.Model.RestName].Set("update", "root")
			}
			if spec.Model.AllowsDelete {
				relationships[spec.Model.RestName].Set("delete", "root")
			}
			if spec.Model.AllowsGet {
				relationships[spec.Model.RestName].Set("get", "root")
			}
		}

		for _, rel := range spec.Relations {

			if rel.AllowsGet {
				relationships[rel.RestName].Set("getmany", spec.Model.RestName)
			}
			if rel.AllowsCreate {
				relationships[rel.RestName].Set("create", spec.Model.RestName)
			}

		}
	}

	return relationships
}

// RelationshipsByResourceName returns the relationships indexed by resource name.
func (s *SpecificationSet) RelationshipsByResourceName() map[string]*Relationship {

	relationships := map[string]*Relationship{}

	for _, spec := range s.Specifications() {
		relationships[spec.Model.ResourceName] = NewRelationship()
	}

	for _, spec := range s.Specifications() {

		if !spec.Model.IsRoot {
			if spec.Model.AllowsUpdate {
				relationships[spec.Model.ResourceName].Set("update", "root")
			}
			if spec.Model.AllowsDelete {
				relationships[spec.Model.ResourceName].Set("delete", "root")
			}
			if spec.Model.AllowsGet {
				relationships[spec.Model.ResourceName].Set("get", "root")
			}
		}

		for _, rel := range spec.Relations {

			childrenSpec := s.specs[rel.RestName]

			if rel.AllowsGet {
				relationships[childrenSpec.Model.ResourceName].Set("getmany", spec.Model.RestName)
			}
			if rel.AllowsCreate {
				relationships[childrenSpec.Model.ResourceName].Set("create", spec.Model.RestName)
			}

		}
	}

	return relationships
}
