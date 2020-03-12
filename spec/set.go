// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// AttributeNameConverterFunc is the type of a attribute name conveter.
type AttributeNameConverterFunc func(name string) string

// AttributeTypeConverterFunc is the type of a attribute type conveter.
type AttributeTypeConverterFunc func(typ AttributeType, subtype string) (converted string, provider string)

// A specificationSet represents a compete set of Specification
type specificationSet struct {
	configuration  *Config
	typeMap        TypeMapping
	validationsMap ValidationMapping
	apiInfo        *APIInfo
	parametersMap  ParameterMapping

	specs map[string]Specification
}

// LoadSpecificationSetFromGithub loads a set of specs from github.
func LoadSpecificationSetFromGithub(
	token string,
	repoURL string,
	refName string,
	internalPath string,
	nameConvertFunc AttributeNameConverterFunc,
	typeConvertFunc AttributeTypeConverterFunc,
	typeMappingName string,
) (SpecificationSet, error) {

	var auth transport.AuthMethod
	if token != "" {
		auth = &http.BasicAuth{
			Username: "Bearer",
			Password: token,
		}
	}

	tmpFolder, err := ioutil.TempDir("", "regolithe-refs-head")
	if err != nil {
		return nil, err
	}
	defer func(f string) { _ = os.RemoveAll(f) }(tmpFolder) // nolint: errcheck

	var (
		ref           plumbing.ReferenceName
		needsCheckout bool
	)

	givenHash := plumbing.NewHash(refName)
	if !givenHash.IsZero() {
		ref = plumbing.NewReferenceFromStrings("refs/heads/master", "").Name()
		needsCheckout = true
	} else {
		ref = plumbing.NewReferenceFromStrings("refs/heads/"+refName, "").Name()
	}

	log.Printf("Retrieving repository: ref=%s repo=%s path=%s", refName, repoURL, internalPath)

	cloneFunc := func(folder string, ref plumbing.ReferenceName) (*git.Repository, error) {
		return git.PlainClone(
			folder,
			false,
			&git.CloneOptions{
				URL:           repoURL,
				Progress:      nil,
				ReferenceName: ref,
				Auth:          auth,
			})
	}

	repo, err := cloneFunc(tmpFolder, ref)

	if err != nil {
		if err == plumbing.ErrReferenceNotFound {

			log.Printf("failed to clone with refs/heads: ref=%s repo=%s path=%s err=%s", refName, repoURL, internalPath, err)

			// Need to recreate a folder, get error repository already created otherwise
			// Happened even if old tmp folder is deleted...
			tmpFolder, err = ioutil.TempDir("", "regolithe-refs-tags")
			if err != nil {
				return nil, err
			}
			defer func(f string) { _ = os.RemoveAll(f) }(tmpFolder) // nolint: errcheck

			ref = plumbing.NewReferenceFromStrings("refs/tags/"+refName, "").Name()
			repo, err = cloneFunc(tmpFolder, ref)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if needsCheckout {
		wt, e := repo.Worktree()
		if e != nil {
			return nil, e
		}

		if err = wt.Checkout(
			&git.CheckoutOptions{
				Hash: givenHash,
			}); err != nil {
			return nil, err
		}
	}

	set, err := LoadSpecificationSet(
		path.Join(tmpFolder, internalPath),
		nameConvertFunc,
		typeConvertFunc,
		typeMappingName,
	)
	if err != nil {
		return nil, err
	}

	return set, nil
}

// LoadSpecificationSet loads and parses all specification in a folder.
func LoadSpecificationSet(
	dirname string,
	nameConvertFunc AttributeNameConverterFunc,
	typeConvertFunc AttributeTypeConverterFunc,
	typeMappingName string,
) (SpecificationSet, error) {

	var loadedRegolitheINI bool

	set := &specificationSet{
		specs: map[string]Specification{},
	}

	filesInfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	baseSpecs := map[string]Specification{}

	for _, info := range filesInfo {

		switch info.Name() {

		case "regolithe.ini":

			set.configuration, err = LoadConfig(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

			loadedRegolitheINI = true

		case "_type.mapping":

			set.typeMap, err = LoadTypeMapping(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

		case "_validation.mapping":

			set.validationsMap, err = LoadValidationMapping(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

		case "_parameter.mapping":

			set.parametersMap, err = LoadGlobalParameters(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

		case "_api.info":
			set.apiInfo, err = LoadAPIInfo(path.Join(dirname, info.Name()))
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
			baseName = strings.TrimPrefix(baseName, "+")

			targetMap[baseName], err = LoadSpecification(path.Join(dirname, info.Name()), false)
			if err != nil {
				return nil, err
			}

			if targetMap[baseName].Model() != nil &&
				targetMap[baseName].Model().RestName != baseName {
				return nil, fmt.Errorf("%s: declared rest_name '%s' must be identical to filename without extension", info.Name(), targetMap[baseName].Model().RestName)
			}
		}
	}

	if !loadedRegolitheINI {
		return nil, fmt.Errorf("unable to find regolithe.ini in folder '%s'", dirname)
	}

	// Massage the specs
	for _, spec := range set.specs {

		s := spec.(*specification)

		var skipInheritOrdering bool
		if len(s.RawDefaultOrder) > 0 && s.RawDefaultOrder[0] == ":no-inherit" {
			s.RawDefaultOrder = s.RawDefaultOrder[1:]
			skipInheritOrdering = true
		}

		// Apply base specs.
		var ordering []string
		for _, ext := range spec.Model().Extends {

			base, ok := baseSpecs[ext]
			if !ok {
				return nil, fmt.Errorf("unable to find base spec '%s' for spec '%s'", ext, spec.Model().RestName)
			}

			if !skipInheritOrdering {
				ordering = append(ordering, base.DefaultOrder()...)
			}

			if err = spec.ApplyBaseSpecifications(base); err != nil {
				return nil, err
			}
		}

		s.RawDefaultOrder = append(ordering, s.RawDefaultOrder...)

		// Link the APIs to corresponding specifications
		for _, rel := range spec.Relations() {

			linked, ok := set.specs[rel.RestName]
			if !ok {
				return nil, fmt.Errorf("unable to find related spec '%s' for spec '%s'", rel.RestName, spec.Model().RestName)
			}

			rel.remoteSpecification = linked
		}

		for _, version := range spec.AttributeVersions() {

			for _, attr := range spec.Attributes(version) {

				if attr.ValidationProviders == nil {
					attr.ValidationProviders = map[string]*ValidationMap{}
				}

				if nameConvertFunc != nil {
					attr.ConvertedName = nameConvertFunc(attr.Name)
				} else {
					attr.ConvertedName = attr.Name
				}

				if typeConvertFunc != nil {
					attr.ConvertedType, attr.TypeProvider = typeConvertFunc(attr.Type, attr.SubType)
				}

				if typeMappingName != "" {

					if set.typeMap != nil && attr.Type == AttributeTypeExt {

						m, err := set.typeMap.Mapping(typeMappingName, attr.SubType)
						if err != nil {
							return nil, fmt.Errorf("unable to apply type mapping '%s' to attribute '%s'", attr.SubType, attr.Name)
						}

						if m != nil {
							attr.ConvertedType = m.Type
							attr.Initializer = m.Initializer
							attr.TypeProvider = m.Import
						} else {
							attr.ConvertedType = string(attr.Type)
						}
					}

					if set.validationsMap != nil {

						for _, validationName := range attr.Validations {

							m, err := set.validationsMap.Mapping(typeMappingName, validationName)
							if err != nil {
								return nil, fmt.Errorf("unable to apply validation mapping '%s' to attribute '%s': %s", validationName, attr.Name, err)
							}
							if m == nil {
								continue
							}

							attr.ValidationProviders[m.Name] = m
						}
					}
				}
			}
		}

		if set.parametersMap != nil {

			if spec.Model() != nil {

				if spec.Model().Get != nil {
					for _, key := range spec.Model().Get.ParameterReferences {
						if spec.Model().Get.ParameterDefinition == nil {
							spec.Model().Get.ParameterDefinition = &ParameterDefinition{}
						}
						if err := spec.Model().Get.ParameterDefinition.extend(set.parametersMap[key], key); err != nil {
							return nil, err
						}
					}
				}

				if spec.Model().Update != nil {
					for _, key := range spec.Model().Update.ParameterReferences {
						if spec.Model().Update.ParameterDefinition == nil {
							spec.Model().Update.ParameterDefinition = &ParameterDefinition{}
						}
						if err := spec.Model().Update.ParameterDefinition.extend(set.parametersMap[key], key); err != nil {
							return nil, err
						}
					}
				}

				if spec.Model().Delete != nil {
					for _, key := range spec.Model().Delete.ParameterReferences {
						if spec.Model().Delete.ParameterDefinition == nil {
							spec.Model().Delete.ParameterDefinition = &ParameterDefinition{}
						}
						if err := spec.Model().Delete.ParameterDefinition.extend(set.parametersMap[key], key); err != nil {
							return nil, err
						}
					}
				}

				for _, r := range spec.Relations() {

					if r.Create != nil {
						for _, key := range r.Create.ParameterReferences {
							if r.Create.ParameterDefinition == nil {
								r.Create.ParameterDefinition = &ParameterDefinition{}
							}
							if err := r.Create.ParameterDefinition.extend(set.parametersMap[key], key); err != nil {
								return nil, err
							}
						}
					}

					if r.Get != nil {
						for _, key := range r.Get.ParameterReferences {
							if r.Get.ParameterDefinition == nil {
								r.Get.ParameterDefinition = &ParameterDefinition{}
							}
							if err := r.Get.ParameterDefinition.extend(set.parametersMap[key], key); err != nil {
								return nil, err
							}
						}
					}

					if r.Update != nil {
						for _, key := range r.Update.ParameterReferences {
							if r.Update.ParameterDefinition == nil {
								r.Update.ParameterDefinition = &ParameterDefinition{}
							}
							if err := r.Update.ParameterDefinition.extend(set.parametersMap[key], key); err != nil {
								return nil, err
							}
						}
					}

					if r.Delete != nil {
						for _, key := range r.Delete.ParameterReferences {
							if r.Delete.ParameterDefinition == nil {
								r.Delete.ParameterDefinition = &ParameterDefinition{}
							}
							if err := r.Delete.ParameterDefinition.extend(set.parametersMap[key], key); err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
	}

	var errs []error
	for _, spec := range set.Specifications() {
		if es := spec.Validate(); es != nil {
			errs = append(errs, es...)
		}
	}

	if len(errs) > 0 {
		return nil, formatValidationErrors(errs)
	}

	return set, nil
}

func (s *specificationSet) Configuration() *Config {

	return s.configuration
}

func (s *specificationSet) TypeMapping() TypeMapping {

	return s.typeMap
}

func (s *specificationSet) ValidationMapping() ValidationMapping {

	return s.validationsMap
}

func (s *specificationSet) APIInfo() *APIInfo {

	return s.apiInfo
}

// Specification returns the Specification with the given name.
func (s *specificationSet) Specification(name string) Specification {

	return s.specs[name]
}

// SpecificationGroup returns the Specifications from the given group.
func (s *specificationSet) SpecificationGroup(groupName string) (specs []Specification) {

	for _, s := range s.specs {
		if s.Model().Group == groupName {
			specs = append(specs, s)
		}
	}

	sort.Slice(specs, func(i int, j int) bool {
		return strings.Compare(specs[i].Model().RestName, specs[j].Model().RestName) == -1
	})

	return
}

// Specifications returns all Specifications.
func (s *specificationSet) Specifications() (specs []Specification) {

	for _, s := range s.specs {
		specs = append(specs, s)
	}

	sort.Slice(specs, func(i int, j int) bool {
		return strings.Compare(specs[i].Model().RestName, specs[j].Model().RestName) == -1
	})
	return
}

// Len returns the number of specifications in the set.
func (s *specificationSet) Len() int {

	return len(s.specs)
}

// Relationships is better
func (s *specificationSet) Relationships() map[string]*Relationship {

	relationships := map[string]*Relationship{}

	for _, spec := range s.Specifications() {
		relationships[spec.Model().EntityName] = NewRelationship()
	}

	for _, spec := range s.Specifications() {

		model := spec.Model()
		if !model.IsRoot {
			if model.Update != nil {
				relationships[model.EntityName].Set("update", "root", model.Update)
			}
			if model.Delete != nil {
				relationships[model.EntityName].Set("delete", "root", model.Delete)
			}
			if model.Get != nil {
				relationships[model.EntityName].Set("get", "root", model.Get)
			}
		}

		for _, rel := range spec.Relations() {

			childrenSpec := s.specs[rel.RestName]

			model := spec.Model()
			relatedModed := childrenSpec.Model()

			if rel.Get != nil {
				relationships[relatedModed.EntityName].Set("getmany", model.RestName, rel.Get)
			}
			if rel.Create != nil {
				relationships[relatedModed.EntityName].Set("create", model.RestName, rel.Create)
			}
		}
	}

	return relationships
}

// RelationshipsByRestName returns the relationships indexed by rest name.
func (s *specificationSet) RelationshipsByRestName() map[string]*Relationship {

	relationships := map[string]*Relationship{}

	for _, spec := range s.Specifications() {
		relationships[spec.Model().RestName] = NewRelationship()
	}

	for _, spec := range s.Specifications() {

		model := spec.Model()

		if !model.IsRoot {
			if model.Update != nil {
				relationships[model.RestName].Set("update", "root", model.Update)
			}
			if model.Delete != nil {
				relationships[model.RestName].Set("delete", "root", model.Delete)
			}
			if model.Get != nil {
				relationships[model.RestName].Set("get", "root", model.Get)
			}
		}

		for _, rel := range spec.Relations() {

			if rel.Get != nil {
				relationships[rel.RestName].Set("getmany", model.RestName, rel.Get)
			}

			if rel.Create != nil {
				relationships[rel.RestName].Set("create", model.RestName, rel.Create)
			}
		}
	}

	return relationships
}

// RelationshipsByResourceName returns the relationships indexed by resource name.
func (s *specificationSet) RelationshipsByResourceName() map[string]*Relationship {

	relationships := map[string]*Relationship{}

	for _, spec := range s.Specifications() {
		relationships[spec.Model().ResourceName] = NewRelationship()
	}

	for _, spec := range s.Specifications() {

		model := spec.Model()

		if !model.IsRoot {
			if model.Update != nil {
				relationships[model.ResourceName].Set("update", "root", model.Update)
			}
			if model.Delete != nil {
				relationships[model.ResourceName].Set("delete", "root", model.Delete)
			}
			if model.Get != nil {
				relationships[model.ResourceName].Set("get", "root", model.Get)
			}
		}

		for _, rel := range spec.Relations() {

			childrenSpec := s.specs[rel.RestName]

			if rel.Get != nil {
				relationships[childrenSpec.Model().ResourceName].Set("getmany", model.RestName, rel.Get)
			}
			if rel.Create != nil {
				relationships[childrenSpec.Model().ResourceName].Set("create", model.RestName, rel.Create)
			}

		}
	}

	return relationships
}

// Groups returns the list of all groups
func (s *specificationSet) Groups() []string {

	done := map[string]struct{}{}

	for _, sp := range s.Specifications() {
		done[sp.Model().Group] = struct{}{}
	}

	out := make([]string, len(done))
	var i int
	for k := range done {
		out[i] = k
		i++
	}

	sort.Strings(out)

	return out
}
