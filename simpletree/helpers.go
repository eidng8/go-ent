package simpletree

import (
	"embed"

	"github.com/ogen-go/ogen"
)

//go:embed *.tmpl
var dir embed.FS

// AttachTo adds the `recurse` parameter to the path item.
func AttachTo(item *ogen.Operation) {
	item.AddParameters(RecurseParam())
}

// RemoveFields removes the specified fields from the properties.
func RemoveFields(props []ogen.Property, fields ...string) []ogen.Property {
	for _, field := range fields {
		for i, prop := range props {
			if prop.Name == field {
				props = append(props[:i], props[i+1:]...)
				break
			}
		}
	}
	return props
}

// RemoveEdges removes the `parent` and `children` fields from the OpenAPI
// schema.
func RemoveEdges(op *ogen.Operation) {
	schema := op.RequestBody.Content["application/json"].Schema
	schema.Properties = RemoveFields(schema.Properties, "parent", "children")
}

// RecurseParam returns the `recurse` parameter.
func RecurseParam() *ogen.Parameter {
	return &ogen.Parameter{
		Name:        "recurse",
		In:          "query",
		Description: "Whether to return all descendants (recurse to last leaf)",
		Required:    false,
		Schema:      &ogen.Schema{Type: "boolean"},
	}
}
