package entitysql

import (
	"context"
	"database/sql"

	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/tlog"
)

type UpdateSpec struct {
	Entity    *EntitySpec
	Scan      Scanner
	Sets      []*FieldSpec
	Predicate func(*Predicate)
}

func NewUpdateSpec(entity string, rows []FieldName) *UpdateSpec {
	return &UpdateSpec{
		Entity: &EntitySpec{
			Name: entity,
			Rows: rows,
		},
	}
}

func NewUpdate(ctx context.Context, drv dialect.Driver, spec *UpdateSpec) error {
	builder := NewDialect(drv.Dialect())
	qb := updateBuilder{UpdateSpec: spec, entityBuilder: entityBuilder{builder: builder}}
	return qb.update(ctx, drv)
}

type updateBuilder struct {
	entityBuilder
	*UpdateSpec
}

func (b *updateBuilder) update(ctx context.Context, drv dialect.Driver) error {
	var res sql.Result
	updater, err := b.updater(ctx)
	if err != nil {
		return err
	}
	tx, err := drv.Tx(ctx)
	if err := func() error {
		query, args := updater.Query()
		tlog.Print(query)
		tlog.Print(args)
		return tx.Exec(ctx, query, args, res)
	}(); err != nil {
		return rollback(tx, err)
	}
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (b *updateBuilder) updater(ctx context.Context) (*Updater, error) {
	updater := Updater{}
	updater.SetDialect(b.builder.dialect)
	updater.SetEntity(b.Entity.Name)
	for _, field := range b.Sets {
		updater.Set(field.Column, field.Value)
	}
	if pred := b.Predicate; pred != nil {
		updater.where = P()
		pred(updater.where)
	}
	return &updater, nil
}
