package scim

import (
	"encoding/json"
	"io/ioutil"
)

func NewResourceTypeFromFile(filepath string) (ResourceType, error) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		return ResourceType{}, err
	}

	return NewResourceTypeFromBytes(raw)
}

func NewResourceTypeFromString(s string) (ResourceType, error) {
	return NewResourceTypeFromBytes([]byte(s))
}

func NewResourceTypeFromBytes(raw []byte) (ResourceType, error) {
	err := resourceTypeSchema.validate(raw)
	if err != nil {
		return ResourceType{}, err
	}

	var resourceType resourceType
	err = json.Unmarshal(raw, &resourceType)
	if err != nil {
		return ResourceType{}, err
	}

	return ResourceType{resourceType}, nil
}

// ResourceType specifies the metadata about a resource type.
type ResourceType struct {
	resourceType resourceType
}

// resourceType specifies the metadata about a resource type. Unlike other core resources, all attributes are
// required unless otherwise specified.
//
// RFC: https://tools.ietf.org/html/rfc7643#section-6
type resourceType struct {
	// ID is the resource type's server unique id. This is often the same value as the "name" attribute.
	// OPTIONAL.
	ID string
	// Name is the resource type name. This name is referenced by the "meta.resourceType" attribute in all resources.
	Name string
	// Description is the resource type's human-readable description.
	// OPTIONAL.
	Description string
	// Endpoint is the resource type's HTTP-addressable endpoint relative to the Base URL of the service provider,
	// e.g., "/Users".
	Endpoint string
	// Schema is the resource type's primary/base schema URI, e.g., "urn:ietf:params:scim:schemas:core:2.0:User". This
	// MUST be equal to the "id" attribute of the associated "Schema" resource.
	Schema string
	// SchemaExtensions is a list of URIs of the resource type's schema extensions.
	// OPTIONAL.
	SchemaExtensions []schemaExtension
}

func (r resourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"schemas":          []string{"urn:ietf:params:scim:schemas:core:2.0:ResourceType"},
		"id":               r.ID,
		"name":             r.Name,
		"description":      r.Description,
		"endpoint":         r.Endpoint,
		"schema":           r.Schema,
		"schemaExtensions": r.SchemaExtensions,
	})
}

// schemaExtension is an URI of one of the resource type's schema extensions.
//
// RFC: https://tools.ietf.org/html/rfc7643#section-6
type schemaExtension struct {
	// Schema is the URI of an extended schema, e.g., "urn:edu:2.0:Staff". This MUST be equal to the "id" attribute
	// of a "Schema" resource.
	Schema string
	// Required is a boolean value that specifies whether or not the schema extension is required for the resource
	// type. If true, a resource of this type MUST include this schema extension and also include any attributes
	// declared as required in this schema extension. If false, a resource of this type MAY omit this schema
	// extension.
	Required bool
}

var resourceTypeSchema schema

func init() {
	if err := json.Unmarshal([]byte(rawResourceTypeSchema), &resourceTypeSchema); err != nil {
		panic(err)
	}
}
