package entitysql

import (
	"context"
	"database/sql"

	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/tlog"
)

type DeleteSpec struct {
	Entity    *EntitySpec
	Predicate func(*Predicate)
	Affected  *int64
}

func NewDeleteSpec(entity string) *DeleteSpec {
	return &DeleteSpec{
		Entity: &EntitySpec{
			Name: entity,
		},
	}
}

func NewDelete(ctx context.Context, drv dialect.Driver, spec *DeleteSpec) error {
	builder := NewDialect(drv.Dialect())
	qb := deleteBuilder{DeleteSpec: spec, entityBuilder: entityBuilder{builder: builder}}
	return qb.delete(ctx, drv)
}

type deleteBuilder struct {
	entityBuilder
	*DeleteSpec
}

func (b *deleteBuilder) delete(ctx context.Context, drv dialect.Driver) error {
	var res sql.Result
	deleter, err := b.deleter(ctx)
	if err != nil {
		return err
	}
	query, args := deleter.Query()
	tlog.Print(query)
	if err := drv.Exec(ctx, query, args, &res); err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	*b.Affected = affected
	return nil
}

func (b *deleteBuilder) deleter(ctx context.Context) (*Deleter, error) {
	deleter := Deleter{}
	deleter.SetDialect(b.builder.dialect)
	deleter.SetEntity(b.Entity.Name)
	if pred := b.Predicate; pred != nil {
		deleter.where = P()
		pred(deleter.where)
	}
	return &deleter, nil
}
