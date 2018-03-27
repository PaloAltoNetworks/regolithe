package spec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/aporeto-inc/regolithe/schema"
	"github.com/xeipuuv/gojsonschema"

	wordwrap "github.com/mitchellh/go-wordwrap"
	yaml "gopkg.in/yaml.v2"
)

const (
	rootModelKey      = "model"
	rootAttributesKey = "attributes"
	rootRelationsKey  = "relations"
)

// A Specification is the interface representing a Regolithe Specification.
type Specification interface {
	Read(reader io.Reader, validate bool) error
	Write(writer io.Writer) error
	Validate() error
	Attribute(name string, version string) *Attribute
	AttributeVersions() []string
	LatestVersion() string
	Attributes(version string) []*Attribute
	ExposedAttributes(version string) []*Attribute
	Relation(restName string) *Relation
	Identifier() *Attribute
	OrderingAttributes(version string) []*Attribute
	TypeProviders() []string
	AttributeInitializers(version string) map[string]interface{}
	Model() *Model
	Relations() []*Relation
	ApplyBaseSpecifications(specs ...Specification) error
}

type specification struct {
	RawAttributes map[string][]*Attribute `yaml:"attributes,omitempty"    json:"attributes,omitempty"`
	RawRelations  []*Relation             `yaml:"relations,omitempty"     json:"relations,omitempty"`
	RawModel      *Model                  `yaml:"model,omitempty"         json:"model,omitempty"`

	attributeMap       map[string]map[string]*Attribute
	relationsMap       map[string]*Relation
	orderingAttributes map[string][]*Attribute
	identifier         *Attribute
	path               string
}

// NewSpecification returns a new specification.
func NewSpecification() Specification {
	return &specification{}
}

// LoadSpecification returns a new specification using the given file path.
func LoadSpecification(specPath string, validate bool) (Specification, error) {

	file, err := os.Open(specPath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint: errcheck

	spec := &specification{}
	spec.path = specPath

	if err = spec.Read(file, validate); err != nil {
		return nil, err
	}

	return spec, nil
}

// Read loads a specifaction from the given io.Reader
func (s *specification) Read(reader io.Reader, validate bool) (err error) {

	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)

	if err = decoder.Decode(s); err != nil {
		return err
	}

	if err = s.buildAttributesInfo(); err != nil {
		return err
	}

	if err = s.buildRelationssInfo(); err != nil {
		return err
	}

	if s.RawModel != nil {
		s.RawModel.EntityNamePlural = Pluralize(s.RawModel.EntityName)
	}

	if validate {
		if err = s.Validate(); err != nil {
			return err
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

	if len(s.RawAttributes) != 0 {

		versionedAttrs := yaml.MapSlice{}

		for _, version := range s.AttributeVersions() {

			currentAttributes := s.RawAttributes[version]
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

		for i, rel := range s.RawRelations {
			for k, v := range rel.Descriptions {
				rel.Descriptions[k] = wordwrap.WrapString(v, 80)
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
	yamlModelKey := []byte(rootModelKey + ":")
	yamlAttrKey := []byte(rootAttributesKey + ":")
	yamlAttrRelation := []byte(rootRelationsKey + ":")

	lines := bytes.Split(data, []byte("\n"))
	lineN := len(lines)

	for i, line := range lines {

		condFirstLine := i == 0
		condFirstIn := bytes.Equal(previousLine, yamlAttrKey) || bytes.Equal(previousLine, yamlAttrRelation)
		condPrefixed := bytes.HasPrefix(line, prfx1)

		if !condFirstLine && !condFirstIn && condPrefixed {
			buf.WriteRune('\n')
		}

		if bytes.Equal(line, yamlModelKey) {
			if !condFirstLine {
				buf.WriteRune('\n')
			}
			buf.WriteString("# Model\n")
		}
		if bytes.Equal(line, yamlAttrKey) {
			if !condFirstLine {
				buf.WriteRune('\n')
			}
			buf.WriteString("# Attributes\n")
		}
		if bytes.Equal(line, yamlAttrRelation) {
			if !condFirstLine {
				buf.WriteRune('\n')
			}
			buf.WriteString("# Relations\n")
		}

		buf.Write(line)
		if i+1 < lineN {
			buf.WriteRune('\n')
		}

		previousLine = line
	}

	_, err = writer.Write(buf.Bytes())
	return err
}

// Validate validates the spec against the schema.
func (s *specification) Validate() error {

	var schemaData []byte
	var err error

	if s.RawModel == nil {
		schemaData, err = schema.Asset("rego-abstract.json")
	} else {
		schemaData, err = schema.Asset("rego-spec.json")
	}

	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(s)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return err
	}

	var errs []error

	if !res.Valid() {
		if s.RawModel != nil {
			errs = append(errs, makeSchemaValidationError(fmt.Sprintf("%s.spec", s.RawModel.RestName), res.Errors())...)
		} else {
			errs = append(errs, makeSchemaValidationError(path.Base(s.path), res.Errors())...)
		}

	}

	if s.RawModel != nil && s.RawModel.Description != "" && s.RawModel.Description[len(s.RawModel.Description)-1] != '.' {
		errs = append(errs, fmt.Errorf("%s.spec: model description must end with a period", s.RawModel.RestName))
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

	return formatValidationErrors(errs)
}

func (s *specification) Model() *Model {
	return s.RawModel
}

func (s *specification) Relations() []*Relation {
	return s.RawRelations
}

// Attribute returns the Attributes with the given name.
func (s *specification) Attribute(name string, version string) *Attribute {

	return s.attributeMap[version][name]
}

// AttributeVersions returns the list of all attribute versions.
func (s *specification) AttributeVersions() []string {

	var out []string
	for v := range s.RawAttributes {
		out = append(out, v)
	}

	return out
}

// LatestVersion returns the latest version
func (s *specification) LatestVersion() string {

	var max int
	var latest string

	for v := range s.RawAttributes {

		vs := strings.TrimPrefix(v, "v")
		vi, err := strconv.Atoi(vs)
		if err != nil {
			panic(fmt.Sprintf("Invalid version '%s'", v))
		}

		if vi > max {
			latest = v
		}
	}

	return latest
}

// Attributes returns the list of attribute sorted by names.
func (s *specification) Attributes(version string) []*Attribute {

	attrs := append([]*Attribute{}, s.RawAttributes[version]...)

	sort.Slice(attrs, func(i int, j int) bool {
		return strings.Compare(attrs[i].Name, attrs[j].Name) == -1
	})

	return attrs
}

// ExposedAttributes returns the exposed attributes.
func (s *specification) ExposedAttributes(version string) []*Attribute {

	var attrs []*Attribute

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

// OrderingAttributes returns all the ordering attribute.
func (s *specification) OrderingAttributes(version string) []*Attribute {
	return append([]*Attribute{}, s.orderingAttributes[version]...)
}

// ApplyBaseSpecifications applies attributes of the given *Specifications to the receiver.
func (s *specification) ApplyBaseSpecifications(specs ...Specification) error {

	for _, candidate := range specs {

		spec, ok := candidate.(*specification)
		if !ok {
			panic("given specification doesn't have the same type")
		}

		if spec.RawAttributes == nil {
			continue
		}

		for version := range spec.RawAttributes {
			for _, attr := range spec.RawAttributes[version] {
				if _, ok := s.attributeMap[version][attr.Name]; !ok {
					s.RawAttributes[version] = append(s.RawAttributes[version], attr)
				}
			}
		}
	}

	if err := s.buildAttributesInfo(); err != nil {
		return err
	}

	return nil
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

// AttributeInitializers returns all initializers of the Attributes.
func (s *specification) AttributeInitializers(version string) map[string]interface{} {

	out := map[string]interface{}{}

	for _, attr := range s.RawAttributes[version] {

		if attr.Initializer != "" {
			out[attr.ConvertedName] = attr.Initializer
			continue
		}

		if attr.DefaultValue != nil {
			if attr.Type == AttributeTypeString || attr.Type == AttributeTypeEnum {
				out[attr.ConvertedName] = `"` + attr.DefaultValue.(string) + `"`
				continue
			}
			out[attr.ConvertedName] = attr.DefaultValue
		}
	}

	return out
}

// buildAttributesInfo builds the attributes map.
func (s *specification) buildAttributesInfo() error {

	s.attributeMap = map[string]map[string]*Attribute{}
	s.orderingAttributes = map[string][]*Attribute{}
	s.identifier = nil

	for version, attrs := range s.RawAttributes {

		if _, ok := s.attributeMap[version]; !ok {
			s.attributeMap[version] = map[string]*Attribute{}
		}

		if _, ok := s.orderingAttributes[version]; !ok {
			s.orderingAttributes[version] = []*Attribute{}
		}

		for _, attr := range attrs {

			attr.linkedSpecification = s

			if _, ok := s.attributeMap[version][attr.Name]; ok {
				if s.RawModel != nil {
					return fmt.Errorf("Specification %s has more than one attribute named %s", s.RawModel.RestName, attr.Name)
				}

				return fmt.Errorf("One abstract has more than one attribute named %s", attr.Name)
			}

			s.attributeMap[version][attr.Name] = attr

			if attr.DefaultOrder {
				s.orderingAttributes[version] = append(s.orderingAttributes[version], attr)
			}

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

// buildRelationssInfo builds the relations map.
func (s *specification) buildRelationssInfo() error {

	s.relationsMap = map[string]*Relation{}
	for _, rel := range s.RawRelations {

		rel.currentSpecification = s

		if _, ok := s.relationsMap[rel.RestName]; ok {
			return fmt.Errorf("Specification has more than one child relation pointing to %s", rel.RestName)
		}

		s.relationsMap[rel.RestName] = rel
	}

	return nil
}
