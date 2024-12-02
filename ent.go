package ent

import (
	"embed"
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

//go:embed templates/*.tmpl
var tmpldir embed.FS

// Extension for extra ent functions.
type Extension struct {
	entc.DefaultExtension
}

func (*Extension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(
			gen.NewTemplate("eidng8_ent_templates").ParseFS(
				tmpldir, "templates/*.tmpl",
			),
		),
	}
}

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
