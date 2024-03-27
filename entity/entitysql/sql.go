package entitysql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

// Scanner 用于扫描返回的数据。
type Scanner func(row dialect.Rows, rows []FieldName) error

type (
	// SqlSpec sql语句信息。
	SqlSpec struct {
		Query string
		Args  []any
	}

	// EntitySpec 实体表的信息。
	EntitySpec struct {
		// Name 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		Name string
		// Columns 实体表的列名，通过这个属性会指定Select中包含的列。
		Columns []FieldName
	}
	// FieldSpec 字段信息。
	FieldName string

	// FieldSpec 字段信息。
	FieldSpec struct {
		Column  string
		Value   driver.Value // value to be stored.
		Default bool
	}

	// CaseSpec Case语句信息。
	CaseSpec struct {
		// Value Case的值。
		Value any
		// When Case的条件。
		When func(*Predicate)
	}
)

// String 返回字段的名称。
func (e FieldName) String() string {
	return string(e)
}

// setColumns 设置字段的值。
//
// Params:
//
//   - fields: 字段信息。
//   - set: 设置字段的值。
func setColumns(fields []*FieldSpec, set func(column string, value driver.Value)) error {
	for _, fi := range fields {
		value := fi.Value
		set(fi.Column, value)
	}
	return nil
}

type (
	// entityBuilder 实体表的生成器。
	entityBuilder struct {
		drv     dialect.ExecQuerier
		builder *DialectBuilder
	}

	// DialectBuilder 方言生成器。
	DialectBuilder struct {
		dialect dialect.DbDriver
	}
)

// NewDialect 创建一个方言生成器。
func NewDialect(dialect dialect.DbDriver) *DialectBuilder {
	return &DialectBuilder{dialect: dialect}
}

// NewEntityBuilder 创建一个实体表的生成器。
func (b *DialectBuilder) Select() *Selector {
	s := Selector{}
	s.SetDialect(b.dialect)
	return &s
}

// Rollback 回滚事务。
func Rollback(tx dialect.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}

/**************** Selector 选择语句生成器 ***************/

type (
	// 选择语句生成器
	Selector struct {
		Builder
		ctx          context.Context
		limit        *int
		selectFields []Selection
		from         []string
		where        *Predicate
	}
	// Selection 选择的字段。
	Selection struct {
		name string
	}
)

// Query 生成一个查询语句。
//
// Returns:
//
//	0: 查询语句。
//	1: 查询参数。
func (s *Selector) Query() (SqlSpec, error) {
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
	if len(b.args) > entity.BatchSize {
		return SqlSpec{}, entity.Err_0100030004
	}
	if s.limit != nil {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(*s.limit))
	}
	return SqlSpec{Query: b.String(), Args: b.args}, nil
}

// Rows 返回字段的名称。
//
// Params:
//
//   - rows: 字段的名称。
//
// Returns:
//
//	0: 字段的名称。
func (s *Selector) Rows(rows ...FieldName) []string {
	names := make([]string, len(rows))
	for i := range rows {
		names[i] = string(rows[i])
	}
	return names
}

// SetDialect
//
// Params:
//
//   - dialect: 数据库方言。
//
// Returns:
//
//	0: 选择语句生成器。
func (s *Selector) SetDialect(dialect dialect.DbDriver) *Selector {
	s.dialect = dialect
	return s
}

// SetContext 设置上下文。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 选择语句生成器。
func (s *Selector) SetContext(ctx context.Context) *Selector {
	s.ctx = ctx
	return s
}

// SetSelect 设置选择的字段。
//
// Params:
//
//   - rows: 字段的名称。
//
// Returns:
//
//	0: 选择语句生成器。
func (s *Selector) SetSelect(rows ...string) *Selector {
	s.selectFields = make([]Selection, len(rows))
	for i := range rows {
		s.selectFields[i] = Selection{name: rows[i]}
	}
	return s
}

// SetFrom 设置查询的表。
//
// Params:
//
//   - from: 表名。
//
// Returns:
//
//	0: 选择语句生成器。
func (s *Selector) SetFrom(from string) *Selector {
	s.from = append(s.from, from)
	return s
}

// SetLimit 设置查询的限制。
//
// Params:
//
//   - limit: 限制条数。
//
// Returns:
//
//	0: 选择语句生成器。
func (s *Selector) SetLimit(limit int) *Selector {
	s.limit = &limit
	return s
}

// appendSelect 添加选择的字段。
//
// Params:
//
//   - b: sql生成器。
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
	// Inserter 插入语句生成器。
	Inserter struct {
		Builder
		ctx context.Context
		// schema 模式名称。
		schema string
		// entity 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		entity   string
		columns  []string
		values   map[string][]any
		rowTotal int
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

// NewInserter 创建一个插入语句生成器。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 插入语句生成器。
func NewInserter(ctx context.Context) *Inserter {
	return &Inserter{
		ctx: ctx,
	}
}

// Insert 生成插入语句。
//
// Returns:
//
//	0: 插入语句。
//	1: 错误信息。
func (i *Inserter) Insert() ([]SqlSpec, error) {
	specs := []SqlSpec{}
	current := 0
	b := i.Builder.new()
	if i.defaults && len(i.columns) == 0 {
		i.writeDefault(&b)
	} else {
		i.setInitialQuery(&b)
		b.WriteByte('(').IdentComma(i.columns...).WriteByte(')')
		b.WriteString(" VALUES ")
		for j := 0; j < i.rowTotal; j++ {
			if current+len(i.columns) > entity.BatchSize {
				specs = append(specs, SqlSpec{Query: b.String(), Args: b.args})
				b = i.Builder.new()
				i.setInitialQuery(&b)
				b.WriteByte('(').IdentComma(i.columns...).WriteByte(')')
				b.WriteString(" VALUES ")
				current = 0
			}
			if current > 0 {
				b.Comma()
			}
			v := []any{}
			for _, column := range i.columns {
				v = append(v, i.values[column][j])
			}
			b.WriteByte('(').Args(v...).WriteByte(')')
			current += len(i.columns)
		}
	}
	// if i.conflict != nil {
	// 	i.writeConflict(&b)
	// }
	joinReturning(&b, i.returning)
	specs = append(specs, SqlSpec{Query: b.String(), Args: b.args})

	return specs, nil
}

// setInitialQuery 设置初始的插入语句。
//
// Params:
//
//   - b: sql生成器。
func (i *Inserter) setInitialQuery(b *Builder) {
	b.WriteString("INSERT INTO ")
	b.WriteSchema(i.schema)
	b.Ident(i.entity).Blank()
}

// Set 设置插入的值。
//
// Params:
//
//   - column: 列名,也就是字段名。
//   - v: 值。
func (i *Inserter) Set(column string, v any) *Inserter {
	if i.values == nil {
		i.values = make(map[string][]any)
	}
	i.values[column] = append(i.values[column], v)
	return i
}

// AddRow 添加一行数据。
//
// Returns:
//
//	0: 插入语句生成器。
func (i *Inserter) AddRow() *Inserter {
	i.rowTotal += 1
	return i
}

// FillDefault 用于批量插入，当前行的数据有一些列没有值，填充Null，这里不做非空字段的判断，
// 由codegen生成的代码保证。
//
// Returns:
//
//	0: 插入语句生成器。
func (i *Inserter) FillDefault() *Inserter {
	for _, column := range i.columns {
		if _, ok := i.values[column]; !ok {
			i.values[column] = make([]any, 0, i.rowTotal)
		}
		for len(i.values[column]) < i.rowTotal {
			i.values[column] = append(i.values[column], IdentDefault)
		}
	}
	return i
}

// SetSchema 设置模式名称。
//
// Params:
//
//   - schema: 模式名称。
//
// Returns:
//
//	0: 插入语句生成器。
func (i *Inserter) SetSchema(schema string) *Inserter {
	i.schema = schema
	return i
}

// SetEntity 设置实体表的名称。
//
// Params:
//
//   - entity: 实体表的名称。
//
// Returns:
//
//	0: 插入语句生成器。
func (i *Inserter) SetEntity(entity string) *Inserter {
	i.entity = entity
	return i
}

// SetColumns 设置列名。
//
// Params:
//
//   - columns: 列名。
//
// Returns:
//
//	0: 插入语句生成器。
func (i *Inserter) SetColumns(columns ...FieldName) *Inserter {
	i.columns = make([]string, len(columns))
	for j := range columns {
		i.columns[j] = string(columns[j])
	}
	return i
}

// SetReturning 设置返回的列名。
//
// Params:
//
//   - returning: 返回的列名。
//
// Returns:
//
//	0: 插入语句生成器。
func (i *Inserter) SetReturning(returning ...FieldName) *Inserter {
	i.returning = returning
	return i
}

// SetDialect 设置数据库方言。
//
// Params:
//
//   - dialect: 数据库方言。
//
// Returns:
//
//	0: 插入语句生成器。
func (i *Inserter) SetDialect(dialect dialect.DbDriver) *Inserter {
	i.dialect = dialect
	return i
}

// writeDefault 写入默认值。
//
// Params:
//
//   - b: sql生成器。
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
	// Updater 更新语句生成器。
	Updater struct {
		Builder
		// ctx 上下文。
		ctx context.Context
		// schema 模式名称。
		schema string
		// entity 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		entity string
		// columns 列名。
		columns [][]string
		// values 更新的列的值。
		values [][][]CaseSpec
		// wheres 更新的条件。
		wheres []*Predicate
		// returning 返回的列名。
		returning []FieldName
		// batchNum 批量更新的数量，这个用于防止参数过多。
		batchNum int
	}
)

// NewUpdater 创建一个更新语句生成器。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 更新语句生成器。
func NewUpdater(ctx context.Context) *Updater {
	updater := &Updater{
		ctx: ctx,
	}
	return updater
}

// Query 生成更新语句。
//
// Returns:
//
//	0: 更新语句。
//	1: 错误信息。
func (u *Updater) Query() ([]SqlSpec, error) {
	specs := []SqlSpec{}
	for i := 0; i < u.batchNum; i++ {
		b := u.Builder.new()
		b.WriteString("UPDATE ")
		b.WriteSchema(u.schema)
		b.Ident(u.entity).WriteString(" SET ")
		u.writeSetter(&b, i)
		if u.wheres[i] != nil {
			b.WriteString(" WHERE ")
			w := u.wheres[i]
			b.Join(w)
		}
		joinReturning(&b, u.returning)
		specs = append(specs, SqlSpec{Query: b.String(), Args: b.args})
	}
	return specs, nil
}

// AddBatch 添加一个批量更新。
//
// Returns:
//
//	0: 更新语句生成器。
func (u *Updater) AddBatch() *Updater {
	u.batchNum++
	return u
}

// Set 设置更新的值。
//
// Params:
//
//   - row: 行数。
//   - column: 列名。
//   - cs: 更新的值。
//
// Returns:
//
//	0: 更新语句生成器。
func (u *Updater) Set(row int, column string, cs []CaseSpec) *Updater {
	if len(u.columns)-1 <= row {
		u.columns = append(u.columns, []string{})
		u.values = append(u.values, [][]CaseSpec{})
	}
	u.columns[row] = append(u.columns[row], column)
	u.values[row] = append(u.values[row], cs)
	return u
}

// SetEntity 设置实体表的名称。
//
// Params:
//
//   - entity: 实体表的名称。
//
// Returns:
//
//	0: 更新语句生成器。
func (u *Updater) SetEntity(entity string) *Updater {
	u.entity = entity
	return u
}

// SetDialect 设置数据库方言。
//
// Params:
//
//   - dialect: 数据库方言。
//
// Returns:
//
//	0: 更新语句生成器。
func (u *Updater) SetDialect(dialect dialect.DbDriver) *Updater {
	u.dialect = dialect
	return u
}

// writeSetter 写入更新的值。
//
// Params:
//
//   - b: sql生成器。
//   - row: 行数。
func (u *Updater) writeSetter(b *Builder, row int) {
	for i, column := range u.columns[row] {
		if i > 0 {
			b.Comma()
		}
		b.Ident(column).WriteString(" = ")
		b.Join(NewCaser(column, u.values[row][i]))
	}
}

/**************** Deleter 删除语句生成器 ***************/

type (
	// Deleter 删除语句生成器。
	Deleter struct {
		Builder
		ctx context.Context
		// schema 模式名称。
		schema string
		// entity 实体表的名称，和[entity.EntityConfig]中的AttrName相同。
		entity string
		// where 删除的条件。
		where *Predicate
	}
)

// NewDeleter 创建一个删除语句生成器。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 删除语句生成器。
func NewDeleter(ctx context.Context) *Deleter {
	return &Deleter{
		ctx: ctx,
	}
}

// Query 生成删除语句。
//
// Returns:
//
//	0: 删除语句。
//	1: 错误信息。
func (d *Deleter) Query() ([]SqlSpec, error) {
	specs := []SqlSpec{}
	b := d.Builder.new()
	d.setInitialQuery(&b)
	if d.where != nil {
		length := d.where.FunsLen()
		batchNum := length/entity.BatchSize + 1
		for i := 0; i < batchNum; i++ {
			b.WriteString(" WHERE ")
			if length < (i+1)*entity.BatchSize {
				b.Join(d.where.Clone(i*entity.BatchSize, -1))
			} else {
				b.Join(d.where.Clone(i*entity.BatchSize, (i+1)*entity.BatchSize))
			}
			specs = append(specs, SqlSpec{Query: b.String(), Args: b.args})
			b = d.Builder.new()
			d.setInitialQuery(&b)
		}
	}
	return specs, nil
}

// setInitialQuery 设置初始的删除语句。
func (d *Deleter) setInitialQuery(b *Builder) {
	b.WriteSchema(d.schema)
	b.Ident(d.entity).Blank()
}

// SetEntity 设置实体表的名称。
//
// Params:
//
//   - entity: 实体表的名称。
//
// Returns:
//
//	0: 删除语句生成器。
func (d *Deleter) SetEntity(entity string) *Deleter {
	d.entity = entity
	return d
}

// SetDialect 设置数据库方言。
//
// Params:
//
//   - dialect: 数据库方言。
//
// Returns:
//
//	0: 删除语句生成器。
func (d *Deleter) SetDialect(dialect dialect.DbDriver) *Deleter {
	d.dialect = dialect
	return d
}

// SetSchema 设置模式名称。
//
// Params:
//
//   - schema: 模式名称。
//
// Returns:
//
//	0: 删除语句生成器。
func (d *Deleter) SetSchema(schema string) *Deleter {
	d.schema = schema
	return d
}

/**************** Predicate Where子句生成器 ***************/

type (
	// Predicate Where子句生成器。
	Predicate struct {
		Builder
		// fns Where子句生成器的函数。
		fns []func(*Builder)
	}
	// PredicateFunc Where子句生成器的函数。
	PredicateFunc func(p *Predicate)
)

// P 创建一个Where子句生成器。
func P(fns ...func(*Builder)) *Predicate {
	return &Predicate{fns: fns}
}

// Clone 复制一个Predicate
//
// Params:
//
//   - begin: 复制从第几个开始， -1表示没有开始位置
//   - end: 复制到第几个结束， -1表示没有结束位置
//
// Returns:
//
//	0: 复制的Predicate。
func (p *Predicate) Clone(begin int, end int) *Predicate {
	if begin < 0 {
		return &Predicate{fns: p.fns[:end]}
	} else if end < 0 {
		return &Predicate{fns: p.fns[begin:]}
	} else {
		return &Predicate{fns: p.fns[begin:end]}
	}
}

// Query 生成查询中的Where子句。
//
// Returns:
//
//	0: 查询中的Where子句。
//	1: 查询中的参数。
func (p *Predicate) Query() (string, []any) {
	if p.Len() > 0 || len(p.args) > 0 {
		p.Reset()
		p.args = nil
	}
	for i, f := range p.fns {
		// 添加这个判断是因为在批量操作中，如果参数过多会执行分批操作，此时可能出现第一个或最后一个是AND或者OR的情况，
		if i == 0 || i == len(p.fns)-1 {
			nb := p.Builder.new()
			f(&nb)
			s := nb.String()
			if strings.Contains(s, "AND") || strings.Contains(s, "OR") {
				continue
			}

		}
		f(&p.Builder)
	}
	return p.String(), p.args
}

// FunsLen 返回Where子句的长度。
//
// Returns:
//
//	0: Where子句的长度。
func (p *Predicate) FunsLen() int {
	return len(p.fns)
}

// Append 添加一个Where子句。
//
// Params:
//
//   - f: Where子句生成器的函数。
func (p *Predicate) Append(f func(*Builder)) *Predicate {
	p.fns = append(p.fns, f)
	return p
}

// And 添加一个AND标识符。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) And() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString(" AND ")
	})
}

// Or 添加一个OR标识符。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Or() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString(" OR ")
	})
}

// Not 添加一个NOT标识符。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Not() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString(" NOT ")
	})
}

// EQ 添加一个等于的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) EQ(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpEQ)
		p.arg(b, v)
		b.Blank()
	})
}

// NEQ 添加一个不等于的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) NEQ(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpNEQ)
		p.arg(b, v)
		b.Blank()
	})
}

// GT 添加一个大于的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) GT(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpGT)
		p.arg(b, v)
		b.Blank()
	})
}

// GTE 添加一个大于等于的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) GTE(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpGTE)
		p.arg(b, v)
		b.Blank()
	})
}

// LT 添加一个小于的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) LT(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpLT)
		p.arg(b, v)
		b.Blank()
	})
}

// LTE 添加一个小于等于的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) LTE(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpLTE)
		p.arg(b, v)
		b.Blank()
	})
}

// In 添加一个IN的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) In(column string, v ...any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpIn)
		b.WriteByte('(')
		for i, a := range v {
			if i > 0 {
				b.Comma()
			}
			p.arg(b, a)
		}
		b.WriteByte(')')
	})
}

// NotIn 添加一个NOT IN的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) NotIn(column string, v ...any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpNotIn)
		b.WriteByte('(')
		for i, a := range v {
			if i > 0 {
				b.Comma()
			}
			p.arg(b, a)
		}
		b.WriteByte(')')
	})
}

// Like 添加一个LIKE的条件。
//
// Params:
//
//   - column: 列名。
//   - v: 值。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Like(column string, v any) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpLike)
		p.arg(b, v)
		b.Blank()
	})
}

// IsNull 添加一个IS NULL的条件。
//
// Params:
//
//   - column: 列名。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) IsNull(column string) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpIsNull)
		b.Blank()
	})
}

// NotNull 添加一个IS NOT NULL的条件。
//
// Params:
//
//   - column: 列名。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) NotNull(column string) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(column)
		b.WriteOp(OpNotNull)
		b.Blank()
	})
}

// Add 添加一个加法的条件。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Add() *Predicate {
	return p.Append(func(b *Builder) {
		b.Blank()
		b.WriteOp(OpAdd)
		b.Blank()
	})
}

// Sub 添加一个减法的条件。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Sub() *Predicate {
	return p.Append(func(b *Builder) {
		b.Blank()
		b.WriteOp(OpSub)
		b.Blank()
	})
}

// Mul 添加一个乘法的条件。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Mul() *Predicate {
	return p.Append(func(b *Builder) {
		b.Blank()
		b.WriteOp(OpMul)
		b.Blank()
	})
}

// Div 添加一个除法的条件。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Div() *Predicate {
	return p.Append(func(b *Builder) {
		b.Blank()
		b.WriteOp(OpDiv)
		b.Blank()
	})
}

// Mod 添加一个取模的条件。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) Mod() *Predicate {
	return p.Append(func(b *Builder) {
		b.Blank()
		b.WriteOp(OpMod)
		b.Blank()
	})
}

// arg 添加一个参数。
//
// Params:
//
//   - b: sql生成器。
//   - a: 参数。
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

/**************** CASE 批量更新中Set语句生成器 ***************/

type (
	// Caser CASE语句生成器。
	Caser struct {
		Builder
		// Cases CASE语句信息。
		Cases  []CaseSpec
		column string
		elser  any
	}
)

// NewCaser 创建一个CASE语句生成器。
//
// Params:
//
//   - column: 列名。
//   - cases: CASE语句信息。
//
// Returns:
//
//	0: CASE语句生成器。
func NewCaser(column string, cases []CaseSpec) *Caser {
	return &Caser{
		Cases:  cases,
		column: column,
	}
}

// Query 生成CASE语句。
//
// Returns:
//
//	0: CASE语句。
//	1: CASE语句的参数。
func (c *Caser) Query() (string, []any) {
	if c.Len() > 0 || len(c.args) > 0 {
		c.Reset()
		c.args = nil
	}
	if len(c.Cases) == 0 {
		return "", nil
	} else if len(c.Cases) == 1 && c.Cases[0].When == nil {
		c.arg(c.Cases[0].Value)
		return c.String(), c.args
	} else {
		c.WriteString("CASE ")
		for _, cs := range c.Cases {
			// 对于多个没有条件的CASE，只有最后一个会生效，
			// 如果要做限制或者判断,请在外部实现。
			if cs.When == nil {
				c.elser = cs.Value
			} else {
				c.When(cs.When)
				c.Then(cs.Value)
			}
		}
		c.Else()
		c.WriteString(" END")
	}
	return c.String(), c.args
}

// When 添加一个WHEN条件。
//
// Params:
//
//   - pred: WHERE子句生成器。
//
// Returns:
//
//	0: CASE语句生成器。
func (c *Caser) When(pred func(*Predicate)) *Caser {
	c.WriteString(" WHEN ")
	p := P()
	pred(p)
	c.Join(p)
	return c
}

// Then 添加一个THEN条件。
//
// Params:
//
//   - v: 值。
//
// Returns:
//
//	0: CASE语句生成器。
func (c *Caser) Then(v any) *Caser {
	c.WriteString(" THEN ")
	c.arg(v)
	return c
}

// Else 添加一个ELSE条件。
//
// Returns:
//
//	0: CASE语句生成器。
func (c *Caser) Else() *Caser {
	c.WriteString(" ELSE ")
	if c.elser == nil {
		c.Ident(c.column)
	} else {
		c.arg(c.elser)
	}
	return c
}

// arg 添加一个参数。
//
// Params:
//
//   - a: 参数。
func (c *Caser) arg(a any) {
	switch a.(type) {
	case *Selector:
		c.Builder.Wrap(func(b *Builder) {
			b.Arg(a)
		})
	default:
		c.Builder.Arg(a)
	}
}
