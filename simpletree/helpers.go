package simpletree

import (
	"embed"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/ogen-go/ogen"
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

func RemoveEdges(op *ogen.Operation) {
	schema := op.RequestBody.Content["application/json"].Schema
	schema.Properties = RemoveFields(schema.Properties, "parent", "children")
}
