package entitysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/dialect"
	"github.com/zodileap/taurus_go/tlog"
)

var ErrReturningUnsupported = errors.New("entitysql RETURNING is not supported by this dialect")

// CreateSpec 用于生成创建语句。
type CreateSpec struct {
	// Entity 表示实体的信息。
	Entity *EntitySpec
	// Scan 用于扫描返回的数据。
	Scan Scanner
	// Fields 用于生成创建语句。
	Fields [][]*FieldSpec
	// Returning 用于返回的字段。
	Returning []FieldName
	// Schema 可选 schema 名称，不设置时使用数据库默认 schema。
	Schema string
}

// NewCreateSpec 创建一个CreateSpec。
//
// Params:
//
//   - entity: 实体的名称。
//   - columns: 实体的字段。
func NewCreateSpec(entity string, columns []FieldName) *CreateSpec {
	return &CreateSpec{
		Entity: &EntitySpec{
			Name:    entity,
			Columns: NewFieldSpecs(columns...),
		},
	}
}

// NewCreate 生成创建语句，并执行。
// 如果执行失败，会回滚事务。
//
// Params:
//
//   - ctx: 上下文。
//   - drv: 数据库连接。
//   - spec: 创建语句的信息。
func NewCreate(ctx context.Context, drv dialect.Tx, spec *CreateSpec) error {
	builder := NewDialect(drv.Dialect())
	entity := entityBuilder{builder: builder, drv: drv}
	qb := createBuilder{CreateSpec: spec, entityBuilder: entity}
	return qb.create(ctx, drv)
}

// CheckRequired 检查字段是否为空。
//
// Params:
//
//   - name: 字段名称。
//   - f: 字段。
func (s *CreateSpec) CheckRequired(dbDriver dialect.DbDriver, name FieldName, f entity.FieldStorager) error {
	v, err := f.SqlParam(dbDriver)
	if v == nil || err != nil {
		return entity.Err_0100030001.Sprintf(s.Entity.Name, name)
	}
	return nil
}

// createBuilder 用于生成创建语句。
type createBuilder struct {
	entityBuilder
	*CreateSpec
}

// create 生成创建语句，并执行。
// 如果执行失败，会回滚事务。
//
// Params:
//
//   - ctx: 上下文。
//   - drv: 数据库连接。
func (b *createBuilder) create(ctx context.Context, drv dialect.Tx) error {
	inserter, err := b.inserter(ctx)
	if err != nil {
		return err
	}
	if len(inserter.returning) == 0 {
		specs, err := inserter.Insert()
		if err != nil {
			return err
		}
		for _, spec := range specs {
			config := entity.GetConfig()
			if *(config.SqlConsole) {
				tlog.Debug(*config.SqlLogger, fmt.Sprintf("sql: %s", spec.Query))
				tlog.Debug(*config.SqlLogger, fmt.Sprintf("args: %v", spec.Args))
			}
			var res sql.Result
			if err := drv.Exec(ctx, spec.Query, spec.Args, &res); err != nil {
				return err
			}
		}
		return nil
	}

	return b.insertReturning(ctx, inserter)
}

// inserter 从CreateSpec提取Inserter生成sql需要的信息。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 插入语句生成器。
//	1: 错误信息。
func (b *createBuilder) inserter(ctx context.Context) (*Inserter, error) {
	inserter := NewInserter(ctx)
	inserter.SetDialect(b.builder.dialect)
	inserter.SetSchema(b.Schema).SetEntity(b.Entity.Name)
	err := b.setColumns(inserter)
	if err != nil {
		return nil, err
	}
	inserter.SetReturning(b.Returning...)
	return inserter, nil
}

// insertReturning 用于执行包含 RETURNING 字段的插入。
//
// Params:
//
//   - ctx: 上下文。
//   - insert: 插入语句生成器。
//
// Returns:
//
//	0: 错误信息。
func (b *createBuilder) insertReturning(ctx context.Context, insert *Inserter) error {
	if insert.Dialect() == dialect.MySQL {
		return fmt.Errorf("%w: %s", ErrReturningUnsupported, insert.Dialect())
	}

	specs, err := insert.Insert()
	if err != nil {
		return err
	}
	for _, spec := range specs {
		rows := dialect.Rows{}
		config := entity.GetConfig()
		if *(config.SqlConsole) {
			tlog.Debug(*config.SqlLogger, fmt.Sprintf("sql: %s", spec.Query))
			tlog.Debug(*config.SqlLogger, fmt.Sprintf("args: %v", spec.Args))
		}
		if err := b.drv.Query(ctx, spec.Query, spec.Args, &rows); err != nil {
			return err
		}
		for rows.Next() {
			scannerFields := make([]ScannerField, len(b.Returning))
			for i, name := range b.Returning {
				scannerFields[i] = ScannerField(name)
			}
			if err := b.Scan(rows, scannerFields); err != nil {
				return err
			}
		}
		if err := rows.Close(); err != nil {
			return err
		}
	}
	return nil
}

// setColumns 用于设置插入的字段。
//
// Params:
//
//   - inserter: 插入语句生成器。
func (b *createBuilder) setColumns(inserter *Inserter) error {
	inserter.SetColumns(b.Entity.Columns...)
	t := b.entityBuilder.builder.Table(b.Entity.Name)
	for _, fields := range b.Fields {
		err := setColumns(fields, func(column string, field FieldSpec) {
			if field.Param != nil {
				inserter.Set(column, t.as, field)
			}

		})
		if err != nil {
			return err
		} else {
			inserter.AddRow()
			// 请不要把FillDefault放在AddRow前面，不然会导致填充失败
			inserter.FillDefault()
		}
	}
	return nil
}
