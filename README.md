# go-ent

[![Go Reference](https://pkg.go.dev/badge/github.com/eidng8/go-ent.svg)](https://pkg.go.dev/github.com/eidng8/go-ent)
[![Go Report Card](https://goreportcard.com/badge/github.com/eidng8/go-ent)](https://goreportcard.com/report/github.com/eidng8/go-ent)
[![License](https://img.shields.io/github/license/eidng8/go-ent)](https://github.com/eidng8/go-ent?tab=MIT-1-ov-file#)

A collection of extensions to be used with [ent ORM](https://entgo.io).

## Pagination

A simple pagination module to use with Ent.

### Usage

```golang
//go:build ignore
// +build ignore

// entc.go
package main

import (
    "log"

    "entgo.io/contrib/entoas"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/eidng8/go-ent/paginate"
    "github.com/ogen-go/ogen"
)

func main() {
    oas, err := entoas.NewExtension(
        entoas.Mutations(
            func(g *gen.Graph, s *ogen.Spec) error {
                ep := s.Paths["/base-uri"]
                paginate.FixParamNames(ep.Get.Parameters)
                paginate.SetResponse(ep.Get)
                return nil
            },
        ),
    )
    if err != nil {
        log.Fatalf("creating entoas extension: %v", err)
    }
    err = entc.Generate("./ent/schema", &gen.Config{}, entc.Extensions(oas))
    if err != nil {
        log.Fatalf("running ent codegen: %v", err)
    }
}
```

```golang
package main

import (
    "context"

    "github.com/eidng8/go-ent/paginate"
    "github.com/gin-gonic/gin"

    "your_project/ent"
    "your_project/ent/your_model"
)

// getPage is a function to get a paginated list of your_model
func getPage(ctx context.Context, query *ent.your_query) (*paginate.PaginatedList[ent.your_model], error) { 
    // get the gin context to be used for the pagination
    gc := ctx.(*gin.Context)
    // creates a context to be used for ent query execution, 
    // e.g. if soft delete from the official site is used
    qc := SkipSoftDelete(context.Background())
    pageParams := paginate.GetPaginationParams(gc)
    // MUST be explicitly sorted, doesn't need this line if the query is already sorted
    query.Order(your_model.ByID())
    // optionally, add more predicates if needed
    query.Where(predicate1, predicate2, ...)
    // call `paginate.GetPage()` function to get the paginated result
    page, err := paginate.GetPage[ent.your_model](gc, qc, query, pageParams)
    if err != nil {
        return nil, err
    }
    return page, nil
}
```


## Simple tree

A simple tree module to use with Ent, with predefined column name and CTE (common table expressions).

### Usage

```golang
//go:build ignore
// +build ignore

// entc.go
package main

import (
    "log"

    "entgo.io/contrib/entoas"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/eidng8/go-ent/simpletree"
    "github.com/ogen-go/ogen"
)

func main() {
    oas, err := entoas.NewExtension(
        entoas.Mutations(
            func(g *gen.Graph, s *ogen.Spec) error {
                // we don't need `parent` and `children` fields in list and detail views
                ep := s.Paths["/base-uri"]
                simpletree.RemoveEdges(ep.Post)
                ep = s.Paths["/base-uri/{id}"]
                simpletree.RemoveEdges(ep.Patch)
                // attach necessary endpoints for tree operations
                ep = s.Paths["/base-uri/{id}/children"]
                simpletree.AttachTo(ep.Get)
                return nil
            },
        ),
    )
    if err != nil {
        log.Fatalf("creating entoas extension: %v", err)
    }
    err = entc.Generate("./ent/schema", &gen.Config{}, entc.Extensions(oas))
    if err != nil {
        log.Fatalf("running ent codegen: %v", err)
    }
}
```

```golang
// ent/schema/aschema.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/eidng8/go-ent/simpletree"
)

type ASchema struct {
    ent.Schema
}

func (ASchema) Fields() []ent.Field {
    return []ent.Field{
        // Primary key must be `id`
        field.Uint32("id").Unique().Immutable(),
    }
}

func (ASchema) Mixin() []ent.Mixin {
    return []ent.Mixin{
        // Comment out this when running `go generate` for the first time
        // Choose different ParentXXXMixin[T]{} according to the type of the primary key
        simpletree.ParentU32Mixin[ASchema]{},
    }
}
```


## Soft delete

A simple soft delete module to use with Ent.

### Usage

```golang
//go:build ignore
// +build ignore

// entc.go
package main

import (
    "log"

    "entgo.io/contrib/entoas"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/eidng8/go-ent/softdelete"
    "github.com/ogen-go/ogen"
)

func main() {
    oas, err := entoas.NewExtension(
        entoas.Mutations(
            func(g *gen.Graph, s *ogen.Spec) error {
                softdelete.AttachTo(s, "/base-uri", "YourListItem")
                return nil
            },
        ),
    )
    if err != nil {
        log.Fatalf("creating entoas extension: %v", err)
    }
    err = entc.Generate("./ent/schema", &gen.Config{}, entc.Extensions(oas))
    if err != nil {
        log.Fatalf("running ent codegen: %v", err)
    }
}
```

```golang
// ent/schema/aschema.go
package schema

import (
    "entgo.io/contrib/entoas"
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/ogen-go/ogen"

    "github.com/eidng8/go-ent/softdelete"
    
    gen "<project>/ent"
    "<project>/ent/intercept"
)

type ASchema struct {
    ent.Schema
}

func (ASchema) Fields() []ent.Field {
    return []ent.Field{
        // fields
    }
}

func (ASchema) Mixin() []ent.Mixin {
    return []ent.Mixin{
        // Comment out this when running `go generate` for the first time
        softdelete.Mixin{},
    }
}

func (ASchema) Interceptors() []ent.Interceptor {
    return []ent.Interceptor{
        // Comment out this when running `go generate` for the first time
        softdelete.Interceptor(intercept.NewQuery),
    }
}

func (ASchema) Hooks() []ent.Hook {
    return []ent.Hook{
        // Comment out this when running `go generate` for the first time
        softdelete.Mutator[*gen.Client](),
    }
}
```
