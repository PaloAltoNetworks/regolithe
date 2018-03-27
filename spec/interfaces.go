package spec

import "io"

// SpecificationSet represents an entire set of specification.
type SpecificationSet interface {

	// Specification returns the Specification with the given name.
	Specification(name string) Specification

	// Specifications returns all Specifications.
	Specifications() (specs []Specification)

	// Len returns the number of specifications in the set.
	Len() int

	// Relationships is better
	Relationships() map[string]*Relationship

	// RelationshipsByRestName returns the relationships indexed by rest name.
	RelationshipsByRestName() map[string]*Relationship

	// RelationshipsByResourceName returns the relationships indexed by resource name.
	RelationshipsByResourceName() map[string]*Relationship

	Configuration() *Config
	ExternalTypes() TypeMapping
	APIInfo() *APIInfo
}

// A Specification is the interface representing a Regolithe Specification.
type Specification interface {

	// Read reads and load the given reader containing a specification.
	// If validates is true, validations will be done.
	Read(reader io.Reader, validate bool) error

	// Writes write the current state of the Specification in the given
	// writer.
	Write(writer io.Writer) error

	// Validate validates the specification content.
	Validate() []error

	// ApplyBaseSpecifications applyes the given abstract specification to
	// the specification.
	ApplyBaseSpecifications(specs ...Specification) error

	// Model returns the Specification model.
	Model() *Model

	// Attribute returns the attribute with the given name in the given version.
	Attribute(name string, version string) *Attribute

	// Attributes returns all attributes for the given version.
	Attributes(version string) []*Attribute

	// ExposedAttributes returns only the exposed attributes in the given version.
	ExposedAttributes(version string) []*Attribute

	// OrderingAttributes returns the list of attributes used for ordering.
	OrderingAttributes(version string) []*Attribute

	// AttributeInitializers returns all the attribute initializers for the
	// given version.
	AttributeInitializers(version string) map[string]interface{}

	// AttributeVersions returns all the versions of attributes.
	AttributeVersions() []string

	// LatestAttributeVersion returns the latest version of the attributes.
	LatestAttributeVersion() string

	// Relations returns the Specification relations.
	Relations() []*Relation

	// Relation returns the relation to the given restName.
	Relation(restName string) *Relation

	// Identitier returns the Attribute used as an identifier.
	Identifier() *Attribute

	// TypeProviders returns all type providers.
	TypeProviders() []string
}
