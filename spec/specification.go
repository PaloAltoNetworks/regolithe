package spec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"github.com/xeipuuv/gojsonschema"
	"go.aporeto.io/regolithe/schema"

	wordwrap "github.com/mitchellh/go-wordwrap"
	yaml "gopkg.in/yaml.v2"
)

const (
	rootModelKey      = "model"
	rootAttributesKey = "attributes"
	rootRelationsKey  = "relations"
)

type versionedAttributes map[string][]*Attribute
type attributeMapping map[string]map[string]*Attribute
type relationMapping map[string]*Relation

type specification struct {
	RawAttributes versionedAttributes `yaml:"attributes,omitempty"    json:"attributes,omitempty"`
	RawRelations  []*Relation         `yaml:"relations,omitempty"     json:"relations,omitempty"`
	RawModel      *Model              `yaml:"model,omitempty"         json:"model,omitempty"`

	attributeMap       attributeMapping
	relationsMap       relationMapping
	orderingAttributes versionedAttributes
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

	if err = s.buildAttributesMapping(); err != nil {
		return err
	}

	if err = s.buildRelationsMapping(); err != nil {
		return err
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

		for i, rel := range s.RawRelations {
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
	yamlAttrKey := []byte(rootAttributesKey + ":")
	yamlAttrRelation := []byte(rootRelationsKey + ":")

	lines := bytes.Split(data, []byte("\n"))
	lineN := len(lines)

	for i, line := range lines {

		condFirstLine := i == 0
		condFirstIn := bytes.Equal(previousLine, yamlAttrKey) || bytes.Equal(previousLine, yamlAttrRelation)
		condPrefixed := bytes.HasPrefix(line, prfx1) ||
			bytes.HasPrefix(line, prfx3) ||
			(bytes.HasPrefix(line, prfx2) && !bytes.HasSuffix(previousLine, sufx1)) ||
			(bytes.HasPrefix(line, prfx4) && !bytes.HasSuffix(previousLine, sufx1))

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
func (s *specification) Validate() []error {

	var schemaData []byte
	var err error

	if s.RawModel == nil {
		schemaData, err = schema.Asset("rego-abstract.json")
	} else {
		schemaData, err = schema.Asset("rego-spec.json")
	}

	if err != nil {
		return []error{err}
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	specLoader := gojsonschema.NewGoLoader(s)

	res, err := gojsonschema.Validate(schemaLoader, specLoader)
	if err != nil {
		return []error{err}
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

// OrderingAttributes returns all the ordering attribute.
func (s *specification) OrderingAttributes(version string) []*Attribute {

	var out []*Attribute // nolint
	versions := s.versionsFrom(version)

	for _, v := range versions {
		out = append(out, s.orderingAttributes[v]...)
	}

	return out
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

		if spec.RawAttributes == nil {
			continue
		}

		for version := range spec.RawAttributes {

			// The spec may have no attributes at all.
			if s.RawAttributes == nil {
				s.RawAttributes = versionedAttributes{}
			}

			for _, attr := range spec.RawAttributes[version] {
				if _, ok := s.attributeMap[version][attr.Name]; !ok {
					s.RawAttributes[version] = append(s.RawAttributes[version], attr)
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

// AttributeInitializers returns all initializers of the Attributes.
func (s *specification) AttributeInitializers(version string) map[string]interface{} {

	out := map[string]interface{}{}

	for _, attr := range s.Attributes(version) {

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

func (s *specification) buildAttributesMapping() error {

	s.attributeMap = attributeMapping{}
	s.orderingAttributes = versionedAttributes{}
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
