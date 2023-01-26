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
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/mitchellh/copystructure"
	wordwrap "github.com/mitchellh/go-wordwrap"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

const (
	rootModelKey        = "model"
	rootIndexesKey      = "indexes"
	rootDefaultOrderKey = "default_order"
	rootAttributesKey   = "attributes"
	rootRelationsKey    = "relations"
)

type versionedAttributes map[string][]*Attribute
type attributeMapping map[string]map[string]*Attribute
type relationMapping map[string]*Relation

type specification struct {
	RawAttributes   versionedAttributes `yaml:"attributes,omitempty"    json:"attributes,omitempty"`
	RawRelations    []*Relation         `yaml:"relations,omitempty"     json:"relations,omitempty"`
	RawModel        *Model              `yaml:"model,omitempty"         json:"model,omitempty"`
	RawIndexes      [][]string          `yaml:"indexes,omitempty"       json:"indexes,omitempty"`
	RawDefaultOrder []string            `yaml:"default_order,omitempty" json:"default_order,omitempty"`

	attributeMap attributeMapping
	relationsMap relationMapping
	identifier   *Attribute
	path         string
}

// NewSpecification returns a new specification.
func NewSpecification() Specification {

	return &specification{}
}

// LoadSpecification returns a new specification using the given file path.
func LoadSpecification(specPath string, validate bool) (Specification, error) {

	file, err := os.Open(specPath) // #nosec
	if err != nil {
		return nil, err
	}
	// #nosec G307
	defer file.Close() // nolint: errcheck

	spec := &specification{}
	spec.path = specPath

	if err = spec.Read(file, validate); err != nil {
		return nil, err
	}

	return spec, nil
}

func massageYAML(in any) any {

	var out any

	switch m := in.(type) {

	case map[any]any:
		c := map[string]any{}
		for k, v := range m {
			c[k.(string)] = massageYAML(v)
		}
		out = c

	case []any:
		c := make([]any, len(m))
		for i, v := range m {
			c[i] = massageYAML(v)
		}
		out = c

	default:
		out = in
	}

	return out
}

// Read loads a specifaction from the given io.Reader
func (s *specification) Read(reader io.Reader, validate bool) (err error) {

	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)

	if err = decoder.Decode(s); err != nil {
		return fmt.Errorf("unable to decode spec yaml: %s", err)
	}

	for _, attrs := range s.RawAttributes {
		for _, attr := range attrs {
			if attr.ExampleValue != nil {
				attr.ExampleValue = massageYAML(attr.ExampleValue)
			}
			if attr.DefaultValue != nil {
				attr.DefaultValue = massageYAML(attr.DefaultValue)
			}
		}
	}

	if err = s.buildAttributesMapping(); err != nil {
		return fmt.Errorf("unable to build attributes mapping: %s", err)
	}

	if err = s.buildRelationsMapping(); err != nil {
		return fmt.Errorf("unable to build relations mapping: %s", err)
	}

	if s.RawModel != nil {
		s.RawModel.EntityNamePlural = Pluralize(s.RawModel.EntityName)
	}

	if validate {
		if errs := s.Validate(); len(errs) != 0 {
			return formatValidationErrors(errs)
		}
	}

	return nil
}

// Write dumps the specification into a []byte.
func (s *specification) Write(writer io.Writer) error {

	repr := yaml.MapSlice{}

	if s.RawModel != nil {
		repr = append(repr, yaml.MapItem{Key: rootModelKey, Value: toYAMLMapSlice(s.RawModel)})
	}

	if len(s.RawDefaultOrder) != 0 {
		repr = append(repr, yaml.MapItem{Key: rootDefaultOrderKey, Value: s.RawDefaultOrder})
	}

	if len(s.RawIndexes) != 0 {
		repr = append(repr, yaml.MapItem{Key: rootIndexesKey, Value: s.RawIndexes})
	}

	if len(s.RawAttributes) != 0 {

		versionedAttrs := yaml.MapSlice{}

		for _, version := range sortVersionStrings(s.AttributeVersions()) {

			currentAttributes := s.RawAttributes[version]
			sortAttributes(currentAttributes)

			attrs := make([]yaml.MapSlice, len(currentAttributes))

			versionedAttrs = append(versionedAttrs, yaml.MapItem{
				Key:   version,
				Value: attrs,
			})

			for i, attr := range currentAttributes {
				attrs[i] = toYAMLMapSlice(attr)
			}
		}

		repr = append(repr, yaml.MapItem{Key: rootAttributesKey, Value: versionedAttrs})
	}

	if len(s.RawRelations) != 0 {

		relations := make([]yaml.MapSlice, len(s.RawRelations))

		rawRelations := s.RawRelations[:]
		sort.Slice(rawRelations, func(i int, j int) bool {
			return strings.Compare(rawRelations[i].RestName, rawRelations[j].RestName) == -1
		})

		for i, rel := range rawRelations {
			if rel.Get != nil {
				rel.Get.Description = wordwrap.WrapString(rel.Get.Description, 80)
				if rel.Get.ParameterDefinition != nil {
					sortParameters(rel.Get.ParameterDefinition.Entries)
				}
			}

			if rel.Create != nil {
				rel.Create.Description = wordwrap.WrapString(rel.Create.Description, 80)
				if rel.Create.ParameterDefinition != nil {
					sortParameters(rel.Create.ParameterDefinition.Entries)
				}
			}

			if rel.Update != nil {
				rel.Update.Description = wordwrap.WrapString(rel.Update.Description, 80)
				if rel.Update.ParameterDefinition != nil {
					sortParameters(rel.Update.ParameterDefinition.Entries)
				}
			}

			if rel.Delete != nil {
				rel.Delete.Description = wordwrap.WrapString(rel.Delete.Description, 80)
				if rel.Delete.ParameterDefinition != nil {
					sortParameters(rel.Delete.ParameterDefinition.Entries)
				}
			}

			relations[i] = toYAMLMapSlice(rel)
		}

		repr = append(repr, yaml.MapItem{Key: rootRelationsKey, Value: relations})
	}

	data, err := yaml.Marshal(repr)
	if err != nil {
		return err
	}

	var previousLine []byte
	buf := &bytes.Buffer{}
	prfx1 := []byte("- ")
	prfx2 := []byte("  - name")
	prfx3 := []byte("  v")
	prfx4 := []byte("      - name")
	sufx1 := []byte(":")
	yamlModelKey := []byte(rootModelKey + ":")
	yamlDefaultOrderKey := []byte(rootDefaultOrderKey + ":")
	yamlIndexesKey := []byte(rootIndexesKey + ":")
	yamlAttrKey := []byte(rootAttributesKey + ":")
	yamlAttrRelation := []byte(rootRelationsKey + ":")

	lines := bytes.Split(data, []byte("\n"))
	lineN := len(lines)

	var inIndexes bool

	for i, line := range lines {

		if bytes.Equal(line, yamlIndexesKey) || bytes.Equal(line, yamlDefaultOrderKey) {
			inIndexes = true
		} else if bytes.Equal(line, yamlAttrKey) {
			inIndexes = false
		}

		condFirstLine := i == 0
		condFirstIn := bytes.Equal(previousLine, yamlAttrKey) || bytes.Equal(previousLine, yamlAttrRelation)
		condPrefixed := bytes.HasPrefix(line, prfx1) ||
			(bytes.HasPrefix(line, prfx3) && !bytes.HasPrefix(line, []byte("  validation"))) ||
			(bytes.HasPrefix(line, prfx2) && !bytes.HasSuffix(previousLine, sufx1)) ||
			(bytes.HasPrefix(line, prfx4) && !bytes.HasSuffix(previousLine, sufx1))

		if !condFirstLine && !condFirstIn && !inIndexes && condPrefixed {
			_, _ = buf.WriteRune('\n')
		}

		if bytes.Equal(line, yamlModelKey) {
			if !condFirstLine {
				_, _ = buf.WriteRune('\n')
			}
			_, _ = buf.WriteString("# Model\n")
		}
		if bytes.Equal(line, yamlDefaultOrderKey) {
			if !condFirstLine {
				_, _ = buf.WriteRune('\n')
			}
			_, _ = buf.WriteString("# Ordering\n")
		}
		if bytes.Equal(line, yamlIndexesKey) {
			if !condFirstLine {
				_, _ = buf.WriteRune('\n')
			}
			_, _ = buf.WriteString("# Indexes\n")
		}
		if bytes.Equal(line, yamlAttrKey) {
			if !condFirstLine {
				_, _ = buf.WriteRune('\n')
			}
			_, _ = buf.WriteString("# Attributes\n")
		}
		if bytes.Equal(line, yamlAttrRelation) {
			if !condFirstLine {
				_, _ = buf.WriteRune('\n')
			}
			_, _ = buf.WriteString("# Relations\n")
		}

		_, _ = buf.Write(line)
		if i+1 < lineN {
			_, _ = buf.WriteRune('\n')
		}

		previousLine = line
	}

	_, err = writer.Write(buf.Bytes())
	return err
}

// Validate validates the spec against the schema.
func (s *specification) Validate() []error {

	var schemaData []byte
	var err error

	if s.RawModel == nil {
		schemaData, err = fs.ReadFile("schema/rego-abstract.json")
	} else {
		schemaData, err = fs.ReadFile("schema/rego-spec.json")
	}

	if err != nil {
		return []error{err}
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(s)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return []error{fmt.Errorf("unable to validate specification: %s", err)}
	}

	var errs []error

	if !res.Valid() {
		if s.RawModel != nil {
			errs = append(errs, makeSchemaValidationError(fmt.Sprintf("%s.spec", s.RawModel.RestName), res.Errors())...)
		} else {
			errs = append(errs, makeSchemaValidationError(path.Base(s.path), res.Errors())...)
		}
	}

	if s.RawModel != nil {
		if e := s.RawModel.Validate(); e != nil {
			errs = append(errs, e...)
		}
	}

	for _, attrs := range s.RawAttributes {
		for _, attr := range attrs {
			if e := attr.Validate(); e != nil {
				errs = append(errs, e...)
			}
		}
	}

	for _, rel := range s.RawRelations {
		if e := rel.Validate(); e != nil {
			errs = append(errs, e...)
		}
	}

	return errs
}

func (s *specification) Model() *Model {

	return s.RawModel
}

func (s *specification) Indexes() [][]string {

	if len(s.RawIndexes) == 1 && s.RawIndexes[0][0] == ":no-inherit" {
		return nil
	}

	sort.Slice(s.RawIndexes, func(i int, j int) bool {
		if s.RawIndexes[i][0][0] == ':' && s.RawIndexes[j][0][0] != ':' {
			return true
		}
		return strings.Compare(s.RawIndexes[i][0], s.RawIndexes[j][0]) != -1
	})

	return s.RawIndexes
}

func (s *specification) DefaultOrder() []string {

	return s.RawDefaultOrder
}

func (s *specification) Relations() []*Relation {
	return s.RawRelations
}

// Attribute returns the Attributes with the given name.
func (s *specification) Attribute(name string, version string) *Attribute {

	versioned, ok := s.attributeMap[version]
	if !ok {
		return nil
	}

	return versioned[name]
}

// AttributeVersions returns the list of all attribute versions.
func (s *specification) AttributeVersions() []string {

	out := make([]string, len(s.RawAttributes))
	var i int
	for v := range s.RawAttributes {
		out[i] = v
		i++
	}

	return out
}

// LatestAttributesVersion returns the latest version
func (s *specification) LatestAttributesVersion() string {

	var max int
	var latest string

	for v := range s.RawAttributes {

		vi, err := versionToInt(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid version '%s'", v))
		}

		if vi >= max {
			latest = v
		}
	}

	return latest
}

// Attributes returns the list of attribute sorted by names.
func (s *specification) Attributes(version string) []*Attribute {

	attrMap := map[string]*Attribute{}

	for _, v := range s.versionsFrom(version) {
		for _, attr := range s.RawAttributes[v] {
			attrMap[attr.Name] = attr
		}
	}

	attrs := make([]*Attribute, len(attrMap))
	var i int
	for _, attr := range attrMap {
		attrs[i] = attr
		i++
	}

	sortAttributes(attrs)

	return attrs
}

// ExposedAttributes returns the exposed attributes.
func (s *specification) ExposedAttributes(version string) []*Attribute {

	var attrs []*Attribute // nolint

	for _, attr := range s.Attributes(version) {
		if !attr.Exposed {
			continue
		}

		attrs = append(attrs, attr)
	}

	return attrs
}

// Relation returns the Relation with the given rest name.
func (s *specification) Relation(restName string) *Relation {
	return s.relationsMap[restName]
}

// Identifier returns all the identifier attribute.
func (s *specification) Identifier() *Attribute {

	return s.identifier
}

// ApplyBaseSpecifications applies attributes of the given *Specifications to the receiver.
func (s *specification) ApplyBaseSpecifications(specs ...Specification) error {

	// The spec has been initialized manually.
	// This should not happen, but let's just be sure we are
	// done with the attr map, or attributes will be duplicated.
	if s.attributeMap == nil {
		if err := s.buildAttributesMapping(); err != nil {
			return err
		}
	}

	for _, candidate := range specs {

		spec, ok := candidate.(*specification)
		if !ok {
			panic("given specification doesn't have the same type")
		}

		if spec.RawAttributes == nil && len(spec.RawIndexes) == 0 {
			continue
		}

		if !s.Model().Detached {
			if len(s.RawIndexes) != 1 || s.RawIndexes[0][0] != ":no-inherit" {
				for _, indexes := range spec.RawIndexes {
					s.RawIndexes = append(s.RawIndexes, indexes)
				}
			}
		}

		for version := range spec.RawAttributes {

			// The spec may have no attributes at all.
			if s.RawAttributes == nil {
				s.RawAttributes = versionedAttributes{}
			}

			for _, attr := range spec.RawAttributes[version] {
				if _, ok := s.attributeMap[version][attr.Name]; !ok {

					attrCopy, err := copystructure.Copy(attr)
					if err != nil {
						return fmt.Errorf("unable to copy attribute '%s' from extension: %w", attr.Name, err)
					}

					s.RawAttributes[version] = append(s.RawAttributes[version], attrCopy.(*Attribute))
				}
			}
		}
	}

	return s.buildAttributesMapping()
}

// TypeProviders returns the unique list of all attributes type providers.
func (s *specification) TypeProviders() []string {

	yes := &struct{}{}
	cache := map[string]*struct{}{}
	var providers []string

	for _, attrs := range s.RawAttributes {
		for _, attr := range attrs {
			if attr.TypeProvider == "" {
				continue
			}

			if _, ok := cache[attr.TypeProvider]; ok {
				continue
			}

			cache[attr.TypeProvider] = yes
			providers = append(providers, attr.TypeProvider)
		}
	}

	return providers
}

// ValidationProviders returns the unique list of all attributes validation providers.
func (s *specification) ValidationProviders() []string {

	yes := &struct{}{}
	cache := map[string]*struct{}{}
	var providers []string

	for _, attrs := range s.RawAttributes {
		for _, attr := range attrs {

			for _, m := range attr.ValidationProviders {

				if _, ok := cache[m.Import]; ok {
					continue
				}

				cache[m.Import] = yes
				if m.Import != "" {
					providers = append(providers, m.Import)
				}
			}
		}
	}

	sort.Strings(providers)
	return providers
}

func (s *specification) buildAttributesMapping() error {

	s.attributeMap = attributeMapping{}
	s.identifier = nil

	for version, attrs := range s.RawAttributes {

		if _, ok := s.attributeMap[version]; !ok {
			s.attributeMap[version] = map[string]*Attribute{}
		}

		for _, attr := range attrs {

			attr.linkedSpecification = s

			if _, ok := s.attributeMap[version][attr.Name]; ok {
				if s.RawModel != nil {
					return fmt.Errorf("specification %s has more than one attribute named %s", s.RawModel.RestName, attr.Name)
				}

				return fmt.Errorf("one abstract has more than one attribute named %s", attr.Name)
			}

			s.attributeMap[version][attr.Name] = attr

			if attr.Identifier {
				if s.identifier != nil {
					return fmt.Errorf("Specification %s has more than one identifier attributes: At least %s and %s", s.RawModel.RestName, s.identifier.Name, attr.Name)
				}
				s.identifier = attr
			}
		}
	}

	return nil
}

func (s *specification) buildRelationsMapping() error {

	s.relationsMap = relationMapping{}
	for _, rel := range s.RawRelations {

		rel.currentSpecification = s

		if _, ok := s.relationsMap[rel.RestName]; ok {
			return fmt.Errorf("Specification has more than one child relation pointing to %s", rel.RestName)
		}

		s.relationsMap[rel.RestName] = rel
	}

	return nil
}

func (s *specification) versionsFrom(version string) []string {

	if version == "" {
		return []string{"v1"}
	}

	initialVersion, err := versionToInt(version)
	if err != nil {
		panic(fmt.Sprintf("Invalid version '%s'", version))
	}

	var versions []int

	for v := range s.RawAttributes {

		currentVersion, err := versionToInt(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid version '%s'", v))
		}

		if currentVersion <= initialVersion {
			versions = append(versions, currentVersion)
		}

	}

	sort.Ints(versions)

	out := make([]string, len(versions))
	for i := range versions {
		out[i] = fmt.Sprintf("v%d", versions[i])
	}

	return out
}
