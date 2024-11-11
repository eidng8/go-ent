package softdelete

import (
	"context"
	"path"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ogen-go/ogen"
)

// ParamTrashed is the name of the query parameter to include trashed items
const ParamTrashed = "trashed"

// AttachTo adds fields, parameters, and endpoints necessary to the soft delete
// pattern to the given OpenAPI `spec` `base`, with the default ID parameter
// name of `{id}`. The `itemSchema` is the name of the schema to attach the
// `deleted_at` field to.
func AttachTo(spec *ogen.Spec, base string, itemSchema string) {
	AttachAs(spec, base, itemSchema, "{id}")
}

// AttachAs adds fields, parameters, and endpoints necessary to the soft delete
// pattern to the given OpenAPI spec, with the given ID parameter name.
func AttachAs(spec *ogen.Spec, base string, itemSchema string, idParam string) {
	AddDeletedAtField(spec.Components.Schemas[itemSchema])
	ep, ok := spec.Paths[base]
	if ok {
		ep.Get.AddParameters(TrashedParam())
	}
	p := path.Join(base, idParam)
	ep, ok = spec.Paths[p]
	if ok {
		ep.Get.AddParameters(TrashedParam())
		ep.Delete.AddParameters(TrashedParam())
	}
	for _, i := range []string{"parent", "children"} {
		pp := path.Join(p, i)
		ep, ok = spec.Paths[pp]
		if ok {
			ep.Get.AddParameters(TrashedParam())
		}
	}
	ep, ok = spec.Paths[path.Join(p, "restore")]
	if !ok {
		AddRestoreEndpoint(spec, p)
	}
}

// IncludeTrashed returns a new context that skips the soft-delete interceptor/mutators.
func IncludeTrashed(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// NewSoftDeleteQueryContext returns a new context that includes the soft delete pattern.
// If `ctx` is nil, it will create a new context.Background().
// calls `IncludeTrashed` if `withTrashed` is `true`.
func NewSoftDeleteQueryContext(
	withTrashed *bool, ctx context.Context,
) context.Context {
	if nil == ctx {
		ctx = context.Background()
	}
	if nil != withTrashed && *withTrashed {
		ctx = IncludeTrashed(ctx)
	}
	return ctx
}

// AddDeletedAtField adds the "deleted_at" field to the oas schema
func AddDeletedAtField(schema *ogen.Schema) {
	schema.Properties = append(
		schema.Properties, ogen.Property{
			Name: FieldDeletedAt,
			Schema: &ogen.Schema{
				Type:        "string",
				Format:      "date-time",
				Nullable:    true,
				Description: "Date and time when the record was deleted",
			},
		},
	)
}

// AddRestoreEndpoint adds the restore endpoint to the OpenAPI spec
func AddRestoreEndpoint(spec *ogen.Spec, basePath string) {
	endpoint := path.Join(basePath, "restore")
	spec.Paths[endpoint] = &ogen.PathItem{
		Post: &ogen.Operation{
			Parameters: []*ogen.Parameter{
				&ogen.Parameter{
					Name:     "id",
					In:       openapi3.ParameterInPath,
					Required: true,
				},
			},
			Responses: map[string]*ogen.Response{
				"204": {Description: "Record with requested ID was restored"},
				"400": {Ref: "#/components/responses/400"},
				"404": {Ref: "#/components/responses/404"},
				"409": {Ref: "#/components/responses/409"},
				"500": {Ref: "#/components/responses/500"},
			},
		},
	}
}

// TrashedParam returns the `trashed` query parameter
func TrashedParam() *ogen.Parameter {
	return &ogen.Parameter{
		Name:        ParamTrashed,
		In:          "query",
		Description: "Whether to include trashed items",
		Required:    false,
		Schema:      &ogen.Schema{Type: "boolean"},
	}
}
