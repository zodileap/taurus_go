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
type Scanner func(row dialect.Rows, selects []ScannerField) error

type ScannerField interface {
	String() string
}

type ScannerBuilder struct {
	args [][]any
}

func NewScannerBuilder(length int) *ScannerBuilder {
	return &ScannerBuilder{
		args: make([][]any, length),
	}
}

func (s *ScannerBuilder) Args() [][]any {
	return s.args
}

func (s *ScannerBuilder) Flatten() []any {
	var args []any
	for _, a := range s.args {
		args = append(args, a...)
	}
	return args
}

func (s *ScannerBuilder) Append(index int, args ...any) {
	s.args[index] = append(s.args[index], args...)
}

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
		Columns []FieldSpec
	}
	// FieldSpec 字段信息。
	FieldName string

	// FieldSpec 字段信息。
	FieldSpec struct {
		// Name 字段的名称。
		Name FieldName
		// NameFormat 字段名称的格式化，比如SELECT ST_AsText(geom)中的`ST_AsText(geom)`。
		NameFormat func(dbType dialect.DbDriver, name string) string
		// Param 字段的值。
		Param entity.FieldValue
		// Default 是否使用默认值。
		Default bool
		Format  func(dbType dialect.DbDriver, param string) string
	}

	// CaseSpec Case语句信息。
	CaseSpec struct {
		// Field Case的值。
		Field FieldSpec
		// When Case的条件。
		When PredicateFunc
	}
)

// String 返回字段的名称。
func (e FieldName) String() string {
	return string(e)
}

func NewFieldSpec(column FieldName) FieldSpec {
	return FieldSpec{
		Name: column,
		NameFormat: func(dbType dialect.DbDriver, name string) string {
			return name
		},
	}
}

func NewFieldSpecs(columns ...FieldName) []FieldSpec {
	var fields []FieldSpec
	for _, column := range columns {
		fields = append(fields, NewFieldSpec(column))
	}
	return fields
}

// Value 用于实现driver.Valuer接口。
func (f FieldSpec) Value() (driver.Value, error) {
	return f.Param, nil
}

// FormatParam 格式化字段的值。实现ParamFormatter接口。
func (f FieldSpec) FormatParam(placeholder string, info *StmtInfo) string {
	return f.Format(info.Dialect, placeholder)
}

// setColumns 设置字段的值。
//
// Params:
//
//   - fields: 字段信息。
//   - set: 设置字段的值。
func setColumns(fields []*FieldSpec, set func(column string, value FieldSpec)) error {
	for _, fi := range fields {
		set(string(fi.Name), *fi)
	}
	return nil
}

type TableView interface {
	// view 是一个标记接口，用于标记一个视图。
	view()
	// C 返回序列化的列名。
	C(string) string
}

type SelectTable struct {
	*Builder
	as     string
	asNum  int
	name   string
	schema string
	quote  bool
}

// Table 创建一个表。
func Table(name string) *SelectTable {
	return &SelectTable{
		Builder: &Builder{},
		quote:   true,
		name:    name,
	}
}

// Schema 设置模式名称。
func (s *SelectTable) Schema(name string) *SelectTable {
	s.schema = name
	return s
}

// As 设置别名。
func (s *SelectTable) As(alias string) *SelectTable {
	s.as = alias
	return s
}

func (s *SelectTable) GetAs() (string, int) {
	return s.as, s.asNum
}

// C 返回序列化的列名。
func (s *SelectTable) C(column string) string {
	name := s.name
	b := &Builder{dialect: s.dialect}
	if s.as != "" {
		name = b.Quote(s.as)
	}
	if s.as == "" {
		b.WriteSchema(s.schema)
	}
	b.Ident(name).WriteByte('.').Ident(column)
	return b.String()
}

// Columns 返回一个序列号后的列名列表。
func (s *SelectTable) Columns(columns ...string) []string {
	names := make([]string, 0, len(columns))
	for _, c := range columns {
		names = append(names, s.C(c))
	}
	return names
}

// Unquote 使表名格式化为原始字符串（无引号）。
// 当不想查询当前数据库下的表时，它非常有用。
// 例如 在 MySQL 中为 "INFORMATION_SCHEMA.TABLE_CONSTRAINTS"。
func (s *SelectTable) Unquote() *SelectTable {
	s.quote = false
	return s
}

// ref 返回引用的表名。
func (s *SelectTable) ref() string {
	if !s.quote {
		return s.name
	}
	b := &Builder{dialect: s.dialect}
	b.WriteSchema(s.schema)
	b.Ident(s.name)
	if s.as != "" {
		b.WriteString(" AS ")
		b.Ident(s.as)
	}
	return b.String()
}

// view 是一个标记接口，用于实现TableView接口。
func (*SelectTable) view() {}

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
	s := Selector{
		Builder: &Builder{},
	}
	s.SetDialect(b.dialect)
	return &s
}

func (b *DialectBuilder) Table(name string) *SelectTable {
	t := Table(name)
	t.SetDialect(b.dialect)
	return t
}

type join struct {
	on    *Predicate
	kind  string
	table TableView
}

func (j join) clone() join {
	if sel, ok := j.table.(*Selector); ok {
		j.table = sel.Clone()
	}
	j.on = j.on.clone()
	return j
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
		*Builder
		ctx          context.Context
		as           string
		limit        *int
		selectFields []Selection
		from         []TableView
		where        *Predicate
		order        []*Order
		joins        []join
		table        *SelectTable
	}
	// Selection 选择的字段。
	Selection struct {
		field  FieldSpec
		entity string
	}
)

func (s Selection) String() string {
	return s.field.Name.String()
}

// Query 生成一个查询语句。
//
// Returns:
//
//	0: 查询语句。
//	1: 查询参数。
func (s *Selector) Query() (SqlSpec, error) {
	if len(s.from)+len(s.joins) > 1 {
		s.Builder.isAs = true
	}
	b := s.Builder.clone()
	b.WriteString("SELECT ")
	if len(s.selectFields) > 0 {
		s.appendSelect(b)
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
		switch t := from.(type) {
		case *SelectTable:
			t.SetDialect(s.dialect)
			b.WriteString(t.ref())
		case *Selector:
			t.SetDialect(s.dialect)
			b.Wrap(func(b *Builder) {
				b.Join(t)
			})
			if t.as != "" {
				b.WriteString(" AS ")
				b.Ident(t.as)
			}
		}
	}
	for _, join := range s.joins {
		b.WriteString(" " + join.kind + " ")
		switch view := join.table.(type) {
		case *SelectTable:
			view.SetDialect(s.dialect)
			b.WriteString(view.ref())
		case *Selector:
			view.SetDialect(s.dialect)
			b.Wrap(func(b *Builder) {
				b.Join(view)
			})
			b.WriteString(" AS ")
			b.Ident(view.as)
		}
		if join.on != nil {
			b.WriteString(" ON ")
			b.Join(join.on)
		}
	}
	if s.where != nil && len(s.where.fns) > 0 {
		b.WriteString(" WHERE ")
		b.Join(s.where)
	}
	batchSize := *(entity.GetConfig().BatchSize)
	if len(b.args) > batchSize {
		return SqlSpec{}, entity.Err_0100030004
	}
	if s.order != nil && len(s.order) > 0 {
		b.WriteString(" ORDER BY ")
		for i, order := range s.order {
			if i > 0 {
				b.Comma()
			}
			order.Query(b)
		}
	}
	if s.limit != nil {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(*s.limit))
	}
	return SqlSpec{Query: b.String(), Args: b.args}, nil
}

func (s *Selector) Clone() *Selector {
	if s == nil {
		return nil
	}
	joins := make([]join, len(s.joins))
	for i := range s.joins {
		joins[i] = s.joins[i].clone()
	}
	return &Selector{
		Builder: s.Builder.clone(),
		ctx:     s.ctx,
		as:      s.as,
		from:    s.from,
		limit:   s.limit,
		where:   s.where.clone(),
		joins:   append([]join{}, joins...),
		order:   append([]*Order{}, s.order...),
	}
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
func (s *Selector) Rows(rows ...FieldSpec) []FieldSpec {
	return rows
}

func (s *Selector) OnP(p *Predicate) *Selector {
	if len(s.joins) > 0 {
		join := &s.joins[len(s.joins)-1]
		switch {
		case join.on == nil:
			join.on = p
		default:
			join.on = AndPred(s.Builder, join.on, p)
		}
	}
	return s
}

// On sets the `ON` clause for the `JOIN` operation.
func (s *Selector) On(c1, c2 string) *Selector {
	s.OnP(P(s.Builder, func(builder *Builder) {
		builder.Ident(c1).WriteOp(OpEQ).Ident(c2)
	}))
	return s
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
//   - entity: 实体表的名称。
//   - rows: 字段的名称。
//
// Returns:
//
//	0: 选择语句生成器。
func (s *Selector) SetSelect(entity string, rows ...FieldSpec) *Selector {
	fields := make([]Selection, len(rows))
	for i := range rows {
		fields[i] = Selection{field: rows[i], entity: entity}
	}
	s.selectFields = append(s.selectFields, fields...)
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
func (s *Selector) SetFrom(t TableView) *Selector {
	s.from = nil
	return s.AppendFrom(t)
}

// AppendFrom appends a new TableView to the `FROM` clause.
func (s *Selector) AppendFrom(t TableView) *Selector {
	s.from = append(s.from, t)
	s.Builder.tables = append(s.Builder.tables, t)
	switch view := t.(type) {
	case *SelectTable:
		if view.as == "" {
			view.as, view.asNum = s.getAs()
		}
	case *Selector:
		if view.as == "" {
			view.as, _ = s.getAs()
		}
	}
	if st, ok := t.(state); ok {
		st.SetDialect(s.dialect)
	}
	if len(s.from) >= 1 {
		s.table = selectTable(s.from[0])
	}
	return s
}

func (s *Selector) SetOrder(order *Order) *Selector {
	s.order = append(s.order, order)
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

func (s *Selector) LeftJoin(t TableView) *Selector {
	return s.join("LEFT JOIN", t)
}

func (s *Selector) C(column string) string {
	// 跳过已经限定的列
	if s.isQualified(column) {
		return column
	}
	if s.as != "" {
		b := &Builder{dialect: s.dialect}
		b.Ident(s.as)
		b.WriteByte('.')
		b.Ident(column)
		return b.String()
	}
	return s.Table().C(column)
}

func (s *Selector) Table() *SelectTable {
	if len(s.from) == 0 {
		return nil
	}
	return selectTable(s.from[0])
}

// selectTable returns a *SelectTable from the given TableView.
func selectTable(t TableView) *SelectTable {
	if t == nil {
		return nil
	}
	switch view := t.(type) {
	case *SelectTable:
		return view
	case *Selector:
		if len(view.from) == 0 {
			return nil
		}
		return selectTable(view.from[0])
	default:
		panic(fmt.Sprintf("unexpected TableView %T", t))
	}
}

// appendSelect 添加选择的字段。
//
// Params:
//
//   - b: sql生成器。
func (s *Selector) appendSelect(b *Builder) {
	for i, col := range s.selectFields {
		if i > 0 {
			b.Comma()
		}
		if b.isAs {
			b.WriteString(b.Quote(col.entity))
			b.WriteString(".")
		}
		b.WriteString(col.field.NameFormat(s.dialect, b.Quote(col.field.Name.String())))
	}
}

// join 在selector中添加一个table
func (s *Selector) join(kind string, t TableView) *Selector {
	s.joins = append(s.joins, join{
		kind:  kind,
		table: t,
	})
	s.Builder.tables = append(s.Builder.tables, t)
	switch view := t.(type) {
	case *SelectTable:
		if view.as == "" {
			view.as, view.asNum = s.getAs()
		}
	case *Selector:
		if view.as == "" {
			view.as, _ = s.getAs()

		}
	}
	if st, ok := t.(state); ok {
		st.SetDialect(s.dialect)
	}
	return s
}

func (s *Selector) getAs() (string, int) {
	return ("t" + strconv.Itoa(len(s.joins)+len(s.from))), len(s.joins) + len(s.from)
}

func (*Selector) view() {}

/**************** Inserter 插入语句生成器 ***************/

type (
	// Inserter 插入语句生成器。
	Inserter struct {
		*Builder
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
		Builder: &Builder{},
		ctx:     ctx,
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
		i.writeDefault(b)
	} else {
		i.setInitialQuery(b)
		b.WriteByte('(').IdentComma(i.columns...).WriteByte(')')
		b.WriteString(" VALUES ")
		batchSize := *(entity.GetConfig().BatchSize)
		for j := 0; j < i.rowTotal; j++ {
			if current+len(i.columns) > batchSize {
				specs = append(specs, SqlSpec{Query: b.String(), Args: b.args})
				b = i.Builder.new()
				i.setInitialQuery(b)
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
	joinReturning(b, i.returning)
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
func (i *Inserter) Set(column string, as string, v any) *Inserter {
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
func (i *Inserter) SetColumns(columns ...FieldSpec) *Inserter {
	i.columns = make([]string, len(columns))
	for j := range columns {
		i.columns[j] = string(columns[j].Name)
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
		*Builder
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
		Builder: &Builder{},
		ctx:     ctx,
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
		u.writeSetter(b, i)
		if u.wheres[i] != nil {
			b.WriteString(" WHERE ")
			w := u.wheres[i]
			b.Join(w)
		}
		joinReturning(b, u.returning)
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
	t := Table(u.entity)
	t.SetDialect(b.dialect)
	for i, column := range u.columns[row] {
		if i > 0 {
			b.Comma()
		}
		b.Ident(column).WriteString(" = ")
		b.Join(NewCaser(u.Builder.clone(), column, u.values[row][i], t.as))
	}
}

/**************** Deleter 删除语句生成器 ***************/

type (
	// Deleter 删除语句生成器。
	Deleter struct {
		*Builder
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
		Builder: &Builder{},
		ctx:     ctx,
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
	d.setInitialQuery(b)
	if d.where != nil {
		length := d.where.FunsLen()
		batchSize := *(entity.GetConfig().BatchSize)
		batchNum := length/batchSize + 1
		for i := 0; i < batchNum; i++ {
			b.WriteString(" WHERE ")
			if length < (i+1)*batchSize {
				b.Join(d.where.Clone(i*batchSize, length))
			} else {
				b.Join(d.where.Clone(i*batchSize, (i+1)*batchSize))
			}
			specs = append(specs, SqlSpec{Query: b.String(), Args: b.args})
			b = d.Builder.new()
			d.setInitialQuery(b)
		}
	}
	return specs, nil
}

// setInitialQuery 设置初始的删除语句。
func (d *Deleter) setInitialQuery(b *Builder) {
	b.WriteString("DELETE FROM ")
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
		*Builder
		depth int
		// fns Where子句生成器的函数。
		fns []func(*Builder)
		// 上一个是否是逻辑运算符。
		lastIsLogic bool
	}
	// PredicateFunc Where子句生成器的函数。
	PredicateFunc func(p *Predicate)
)

func AndPred(builder *Builder, preds ...*Predicate) *Predicate {
	p := P(builder)
	return p.Append(func(b *Builder) {
		p.mayWrap(preds, b, "AND")
	})
}

// P 创建一个Where子句生成器。
func P(b *Builder, fns ...func(*Builder)) *Predicate {
	return &Predicate{Builder: b, fns: fns}
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
	var fns []func(*Builder)
	if begin < 0 {
		fns = p.fns[:end]
	} else if end < 0 {
		fns = p.fns[begin:]
	} else {
		fns = p.fns[begin:end]
	}
	return &Predicate{fns: fns, Builder: p.Builder.clone(), lastIsLogic: p.lastIsLogic}
}

func (p *Predicate) clone() *Predicate {
	if p == nil {
		return p
	}
	return &Predicate{fns: append([]func(*Builder){}, p.fns...), Builder: p.Builder.clone(), lastIsLogic: p.lastIsLogic}
}

// Query 生成查询中的Where子句。
//
// Returns:
//
//	0: 查询中的Where子句。
//	1: 查询中的参数。
func (p *Predicate) Query() (SqlSpec, error) {
	if p.Len() > 0 || len(p.args) > 0 {
		p.Reset()
		p.args = nil
	}
	for i, f := range p.fns {
		// 添加这个判断是因为在批量操作中，如果参数过多会执行分批操作，此时可能出现第一个或最后一个是AND或者OR的情况，
		if i == 0 || i == len(p.fns)-1 {
			nb := p.Builder.new()
			f(nb)
			s := nb.String()
			if strings.Contains(s, "AND") || strings.Contains(s, "OR") {
				continue
			}

		}
		f(p.Builder)
	}
	return SqlSpec{
		Query: p.String(),
		Args:  p.args,
	}, nil
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

func (p *Predicate) mayWrap(preds []*Predicate, b *Builder, op string) {
	switch n := len(preds); {
	case n == 1:
		b.Join(preds[0])
		return
	case n > 1 && p.depth != 0:
		b.WriteByte('(')
		defer b.WriteByte(')')
	}
	for i := range preds {
		preds[i].depth = p.depth + 1
		if i > 0 {
			b.WriteByte(' ')
			b.WriteString(op)
			b.WriteByte(' ')
		}
		if len(preds[i].fns) > 1 {
			b.Wrap(func(b *Builder) {
				b.Join(preds[i])
			})
		} else {
			b.Join(preds[i])
		}
	}
}

// And 添加一个AND标识符。
//
// Returns:
//
//	0: Where子句生成器。
func (p *Predicate) And() *Predicate {
	p.lastIsLogic = true
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
	p.lastIsLogic = true
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
	p.lastIsLogic = false
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
func (p *Predicate) EQ(column string, as string, v any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) NEQ(column string, as string, v any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) GT(column string, as string, v any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) GTE(column string, as string, v any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) LT(column string, as string, v any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) LTE(column string, as string, v any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) In(column string, as string, v ...any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) NotIn(column string, as string, v ...any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) Like(column string, as string, v any) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) IsNull(column string, as string) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
func (p *Predicate) NotNull(column string, as string) *Predicate {
	p.lastIsLogic = false
	return p.Append(func(b *Builder) {
		if b.isAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
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
	p.lastIsLogic = false
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
	p.lastIsLogic = false
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
	p.lastIsLogic = false
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
	p.lastIsLogic = false
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
	p.lastIsLogic = false
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

func (p *Predicate) isLogic() bool {
	return p.lastIsLogic
}

func (pf PredicateFunc) isOp(p *Predicate) bool {
	np := p.clone()
	pf(np)
	return np.isLogic()
}

/**************** CASE 批量更新中Set语句生成器 ***************/

type (
	// Caser CASE语句生成器。
	Caser struct {
		*Builder
		// Cases CASE语句信息。
		Cases  []CaseSpec
		column string
		elser  any
		as     string
	}
)

// NewCaser 创建一个CASE语句生成器。
//
// Params:
//
//   - column: 列名。
//   - cases: CASE语句信息。
//   - as: 实体别名。
//
// Returns:
//
//	0: CASE语句生成器。
func NewCaser(buider *Builder, column string, cases []CaseSpec, as string) *Caser {
	return &Caser{
		Builder: buider,
		Cases:   cases,
		column:  column,
		as:      as,
	}
}

// Query 生成CASE语句。
//
// Returns:
//
//	0: CASE语句。
//	1: CASE语句的参数。
func (c *Caser) Query() (SqlSpec, error) {
	if c.Len() > 0 || len(c.args) > 0 {
		c.Reset()
		c.args = nil
	}
	if len(c.Cases) == 0 {
		return SqlSpec{}, nil

	} else if len(c.Cases) == 1 && c.Cases[0].When == nil {
		c.arg(c.Cases[0].Field)
		return SqlSpec{
			Query: c.String(),
			Args:  c.args,
		}, nil
	} else {
		c.WriteString("CASE ")
		for _, cs := range c.Cases {
			// 对于多个没有条件的CASE，只有最后一个会生效，
			// 如果要做限制或者判断,请在外部实现。
			if cs.When == nil {
				c.elser = cs.Field
			} else {
				c.When(cs.When)
				c.Then(cs.Field)
			}
		}
		c.Else()
		c.WriteString(" END")
	}
	return SqlSpec{
		Query: c.String(),
		Args:  c.args,
	}, nil

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
func (c *Caser) When(pred PredicateFunc) *Caser {
	c.WriteString(" WHEN ")
	p := P(c.Builder.clone())
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

type (
	Order struct {
		OrderOptions
		Builder
		Column string
		// Orders 排序信息。
	}

	OrderOptions struct {
		Desc       bool
		As         string // Optional alias.
		NullsFirst bool   // Whether to sort nulls first.
		NullsLast  bool   // Whether to sort nulls last.
	}

	OrderFunc func(*Order)
)

func O() *Order {
	return &Order{}
}

func (o *Order) Query(b *Builder) {
	if b.isAs && o.As != "" {
		b.WriteString(b.Quote(o.As))
		b.WriteString(".")
	}
	b.Ident(o.Column)
	if o.OrderOptions.Desc {
		b.WriteString(" DESC")
	}
	if o.OrderOptions.NullsFirst {
		b.WriteString(" NULLS FIRST")
	}
	if o.OrderOptions.NullsLast {
		b.WriteString(" NULLS LAST")
	}
}

func (o *Order) SetDialect(dialect dialect.DbDriver) *Order {
	o.dialect = dialect
	return o
}

func (o *Order) SetAs(schema string) *Order {
	o.OrderOptions.As = schema
	return o
}

func (o *Order) SetColumn(column string) *Order {
	o.Column = column
	return o
}

func (o *Order) SetOp(op string) *Order {
	o.Column = op
	return o
}

func (o *Order) Desc() *Order {
	o.OrderOptions.Desc = true
	return o
}

func (o *Order) Asc() *Order {
	o.OrderOptions.Desc = false
	return o
}

func (o *Order) NullsFirst() *Order {
	o.OrderOptions.NullsFirst = true
	return o
}

func (o *Order) NullsLast() *Order {
	o.OrderOptions.NullsLast = true
	return o
}
