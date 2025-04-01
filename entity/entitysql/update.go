package entitysql

import (
	"context"
	"fmt"

	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/dialect"
	"github.com/zodileap/taurus_go/tlog"
)

// UpdateSpec 用于生成更新语句。
type UpdateSpec struct {
	Entity *EntitySpec
	Scan   Scanner
	// Sets 更新中Set部分，根据不同的列设置不同的值。
	Sets []map[string][]CaseSpec
	// Predicate 更新中Where部分，根据不同的行设置不同的条件。
	Predicate []PredicateFunc
}

// NewUpdateSpec 创建一个UpdateSpec。
func NewUpdateSpec(entity string, rows []FieldName) *UpdateSpec {
	return &UpdateSpec{
		Entity: &EntitySpec{
			Name:    entity,
			Columns: NewFieldSpecs(rows...),
		},
		Sets: make([]map[string][]CaseSpec, 0),
	}
}

// NewUpdate 生成更新语句，并执行。
// 如果执行失败，会回滚事务。
//
// Params:
//
//   - ctx: 上下文。
//   - drv: 数据库连接。
//   - spec: 更新语句的信息。
func NewUpdate(ctx context.Context, drv dialect.Tx, spec *UpdateSpec) error {
	builder := NewDialect(drv.Dialect())
	qb := updateBuilder{UpdateSpec: spec, entityBuilder: entityBuilder{builder: builder}}
	return qb.update(ctx, drv)
}

// updateBuilder 更新语句生成器。
type updateBuilder struct {
	entityBuilder
	*UpdateSpec
}

// update 生成更新语句，并执行。
//
// Params:
//
//   - ctx: 上下文。
//   - drv: 数据库连接。
func (b *updateBuilder) update(ctx context.Context, drv dialect.Tx) error {
	updater, err := b.updater(ctx)
	if err != nil {
		return err
	}
	specs, err := updater.Query()
	if err != nil {
		return err
	}
	for _, spec := range specs {
		config := entity.GetConfig()
		if *(config.SqlConsole) {
			tlog.Debug(*config.SqlLogger, fmt.Sprintf("sql: %s", spec.Query))
		}
		var rows dialect.Rows
		if err := drv.Query(ctx, spec.Query, spec.Args, &rows); err != nil {
			return err
		}
		for rows.Next() {
			scanneeFields := make([]ScannerField, len(b.Entity.Columns))
			for i, column := range b.Entity.Columns {
				scanneeFields[i] = ScannerField(column.Name)
			}
			err := b.Scan(rows, scanneeFields)
			if err != nil {
				return err
			}
		}
		// 避免出现pq: unexpected Parse response 'C'
		rows.Close()
	}
	return nil
}

// updater 生成更新语句。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//		0: 更新语句生成器。
//	 1: 错误信息。
func (b *updateBuilder) updater(ctx context.Context) (*Updater, error) {
	updater := NewUpdater(ctx)
	// t := b.entityBuilder.builder.Table(b.Entity.Name)
	updater.SetDialect(b.builder.dialect)
	updater.SetEntity(b.Entity.Name)
	for row, cs := range b.Sets {
		updater.AddBatch()
		for column, c := range cs {
			updater.Set(row, column, c)
		}
		pred := b.Predicate[row]
		if pred != nil {
			w := P(updater.Builder)
			pred(w)
			updater.wheres = append(updater.wheres, w)
		} else {
			updater.wheres = append(updater.wheres, nil)
		}
	}
	return updater, nil
}
