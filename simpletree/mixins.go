package simpletree

import (
	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

const columnName = "parent_id"

// ParentU8Mixin is a mixin for uint8 primary key.
type ParentU8Mixin[T ent.Interface] struct {
	ent.Schema
}

func (ParentU8Mixin[T]) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8(columnName).Optional().Nillable(),
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
		field.Uint16(columnName).Optional().Nillable(),
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
		field.Uint32(columnName).Optional().Nillable(),
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
		field.Uint64(columnName).Optional().Nillable(),
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
		field.Int8(columnName).Optional().Nillable(),
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
		field.Int16(columnName).Optional().Nillable(),
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
		field.Int32(columnName).Optional().Nillable(),
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
		field.Int64(columnName).Optional().Nillable(),
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
	return []ent.Field{
		field.String(columnName).Optional().Nillable(),
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
		field.UUID(columnName, uuid.New()).Optional().Nillable(),
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
