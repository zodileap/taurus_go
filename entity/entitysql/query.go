package entitysql

import (
	"context"
	"fmt"

	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/dialect"
	"github.com/zodileap/taurus_go/tlog"
)

type (
	// QueryContext 用于Query，包含了查询语句上下的信息，
	// 比如Limit、Offset、Order等，用于生成查询语句。
	QueryContext struct {
		// Limit 限制查询语句返回的记录数。
		// 调用Limit方法时，会将Limit设置为指定的值。
		// 比如：Limit(10), sql: Select * from user limit 10。
		Limit  *int
		Fields []FieldName
	}
	// QueryContextKey 用于在context中存储QueryContext。
	QueryContextKey struct{}

	RelationDesc struct {
		Orders     []OrderFunc
		Predicates []PredicateFunc
		To         RelationTable
		Join       RelationTable
	}

	RelationTable struct {
		Table   string
		Field   string
		Columns []FieldName
	}

	// Relation 用于生成联表查询。
	Relation func(*Selector)
)

func (r *RelationDesc) Reset() {
	r.Orders = []OrderFunc{}
	r.Predicates = []PredicateFunc{}
}

func AddRelBySelector(s *Selector, t *SelectTable, desc RelationDesc) *SelectTable {
	var (
		build = NewDialect(s.Dialect())
	)
	joinT := build.Table(desc.Join.Table).Schema(s.Table().schema)
	s.LeftJoin(joinT).On(t.C(desc.To.Field), joinT.C(desc.Join.Field))
	s.SetSelect(joinT.as, s.Rows(NewFieldSpecs(desc.Join.Columns...)...)...)

	if orders := desc.Orders; len(orders) > 0 {
		for _, order := range orders {
			o := O()
			o.SetDialect(joinT.dialect)
			o.SetAs(joinT.as)
			order(o)
			s.SetOrder(o)
		}
	}

	if ps := desc.Predicates; len(ps) > 0 {
		if s.where.FunsLen() > 0 && !ps[0].isOp(s.where) {
			s.where.And()
		}
		for _, p := range ps {
			p(s.where)
		}
	}

	return joinT
}

// NewQueryContext 将QueryContext添加到context中，并返回一个新的context。
//
// Params:
//
//   - parent: 父context。
//   - c: 用于Query，包含了查询语句上下的信息。
//
// Returns:
//
//	0: 新的context。
func NewQueryContext(parent context.Context, c *QueryContext) context.Context {
	return context.WithValue(parent, QueryContextKey{}, c)
}

// QueryFromContext 从context中获取QueryContext。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: QueryContext。
func QueryFromContext(ctx context.Context) *QueryContext {
	c, _ := ctx.Value(QueryContextKey{}).(*QueryContext)
	return c
}

// QuerySpec 包含了实体的查询的信息。
// 通过把Query中的信息转换为QuerySpec，把QuerySpec传给NewQuery发起查询,
// QuerySpec是把查询信息转为sql的中间件。
type QuerySpec struct {
	// Entity 表示实体的信息。
	Entity *EntitySpec
	// Scan 用于扫描返回的数据。
	Scan Scanner
	// Limit 限制查询语句返回的记录数。
	Limit int
	// Predicate 查询语句的条件，用于生成where子句。
	Predicate PredicateFunc
	// Rels 用于生成联表查询。
	Rels   []Relation
	Orders []OrderFunc
}

// NewQuerySpec 创建一个QuerySpec。
//
// Params:
//   - entity：实体的名称。
//   - rows：实体的字段。
//
// Returns:
//
//	0: QuerySpec。
func NewQuerySpec(entity string, rows []FieldName) *QuerySpec {
	columns := NewFieldSpecs(rows...)
	return &QuerySpec{
		Entity: &EntitySpec{
			Name:    entity,
			Columns: columns,
		},
	}
}

// NewQuery 查询一个实体，并将返回的结果扫描到指定的值中。
func NewQuery(ctx context.Context, drv dialect.Driver, spec *QuerySpec) error {
	builder := NewDialect(drv.Dialect())
	qb := queryBuilder{QuerySpec: spec, entityBuilder: entityBuilder{builder: builder}}
	return qb.query(ctx, drv)
}

// queryBuilder 查询语句生成器。
type queryBuilder struct {
	entityBuilder
	*QuerySpec
}

// query 查询一个实体，并将返回的结果扫描到指定的值中。
//
// Params:
//
//   - ctx: 上下文。
//   - drv: 数据库连接。
func (b *queryBuilder) query(ctx context.Context, drv dialect.Driver) error {
	selector, err := b.selector(ctx)
	if err != nil {
		return err
	}
	spec, err := selector.Query()
	if err != nil {
		return err
	}
	config := entity.GetConfig()
	if *(config.SqlConsole) {
		tlog.Debug(*config.SqlLogger, fmt.Sprintf("sql: %s", spec.Query))
		tlog.Debug(*config.SqlLogger, fmt.Sprintf("args: %v", spec.Args))
	}
	var rows dialect.Rows
	err = drv.Query(ctx, spec.Query, spec.Args, &rows)
	if err != nil {
		return err
	}
	for rows.Next() {
		ScannerFields := make([]ScannerField, len(b.Entity.Columns))
		for i, c := range b.Entity.Columns {
			ScannerFields[i] = c.Name
		}
		err := b.Scan(rows, ScannerFields)
		if err != nil {
			return err
		}
	}

	return nil
}

// selector 生成查询语句。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 选择语句生成器。
func (b *queryBuilder) selector(ctx context.Context) (*Selector, error) {
	selector := b.builder.Select()
	t := b.builder.Table(b.Entity.Name)
	selector.SetFrom(t)
	selector.SetSelect(t.as, selector.Rows(b.Entity.Columns...)...)
	selector.SetContext(ctx)
	if b.Limit != 0 {
		selector.SetLimit(b.Limit)
	}

	selector.where = P(selector.Builder)
	if pred := b.Predicate; pred != nil {
		pred(selector.where)
	}
	if orders := b.Orders; len(orders) > 0 {
		for _, order := range orders {
			o := O()
			o.SetDialect(b.builder.dialect)
			o.SetAs(t.as)
			order(o)
			selector.SetOrder(o)
		}
	}
	if rels := b.Rels; len(rels) > 0 {
		for _, rel := range rels {
			rel(selector)
		}
	}

	return selector, nil
}
