# go-ent

[![Go Reference](https://pkg.go.dev/badge/github.com/eidng8/go-ent.svg)](https://pkg.go.dev/github.com/eidng8/go-ent)
[![Go Report Card](https://goreportcard.com/badge/github.com/eidng8/go-ent)](https://goreportcard.com/report/github.com/eidng8/go-ent)
[![License](https://img.shields.io/github/license/eidng8/go-ent)](https://github.com/eidng8/go-ent?tab=MIT-1-ov-file#)

A collection of extensions to be used with [ent ORM](https://entgo.io).

## Pagination

A simple pagination module to use with Ent.

### Usage

```golang
package main

import (
	"context"

	"github.com/eidng8/go-ent/paginate"
	"github.com/gin-gonic/gin"

	"your_project/ent"
	"your_project/ent/your_model"
)

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
