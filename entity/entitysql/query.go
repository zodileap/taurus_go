package entitysql

import (
	"context"

	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/tlog"
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
	QueryContextKey struct{}
)

// NewQueryContext 将QueryContext添加到context中，并返回一个新的context。
func NewQueryContext(parent context.Context, c *QueryContext) context.Context {
	return context.WithValue(parent, QueryContextKey{}, c)
}

// QueryFromContext 从context中获取QueryContext。
func QueryFromContext(ctx context.Context) *QueryContext {
	c, _ := ctx.Value(QueryContextKey{}).(*QueryContext)
	return c
}

// QuerySpec 包含了实体的查询的信息。
// 通过把Query中的信息转换为QuerySpec，把QuerySpec传给NewQuery发起查询,
// QuerySpec是把查询信息转为sql的中间件。
type QuerySpec struct {
	Entity    *EntitySpec
	Scan      Scanner
	Limit     int
	Predicate func(*Predicate)
}

// NewQuerySpec 创建一个QuerySpec。
//
// 参数：
//   - entity：实体的名称。
func NewQuerySpec(entity string, rows []FieldName) *QuerySpec {
	return &QuerySpec{
		Entity: &EntitySpec{
			Name: entity,
			Rows: rows,
		},
	}
}

// NewQuery 查询一个实体，并将返回的结果扫描到指定的值中。
func NewQuery(ctx context.Context, drv dialect.Driver, spec *QuerySpec) error {
	builder := NewDialect(drv.Dialect())
	qb := queryBuilder{QuerySpec: spec, entityBuilder: entityBuilder{builder: builder}}
	return qb.query(ctx, drv)
}

type queryBuilder struct {
	entityBuilder
	*QuerySpec
}

func (b *queryBuilder) query(ctx context.Context, drv dialect.Driver) error {
	selector, err := b.selector(ctx)
	if err != nil {
		return err
	}
	query, args := selector.Query()
	tlog.Print(query)
	var rows dialect.Rows
	err = drv.Query(context.Background(), query, args, &rows)
	if err != nil {
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

func (b *queryBuilder) selector(ctx context.Context) (*Selector, error) {
	selector := b.builder.Select()
	selector.SetFrom(b.Entity.Name)
	selector.SetSelect(selector.Rows(b.Entity.Rows...)...)
	selector.SetContext(ctx)
	if b.Limit != 0 {
		selector.SetLimit(b.Limit)
	}
	if pred := b.Predicate; pred != nil {
		selector.where = P()
		pred(selector.where)
	}

	return selector, nil
}
