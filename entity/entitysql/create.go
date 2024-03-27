package entitysql

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/tlog"
)

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
	// TODO
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
			Columns: columns,
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
func (s *CreateSpec) CheckRequired(name FieldName, f entity.FieldStorager) error {
	if f.Value() == nil {
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
			var res sql.Result
			if err := drv.Exec(ctx, spec.Query, spec.Args, &res); err != nil {
				return err
			}
		}
		return nil
	} else {
		return b.insertLastID(ctx, inserter)
	}
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

// insertLastID 用于有RETURNING子句的查询。
// 如果数据库不支持RETURNING子句，会使用LastInsertId，比如MySQL。
//
// Params:
//
//   - ctx: 上下文。
//   - insert: 插入语句生成器。
//
// Returns:
//
//	0: 错误信息。
func (b *createBuilder) insertLastID(ctx context.Context, insert *Inserter) error {
	specs, err := insert.Insert()
	if err != nil {
		return err
	}
	for _, spec := range specs {
		// MySQL 不支持 RETURNING 子句。
		if insert.Dialect() != dialect.MySQL {
			rows := dialect.Rows{}
			tlog.Print(spec.Query)
			tlog.Print(spec.Args)
			if err := b.drv.Query(ctx, spec.Query, spec.Args, &rows); err != nil {
				return err
			}
			for rows.Next() {
				tlog.Print("scan")
				err := b.Scan(rows, b.Returning)
				if err != nil {
					return err
				}
			}
		} else {
			// MySQL.
			var res sql.Result
			if err := b.drv.Exec(ctx, spec.Query, spec.Args, &res); err != nil {
				return err
			}
		}

		// 如果是数字类型，可以使用LastInsertId。
		// 如果没有自增主键会报错，同时只可能有一个自增主键。
		// TODO 还没有完成
		// if c.ID.Type.Numeric() {
		// 	id, err := res.LastInsertId()
		// 	if err != nil {
		// 		return err
		// 	}
		// 	c.ID.Value = id
		// }
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
	for _, fields := range b.Fields {
		err := setColumns(fields, func(column string, value driver.Value) {
			inserter.Set(column, value)
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
