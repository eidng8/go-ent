package softdelete

import (
	"context"
	"fmt"
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// FieldDeletedAt holds the column name for the "deleted_at" field.
const FieldDeletedAt = "deleted_at"

type softDeleteKey struct{}

// Mixin implements the soft delete pattern for schemas.
type Mixin struct {
	mixin.Schema
}

// Fields of the SoftDeleteMixin.
// Once you declare "deleted_at" in here, you MUST DELETE IT from the entity that will use that Mixin
func (Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time(FieldDeletedAt).Optional().Nillable().
			Annotations(entoas.Skip(true)),
	}
}

// Interceptor returns a new ent.Interceptor that implements the soft delete pattern.
func Interceptor[Q interface{ WhereP(...func(*sql.Selector)) }](
	f func(ent.Query) (Q, error),
) ent.Interceptor {
	return ent.InterceptFunc(
		func(next ent.Querier) ent.Querier {
			return ent.QuerierFunc(
				func(ctx context.Context, query ent.Query) (
					ent.Value, error,
				) {
					if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
						return next.Query(ctx, query)
					}
					q, err := f(query)
					if err != nil {
						return nil, err
					}
					q.WhereP(sql.FieldIsNull(FieldDeletedAt))
					return next.Query(ctx, query)
				},
			)
		},
	)
}

// Mutator returns a new ent.Hook that implements the soft delete pattern.
func Mutator[M interface {
	Mutate(context.Context, ent.Mutation) (ent.Value, error)
}]() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(
			func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
					return next.Mutate(ctx, m)
				}
				mx, ok := m.(interface {
					Op() ent.Op
					SetOp(ent.Op)
					WhereP(...func(*sql.Selector))
				})
				if !ok {
					return nil, fmt.Errorf(
						"unexpected mutation type %T %#v", m, m,
					)
				}
				mx.WhereP(sql.FieldIsNull(FieldDeletedAt))
				if mx.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
					err := m.SetField(FieldDeletedAt, time.Now())
					if err != nil {
						return nil, err
					}
					mx.SetOp(ent.OpUpdate)
					md, ok := m.(interface{ Client() M })
					if !ok {
						return nil, fmt.Errorf(
							"unexpected mutation type %T %#v", m, m,
						)
					}
					return md.Client().Mutate(ctx, m)
				}
				return next.Mutate(ctx, m)
			},
		)
	}
}
