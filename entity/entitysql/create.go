package entitysql

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type CreateSpec struct {
	Entity    *EntitySpec
	Scan      Scanner
	Fields    []*FieldSpec
	Returning []FieldName
	// TODO
	Schema string
}

func NewCreateSpec(entity string, rows []FieldName) *CreateSpec {
	return &CreateSpec{
		Entity: &EntitySpec{
			Name: entity,
			Rows: rows,
		},
	}
}

func NewCreate(ctx context.Context, drv dialect.Driver, spec *CreateSpec) error {
	builder := NewDialect(drv.Dialect())
	entity := entityBuilder{builder: builder, drv: drv}
	qb := createBuilder{CreateSpec: spec, entityBuilder: entity}
	return qb.create(ctx, drv)
}

// CheckRequired 检查字段是否为空。
func (s *CreateSpec) CheckRequired(name FieldName, f entity.FieldStorager) error {
	if f.Value() == nil {
		return entity.Err_0100030001.Sprintf(s.Entity.Name, name)
	}
	return nil
}

type createBuilder struct {
	entityBuilder
	*CreateSpec
}

// create 生成创建语句，并执行。
// 如果执行失败，会回滚事务。
func (b *createBuilder) create(ctx context.Context, drv dialect.Driver) error {
	var res sql.Result
	inserter, err := b.inserter(ctx)
	if err != nil {
		return err
	}
	tx, err := b.mayTx(ctx, drv)
	if err := func() error {
		if len(inserter.returning) == 0 {
			Insert, args := inserter.Insert()
			return drv.Exec(ctx, Insert, args, res)
		} else {
			return b.insertLastID(ctx, inserter)
		}
	}(); err != nil {
		return rollback(tx, err)
	}
	return tx.Commit()
}

// inserter 从CreateSpec提取Inserter生成sql需要的信息。
func (b *createBuilder) inserter(ctx context.Context) (*Inserter, error) {
	inserter := Inserter{}
	inserter.SetDialect(b.builder.dialect)
	inserter.SetSchema(b.Schema).SetEntity(b.Entity.Name)
	err := b.setColumns(&inserter)
	if err != nil {
		return nil, err
	}
	inserter.SetReturning(b.Returning...)
	return &inserter, nil
}

// insertLastID 用于有RETURNING子句的查询。
// 如果数据库不支持RETURNING子句，会使用LastInsertId，比如MySQL。
func (b *createBuilder) insertLastID(ctx context.Context, insert *Inserter) error {
	query, args := insert.Insert()
	// MySQL 不支持 RETURNING 子句。
	if insert.Dialect() != dialect.MySQL {
		rows := dialect.Rows{}
		if err := b.drv.Query(ctx, query, args, &rows); err != nil {
			return err
		}
		for rows.Next() {
			err := b.Scan(rows, b.Entity.Rows)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// MySQL.
	var res sql.Result
	if err := b.drv.Exec(ctx, query, args, &res); err != nil {
		return err
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
	return nil
}

func (b *createBuilder) setColumns(inserter *Inserter) error {
	err := setColums(b.Fields, func(column string, value driver.Value) {
		inserter.Set(column, value)
	})
	return err
}

// mayTx 打开一个新事务
func (b *createBuilder) mayTx(ctx context.Context, drv dialect.Driver) (dialect.Tx, error) {
	tx, err := drv.Tx(ctx)
	if err != nil {
		return nil, err
	}
	b.drv = tx
	return tx, nil
}
