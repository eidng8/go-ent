package softdelete

import (
	"context"
	"fmt"
	"path"

	"github.com/go-faster/errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/ogen-go/ogen"
)

// ParamTrashed is the name of the query parameter to include trashed items
const ParamTrashed = "trashed"

// AttachTo adds fields, parameters, and endpoints necessary to the soft delete
// pattern to the given OpenAPI spec, with the given ID parameter name.
func AttachTo(
	spec *ogen.Spec, base string, item *ogen.Schema, idParam *ogen.Parameter,
) error {
	AddDeletedAtField(item)
	ep, exists := spec.Paths[base]
	if exists {
		ep.Get.AddParameters(TrashedParam())
	}
	p := path.Join(base, fmt.Sprintf("{%s}", idParam.Name))
	ep, exists = spec.Paths[p]
	if exists {
		ep.Get.AddParameters(TrashedParam())
		ep.Delete.AddParameters(TrashedParam())
	}
	ep, exists = spec.Paths[path.Join(p, "restore")]
	if !exists {
		return AddRestoreEndpoint(spec, p, idParam)
	}
	return nil
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
		schema.Properties,
		ogen.Property{
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
func AddRestoreEndpoint(
	spec *ogen.Spec, basePath string, idParam *ogen.Parameter,
) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	marshal, err := json.Marshal(idParam)
	if err != nil {
		return errors.Errorf("failed to marshal idParam: %v", err)
	}
	var param ogen.Parameter
	err = json.Unmarshal(marshal, &param)
	if err != nil {
		return errors.Errorf("failed to unmarshal idParam: %v", err)
	}
	endpoint := path.Join(basePath, "restore")
	spec.Paths[endpoint] = &ogen.PathItem{
		Post: &ogen.Operation{
			Summary:     "Restore a trashed record",
			Description: "Restore a record that was previously soft deleted",
			OperationID: "RestoreSimpleTree",
			Parameters:  []*ogen.Parameter{&param},
			Responses: map[string]*ogen.Response{
				"204": {Description: "Record with requested ID was restored"},
				"400": {Ref: "#/components/responses/400"},
				"404": {Ref: "#/components/responses/404"},
				"409": {Ref: "#/components/responses/409"},
				"500": {Ref: "#/components/responses/500"},
			},
		},
	}
	return nil
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
