package ent

import (
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Timestamps returns definitions of created_at and updated_at fields.
func Timestamps() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Optional().Nillable().Default(time.Now).
			Immutable().Annotations(
			// removes the field from the create & update OpenAPI endpoint
			entoas.ReadOnly(true),
		),
		field.Time("updated_at").Optional().Nillable().UpdateDefault(time.Now).
			Immutable().Annotations(
			// removes the field from the create & update OpenAPI endpoint
			entoas.ReadOnly(true),
		),
	}
}
