package entitysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/tlog"
)

// DeleteSpec 用于生成删除语句。
type DeleteSpec struct {
	Entity    *EntitySpec
	Predicate PredicateFunc
	Affected  *int64
}

// NewDeleteSpec 创建一个DeleteSpec。
//
// Params:
//
//   - entity: 实体的名称。
func NewDeleteSpec(entity string) *DeleteSpec {
	return &DeleteSpec{
		Entity: &EntitySpec{
			Name: entity,
		},
	}
}

// NewDelete 生成删除语句，并执行。
// 如果执行失败，会回滚事务。
//
// Params:
//
//   - ctx: 上下文。
//   - drv: 数据库连接。
func NewDelete(ctx context.Context, drv dialect.Tx, spec *DeleteSpec) error {
	builder := NewDialect(drv.Dialect())
	qb := deleteBuilder{DeleteSpec: spec, entityBuilder: entityBuilder{builder: builder}}
	return qb.delete(ctx, drv)
}

// deleteBuilder 删除语句生成器。
type deleteBuilder struct {
	entityBuilder
	*DeleteSpec
}

// delete 生成删除语句，并执行。
//
// Params:
//
//   - ctx: 上下文。
//   - drv: 数据库连接。
func (b *deleteBuilder) delete(ctx context.Context, drv dialect.Tx) error {
	deleter, err := b.deleter(ctx)
	if err != nil {
		return err
	}
	specs, err := deleter.Query()
	if err != nil {
		return err
	}
	for _, spec := range specs {
		config := entity.GetConfig()
		if *(config.SqlConsole) {
			tlog.Debug(*config.SqlLogger, fmt.Sprintf("sql: %s", spec.Query))
		}
		var res sql.Result
		if err := drv.Exec(ctx, spec.Query, spec.Args, &res); err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		*b.Affected += affected
	}
	return nil
}

// deleter 生成删除语句。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 删除语句生成器。
func (b *deleteBuilder) deleter(ctx context.Context) (*Deleter, error) {
	deleter := NewDeleter(ctx)
	// t := b.entityBuilder.builder.Table(b.Entity.Name)
	deleter.SetDialect(b.builder.dialect)
	deleter.SetEntity(b.Entity.Name)
	if pred := b.Predicate; pred != nil {
		deleter.where = P(deleter.Builder)
		pred(deleter.where)
	}
	return deleter, nil
}
