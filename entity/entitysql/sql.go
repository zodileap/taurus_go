package entitysql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/yohobala/taurus_go/entity/dialect"
)

type Scanner func(row dialect.Rows, rows []FieldName) error

type (
	EntitySpec struct {
		// Name 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		Name string
		// Columns 实体表的列名，通过这个属性会指定Select中包含的列。
		Rows []FieldName
	}
	FieldName string

	FieldSpec struct {
		Column string
		Value  driver.Value // value to be stored.
	}
)

func (e FieldName) String() string {
	return string(e)
}

func setColums(fields []*FieldSpec, set func(column string, value driver.Value)) error {
	for _, fi := range fields {
		value := fi.Value
		set(fi.Column, value)
	}
	return nil
}

type (
	entityBuilder struct {
		drv     dialect.ExecQuerier
		builder *DialectBuilder
	}

	DialectBuilder struct {
		dialect dialect.DbDriver
	}
)

func NewDialect(dialect dialect.DbDriver) *DialectBuilder {
	return &DialectBuilder{dialect: dialect}
}

func (b *DialectBuilder) Select() *Selector {
	s := Selector{}
	s.SetDialect(b.dialect)
	return &s
}

func rollback(tx dialect.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}

/**************** Selector 选择语句生成器 ***************/

type (
	Selector struct {
		Builder
		ctx          context.Context
		limit        *int
		selectFields []Selection
		from         []string
		where        *Predicate
	}
	Selection struct {
		name string
	}
)

func (s *Selector) Query() (string, []any) {
	b := s.Builder.clone()
	b.WriteString("SELECT ")
	if len(s.selectFields) > 0 {
		s.appendSelect(&b)
	} else {
		b.WriteString(" * ")
	}
	if len(s.from) > 0 {
		b.WriteString(" FROM ")
	}
	for i, from := range s.from {
		if i > 0 {
			b.Comma()
		}
		b.WriteString(fmt.Sprintf(`"%s"`, from))
	}
	if s.where != nil {
		b.WriteString(" WHERE ")
		b.Join(s.where)
	}
	if s.limit != nil {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(*s.limit))
	}
	return b.String(), b.args
}

func (s *Selector) Rows(rows ...FieldName) []string {
	names := make([]string, len(rows))
	for i := range rows {
		names[i] = string(rows[i])
	}
	return names
}

func (s *Selector) SetDialect(dialect dialect.DbDriver) *Selector {
	s.dialect = dialect
	return s
}

func (s *Selector) SetContext(ctx context.Context) *Selector {
	s.ctx = ctx
	return s
}

func (s *Selector) SetSelect(rows ...string) *Selector {
	s.selectFields = make([]Selection, len(rows))
	for i := range rows {
		s.selectFields[i] = Selection{name: rows[i]}
	}
	return s
}

func (s *Selector) SetFrom(from string) *Selector {
	s.from = append(s.from, from)
	return s
}

func (s *Selector) SetLimit(limit int) *Selector {
	s.limit = &limit
	return s
}

func (s *Selector) appendSelect(b *Builder) {
	for i, field := range s.selectFields {
		if i > 0 {
			b.Comma()
		}
		b.WriteString(b.Quote(field.name))
	}
}

/**************** Inserter 插入语句生成器 ***************/

type (
	Inserter struct {
		Builder
		ctx context.Context
		// schema 模式名称。
		schema string
		// entity 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		entity  string
		columns []string
		values  [][]any
		// returning 返回的列名。
		returning []FieldName
		// defaults 是否使用默认值。
		defaults bool

		// 定义当尝试插入的数据发生当发生冲突时的处理方式。
		// conflict *conflict
	}
	// conflict struct {
	// 	target struct {
	// 		// constraint 约束名称。
	// 		constraint string
	// 		// columns 需要检查冲突的列名。
	// 		columns []string
	// 		// where 检查冲突的条件。
	// 		where *Predicate
	// 	}
	// 	action struct {
	// 		// nothing 是否忽略冲突。
	// 		nothing bool
	// 		// where 执行更新操作的条件。
	// 		where *Predicate
	// 		// update 用于构建在冲突发生时要执行的更新操作。
	// 		update []func(*UpdateSet)
	// 	}
	// }
)

func (i *Inserter) Insert() (string, []any) {
	b := i.Builder.clone()
	b.WriteString("INSERT INTO ")
	b.WriteSchema(i.schema)
	b.Ident(i.entity).Blank()
	if i.defaults && len(i.columns) == 0 {
		i.writeDefault(&b)
	} else {
		b.WriteByte('(').IdentComma(i.columns...).WriteByte(')')
		b.WriteString(" VALUES ")
		for j, v := range i.values {
			if j > 0 {
				b.Comma()
			}
			b.WriteByte('(').Args(v...).WriteByte(')')
		}
	}
	// if i.conflict != nil {
	// 	i.writeConflict(&b)
	// }
	joinReturning(&b, i.returning)
	return b.String(), b.args
}

func (i *Inserter) Set(column string, v any) *Inserter {
	i.columns = append(i.columns, column)
	if len(i.values) == 0 {
		i.values = append(i.values, []any{v})
	} else {
		i.values[0] = append(i.values[0], v)
	}
	return i
}

func (i *Inserter) SetSchema(schema string) *Inserter {
	i.schema = schema
	return i
}

func (i *Inserter) SetEntity(entity string) *Inserter {
	i.entity = entity
	return i
}

func (i *Inserter) SetReturning(returning ...FieldName) *Inserter {
	i.returning = returning
	return i
}

func (i *Inserter) SetDialect(dialect dialect.DbDriver) *Inserter {
	i.dialect = dialect
	return i
}

func (i *Inserter) writeDefault(b *Builder) {
	switch i.Dialect() {
	case dialect.MySQL:
		b.WriteString("VALUES ()")
	case dialect.PostgreSQL:
		b.WriteString("DEFAULT VALUES")
	}
}

// func (i *Inserter) writeConflict(b *Builder) {
// 	switch i.Dialect() {
// 	case dialect.MySQL:
// 		// 当尝试插入的行在表中已经存在（基于主键或唯一索引）时，更新该行的某些字段
// 		b.WriteString(" ON DUPLICATE KEY UPDATE ")
// 		// Fallback to ResolveWithIgnore() as MySQL
// 		// does not support the "DO NOTHING" clause.
// 		if i.conflict.action.nothing {
// 			i.OnConflict(ResolveWithIgnore())
// 		}
// 	case dialect.PostgreSQL:
// 		b.WriteString(" ON CONFLICT")
// 		switch t := i.conflict.target; {
// 		case t.constraint != "" && len(t.columns) != 0:
// 			b.AddError(fmt.Errorf("duplicate CONFLICT clauses: %q, %q", t.constraint, t.columns))
// 		case t.constraint != "":
// 			b.WriteString(" ON CONSTRAINT ").Ident(t.constraint)
// 		case len(t.columns) != 0:
// 			b.WriteString(" (").IdentComma(t.columns...).WriteByte(')')
// 		}
// 		if p := i.conflict.target.where; p != nil {
// 			b.WriteString(" WHERE ").Join(p)
// 		}
// 		if i.conflict.action.nothing {
// 			b.WriteString(" DO NOTHING")
// 			return
// 		}
// 		b.WriteString(" DO UPDATE SET ")
// 	}
// 	if len(i.conflict.action.update) == 0 {
// 		b.AddError(errors.New("missing action for 'DO UPDATE SET' clause"))
// 	}
// 	u := &UpdateSet{UpdateBuilder: Dialect(i.dialect).Update(i.table), columns: i.columns}
// 	u.Builder = *b
// 	for _, f := range i.conflict.action.update {
// 		f(u)
// 	}
// 	u.writeSetter(b)
// 	if p := i.conflict.action.where; p != nil {
// 		p.qualifier = i.table
// 		b.WriteString(" WHERE ").Join(p)
// 	}
// }

// // Predicate 是一个where的条件。
// type Predicate struct {
// 	Builder
// 	depth int
// 	fns   []func(*Builder)
// }

// // UpdateSet 描述了`DO UPDATE`情况下的更新语句。
// type UpdateSet struct {
// 	*UpdateBuilder
// 	columns []string
// }

/**************** Updater 更新语句生成器 ***************/

type (
	Updater struct {
		Builder
		ctx context.Context
		// schema 模式名称。
		schema string
		// entity 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		entity  string
		columns []string
		values  []any
		where   *Predicate
		// returning 返回的列名。
		returning []FieldName
	}
)

func (u *Updater) Query() (string, []any) {
	b := u.Builder.clone()
	b.WriteString("UPDATE ")
	b.WriteSchema(u.schema)
	b.Ident(u.entity).WriteString(" SET ")
	u.writeSetter(&b)
	if u.where != nil {
		b.WriteString(" WHERE ")
		b.Join(u.where)
	}
	joinReturning(&b, u.returning)
	return b.String(), b.args
}

func (u *Updater) Set(column string, v any) *Updater {
	u.columns = append(u.columns, column)
	u.values = append(u.values, v)
	return u
}

func (u *Updater) SetEntity(entity string) *Updater {
	u.entity = entity
	return u
}

func (u *Updater) SetDialect(dialect dialect.DbDriver) *Updater {
	u.dialect = dialect
	return u
}

func (u *Updater) writeSetter(b *Builder) {
	for i, column := range u.columns {
		if i > 0 {
			b.Comma()
		}
		b.Ident(column).WriteString(" = ")
		switch v := u.values[i].(type) {
		case Querier:
			b.Join(v)
		default:
			b.Arg(v)
		}
	}
}

/**************** Deleter 删除语句生成器 ***************/

type (
	Deleter struct {
		Builder
		ctx context.Context
		// schema 模式名称。
		schema string
		// entity 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		entity string
		where  *Predicate
	}
)

func (d *Deleter) Query() (string, []any) {
	b := d.Builder.clone()
	b.WriteString("DELETE FROM ")
	b.WriteSchema(d.schema)
	b.Ident(d.entity).Blank()
	if d.where != nil {
		b.WriteString(" WHERE ")
		b.Join(d.where)
	}
	return b.String(), b.args
}

func (d *Deleter) SetEntity(entity string) *Deleter {
	d.entity = entity
	return d
}

func (d *Deleter) SetDialect(dialect dialect.DbDriver) *Deleter {
	d.dialect = dialect
	return d
}

func (d *Deleter) SetSchema(schema string) *Deleter {
	d.schema = schema
	return d
}

/**************** Predicate Where子句生成器 ***************/

type (
	Predicate struct {
		Builder
		fns []func(*Builder)
	}
)

func P(fns ...func(*Builder)) *Predicate {
	return &Predicate{fns: fns}
}

func (p *Predicate) Query() (string, []any) {
	if p.Len() > 0 || len(p.args) > 0 {
		p.Reset()
		p.args = nil
	}
	for _, f := range p.fns {
		f(&p.Builder)
	}
	return p.String(), p.args
}

func (p *Predicate) Append(f func(*Builder)) *Predicate {
	p.fns = append(p.fns, f)
	return p
}

func (p *Predicate) And() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString(" AND ")
	})
}

func (p *Predicate) Or() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString(" OR ")
	})
}

func (p *Predicate) Not() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString(" NOT ")
	})
}

func (p *Predicate) EQ(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpEQ)
		p.arg(b, v)
		b.Blank()
	})
}

func (*Predicate) arg(b *Builder, a any) {
	switch a.(type) {
	case *Selector:
		b.Wrap(func(b *Builder) {
			b.Arg(a)
		})
	default:
		b.Arg(a)
	}
}
