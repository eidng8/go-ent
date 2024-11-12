package simpletree

import (
	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen"
)

const columnName = "parent_id"

// ParentU8Mixin is a mixin for uint8 primary key.
type ParentU8Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentU8Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "uint8",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("255"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentU8Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentU16Mixin is a mixin for uint16 primary key.
type ParentU16Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentU16Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Uint16(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "uint16",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("65535"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentU16Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentU32Mixin is a mixin for uint32 primary key.
type ParentU32Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentU32Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "uint32",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("4294967295"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentU32Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentU64Mixin is a mixin for uint64 primary key.
type ParentU64Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentU64Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "uint64",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("18446744073709551615"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentU64Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentI8Mixin is a mixin for int8 primary key.
type ParentI8Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentI8Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Int8(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "int8",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("127"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentI8Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentI16Mixin is a mixin for int16 primary key.
type ParentI16Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentI16Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Int16(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "int16",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("32767"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentI16Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentI32Mixin is a mixin for int32 primary key.
type ParentI32Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentI32Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Int32(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "int32",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("2147483647"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentI32Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentI64Mixin is a mixin for int64 primary key.
type ParentI64Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentI64Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Int64(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "integer",
					Format:      "int64",
					Minimum:     ogen.Num("1"),
					Maximum:     ogen.Num("9223372036854775807"),
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentI64Mixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentStringMixin is a mixin for string primary key.
type ParentStringMixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentStringMixin[T]) Fields() []ent.Field {
	u1 := uint64(1)
	return []ent.Field{
		field.String(columnName).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "string",
					MinLength:   &u1,
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentStringMixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

// ParentUuidMixin is a mixin for UUID primary key.
type ParentUuidMixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentUuidMixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.UUID(columnName, uuid.New()).Optional().Nillable().Annotations(
			// adds constraints to the generated OpenAPI specification
			entoas.Schema(
				&ogen.Schema{
					Type:        "string",
					Format:      "uuid",
					Description: "Parent record ID",
				},
			),
		),
	}
}

func (ParentUuidMixin[T]) Edges() []ent.Edge {
	return getEdges[T]()
}

func getEdges[T ent.Interface]() []ent.Edge {
	return []ent.Edge{
		edge.To("children", T.Type).
			Annotations(
				entsql.OnDelete(entsql.Restrict),
				entoas.ReadOnly(true),
				entoas.Skip(true),
			).
			From("parent").Field("parent_id").Unique(),
	}
}
