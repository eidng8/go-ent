package cte

import (
	"embed"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed *.tmpl
var dir embed.FS

type Extension struct {
	entc.DefaultExtension
}

func (*Extension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(gen.NewTemplate("query_cte").ParseFS(dir, "*.tmpl")),
	}
}
