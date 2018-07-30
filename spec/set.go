package spec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	git "gopkg.in/src-d/go-git.v4"
)

// A specificationSet represents a compete set of Specification
type specificationSet struct {
	configuration *Config
	externalTypes TypeMapping
	apiInfo       *APIInfo

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
	defer func(f string) { os.RemoveAll(f) }(tmpFolder) // nolint: errcheck

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

	logrus.WithFields(logrus.Fields{
		"ref":  refName,
		"repo": repoURL,
		"path": internalPath,
	}).Info("Retrieving repository")

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
			logrus.WithFields(logrus.Fields{
				"err":  err,
				"ref":  refName,
				"repo": repoURL,
				"path": internalPath,
			}).Warn("Trying to clone with refs/tags - failed to clone with refs/heads")

			// Need to recreate a folder, get error repository already created otherwise
			// Happened even if old tmp folder is deleted...
			tmpFolder, err = ioutil.TempDir("", "regolithe-refs-tags")
			if err != nil {
				return nil, err
			}
			defer func(f string) { os.RemoveAll(f) }(tmpFolder) // nolint: errcheck

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

	var loadedMonolitheINI bool

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

		case "monolithe.ini":

			set.configuration, err = LoadConfig(path.Join(dirname, info.Name()))
			if err != nil {
				return nil, err
			}

			loadedMonolitheINI = true

		case "_type.mapping":

			set.externalTypes, err = LoadTypeMapping(path.Join(dirname, info.Name()))
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

			targetMap[baseName], err = LoadSpecification(path.Join(dirname, info.Name()), false)
			if err != nil {
				return nil, err
			}

			if targetMap[baseName].Model() != nil && targetMap[baseName].Model().RestName != baseName {
				return nil, fmt.Errorf("%s: declared rest_name '%s' must be identical to filename without extension", info.Name(), targetMap[baseName].Model().RestName)
			}
		}
	}

	if !loadedMonolitheINI {
		return nil, fmt.Errorf("Could not find monolithe.ini in folder %s", dirname)
	}

	// Massage the specs
	for _, spec := range set.specs {

		// Apply base specs.
		for _, ext := range spec.Model().Extends {

			base, ok := baseSpecs[ext]
			if !ok {
				return nil, fmt.Errorf("Unable to find base spec '%s' for spec '%s'", ext, spec.Model().RestName)
			}

			if err = spec.ApplyBaseSpecifications(base); err != nil {
				return nil, err
			}
		}

		// Link the APIs to corresponding specifications
		for _, rel := range spec.Relations() {

			linked, ok := set.specs[rel.RestName]
			if !ok {
				return nil, fmt.Errorf("Unable to find related spec '%s' for spec '%s'", rel.RestName, spec.Model().RestName)
			}

			rel.remoteSpecification = linked
		}

		if set.externalTypes != nil {

			for _, version := range spec.AttributeVersions() {

				for _, attr := range spec.Attributes(version) {

					if nameConvertFunc != nil {
						attr.ConvertedName = nameConvertFunc(attr.Name)
					} else {
						attr.ConvertedName = attr.Name
					}

					if typeConvertFunc != nil {
						attr.ConvertedType, attr.TypeProvider = typeConvertFunc(attr.Type, attr.SubType)
					}

					if attr.Type != AttributeTypeExt {
						continue
					}

					if typeMappingName != "" {
						m, err := set.externalTypes.Mapping(typeMappingName, attr.SubType)
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

func (s *specificationSet) ExternalTypes() TypeMapping {

	return s.externalTypes
}

func (s *specificationSet) APIInfo() *APIInfo {

	return s.apiInfo
}

// Specification returns the Specification with the given name.
func (s *specificationSet) Specification(name string) Specification {

	return s.specs[name]
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
				relationships[model.EntityName].Set("update", "root")
			}
			if model.Delete != nil {
				relationships[model.EntityName].Set("delete", "root")
			}
			if model.Get != nil {
				relationships[model.EntityName].Set("get", "root")
			}
		}

		for _, rel := range spec.Relations() {

			childrenSpec := s.specs[rel.RestName]

			model := spec.Model()
			relatedModed := childrenSpec.Model()

			if rel.Get != nil {
				relationships[relatedModed.EntityName].Set("getmany", model.RestName)
			}
			if rel.Create != nil {
				relationships[relatedModed.EntityName].Set("create", model.RestName)
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
				relationships[model.RestName].Set("update", "root")
			}
			if model.Delete != nil {
				relationships[model.RestName].Set("delete", "root")
			}
			if model.Get != nil {
				relationships[model.RestName].Set("get", "root")
			}
		}

		for _, rel := range spec.Relations() {

			if rel.Get != nil {
				relationships[rel.RestName].Set("getmany", model.RestName)
			}

			if rel.Create != nil {
				relationships[rel.RestName].Set("create", model.RestName)
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
				relationships[model.ResourceName].Set("update", "root")
			}
			if model.Delete != nil {
				relationships[model.ResourceName].Set("delete", "root")
			}
			if model.Get != nil {
				relationships[model.ResourceName].Set("get", "root")
			}
		}

		for _, rel := range spec.Relations() {

			childrenSpec := s.specs[rel.RestName]

			if rel.Get != nil {
				relationships[childrenSpec.Model().ResourceName].Set("getmany", model.RestName)
			}
			if rel.Create != nil {
				relationships[childrenSpec.Model().ResourceName].Set("create", model.RestName)
			}

		}
	}

	return relationships
}
