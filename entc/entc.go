package entc

import (
	"embed"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed templates/*/*.tmpl
var tmpldir embed.FS

// ClientExtension for extra ent functions.
type ClientExtension struct {
	entc.DefaultExtension
}

func (*ClientExtension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(
			gen.NewTemplate("eidng8_ent_templates").
				ParseFS(tmpldir, "templates/client/*.tmpl"),
		),
	}
}

// SimpleTreeExtension for simple-tree.
type SimpleTreeExtension struct {
	entc.DefaultExtension
}

func (*SimpleTreeExtension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(
			gen.NewTemplate("query_cte").
				ParseFS(tmpldir, "templates/simpletree/*.tmpl"),
		),
	}
}
