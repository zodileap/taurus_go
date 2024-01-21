package entitysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yohobala/taurus_go/entity/dialect"
)

type Builder struct {
	// 用于构建查询的字符串构建器。
	sb *strings.Builder
	// 使用的数据库驱动
	dialect dialect.DbDriver // configured dialect.
	// 查询的参数
	args []any
	// 查询树中总共出现的参数数量，在复杂的查询中，可能会出现多个相同的参数。
	// 所以数量可能会大于len(args)。
	total int
	// 限定符作为标识符（如表名）的前缀。
	qualifier string
}

// Querier 封装了entity中不同构建器所实现的基本查询方法.
// 假设有两个实现了 Querier 接口的对象，它们的 Query 方法分别返回以下内容：

// 第一个对象返回：("SELECT * FROM users WHERE age > ?", [30])
// 第二个对象返回：("AND name = ?", ["John"])
// 当使用 Join 方法将这两个对象组合时，生成的SQL查询将是：
// SELECT * FROM users WHERE age > ? AND name = ?
type Querier interface {
	// Query 返回元素的查询表示形式以及与之相关的参数
	Query() (string, []any)
}

// raw 插入不需要要转义的原始字符串。
type raw struct{ s string }

func Raw(s string) Querier { return &raw{s} }

// Query 返回原始字符串，不需要参数。
func (r *raw) Query() (string, []any) { return r.s, nil }

// state 封装了所有用于设置和获取更新状态的所有方法。
type state interface {
	Dialect() dialect.DbDriver
	SetDialect(dialect.DbDriver)
	Total() int
	SetTotal(int)
}

// String 把生成器中的查询语句转换为字符串。
func (b *Builder) String() string {
	if b.sb == nil {
		return ""
	}
	return b.sb.String()
}

func (b *Builder) Len() int {
	if b.sb == nil {
		return 0
	}
	return b.sb.Len()
}

func (b *Builder) Reset() {
	if b.sb != nil {
		b.sb.Reset()
	}
}

// SetDialect 设置生成器使用的数据库驱动。满足state接口。
func (b *Builder) SetDialect(dialect dialect.DbDriver) {
	b.dialect = dialect
}

// Dialect 返回生成器使用的数据库驱动。满足state接口。
func (b Builder) Dialect() dialect.DbDriver {
	return b.dialect
}

// Total 返回查询树中总共出现的参数数量。满足state接口。
func (b Builder) Total() int {
	return b.total
}

// SetTotal 设置查询树中总共出现的参数数量。满足state接口。
// 用于在查询、表达式中传递参数数量。
func (b *Builder) SetTotal(total int) {
	b.total = total
}

// Quote 根据配置的dialect，选择不同的字符引用SQL标识符。默认为"`"（通常用于MySQL)。
// 用于区分关键字，特殊字符等。
func (b *Builder) Quote(ident string) string {
	quote := "`"
	switch {
	// 如果是PostgreSQL，使用双引号。
	case b.postgres():
		if strings.Contains(ident, "`") {
			return strings.ReplaceAll(ident, "`", `"`)
		}
		quote = `"`
	// 未知的dialect，使用原始的标识符。
	case string(b.dialect) == "" && strings.ContainsAny(ident, "`\""):
		return ident
	}
	return quote + ident + quote
}

// Ident 添加标识符到查询中。标识符可以是表名、列名、别名等。
func (b *Builder) Ident(s string) *Builder {
	switch {
	// 忽略空字符串。
	case len(s) == 0:
	// 添加限定符和标识符
	case !strings.HasSuffix(s, "*") && !b.isIdent(s) && !isFunc(s) && !isModifier(s) && !isAlias(s):
		if b.qualifier != "" {
			b.WriteString(b.Quote(b.qualifier)).WriteByte('.')
		}
		b.WriteString(b.Quote(s))
	// 函数、修饰符、别名的特殊处理（针对PostgreSQL）
	case (isFunc(s) || isModifier(s) || isAlias(s)) && b.postgres():
		b.WriteString(strings.ReplaceAll(s, "`", `"`))
	default:
		b.WriteString(s)
	}
	return b
}

// IdentComma 添加标识符到查询中，用逗号分隔
func (b *Builder) IdentComma(s ...string) *Builder {
	for i := range s {
		if i > 0 {
			b.Comma()
		}
		b.Ident(s[i])
	}
	return b
}

type (
	// StmtInfo 保存SQL语句或数据库连接的上下文信息
	StmtInfo struct {
		// Dialect 数据库驱动
		Dialect dialect.DbDriver
	}
	// ParamFormatter 定义了FormatParam方法，用于格式化占位符。
	ParamFormatter interface {
		// 这个接口可以被用于特定的场景，
		// 例如当你使用的数据库需要一种特殊的参数格式时。
		// 例子：如果你在MySQL中使用地理空间数据，
		// 你可能需要将标准的参数占位符（如?）转换为特定的函数调用（如ST_GeomFromWKB(?)）。
		// 参数
		//   - placeholder: 标准的占位符，如?
		//   - info: 保存SQL语句或数据库连接的上下文信息
		FormatParam(placeholder string, info *StmtInfo) string
	}
)

// Arg 添加一个参数到生成器中。
func (b *Builder) Arg(a any) *Builder {
	switch v := a.(type) {
	case nil:
		b.WriteString("NULL")
		return b
	// 如果是原始字符串，直接把字符串添加到查询中。
	case *raw:
		b.WriteString(v.s)
		return b
	// 如果是查询构建器，把它们拼接到查询中。
	case Querier:
		b.Join(v)
		return b
	}
	// 默认的占位符参数（MySQL和SQLite）。
	format := "?"
	if b.postgres() {
		// Postgres参数使用语法$n引用
		format = "$" + strconv.Itoa(b.total+1)
	}
	// 如果参数实现了ParamFormatter接口，使用它的FormatParam方法。
	// 比如 postgesSQL中的PostGIS的ST_GeomFromGeoJSON($1)
	if f, ok := a.(ParamFormatter); ok {
		format = f.FormatParam(format, &StmtInfo{
			Dialect: b.dialect,
		})
	}
	return b.Argf(format, a)
}

// Argf 添加多个参数到生成器中。
func (b *Builder) Args(args ...any) *Builder {
	for i := range args {
		if i > 0 {
			b.Comma()
		}
		b.Arg(args[i])
	}
	return b
}

// Argf 将输入参数以给定的格式添加到生成器中。
//
//	Argf("JSON(?)", b).
//	Argf("ST_GeomFromText(?)", geom)
func (b *Builder) Argf(format string, a any) *Builder {
	switch a := a.(type) {
	case nil:
		b.WriteString("NULL")
		return b
	case *raw:
		b.WriteString(a.s)
		return b
	case Querier:
		b.Join(a)
		return b
	}
	b.total++
	b.args = append(b.args, a)
	b.WriteString(format)
	return b
}

// Wrap 获取一个回调函数，把它包装在括号中，并添加到查询中。
func (b *Builder) Wrap(f func(*Builder)) *Builder {
	nb := &Builder{dialect: b.dialect, total: b.total, sb: &strings.Builder{}}
	nb.WriteByte('(')
	f(nb)
	nb.WriteByte(')')
	b.WriteString(nb.String())
	b.args = append(b.args, nb.args...)
	b.total = nb.total
	return b
}

// WriteByte 添加一个字节到查询中。
func (b *Builder) WriteByte(c byte) *Builder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}
	b.sb.WriteByte(c)
	return b
}

// WriteString 添加一个字符串到查询中。
func (b *Builder) WriteString(s string) *Builder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}
	b.sb.WriteString(s)
	return b
}

// WriteSchema 添加一个模式到查询中。
func (b *Builder) WriteSchema(schema string) *Builder {
	if schema != "" {
		b.Ident(schema).WriteByte('.')
	}
	return b
}

func (b *Builder) WriteOp(op Op) *Builder {
	switch {
	case op >= OpEQ && op <= OpLike || op >= OpAdd && op <= OpMod:
		b.Blank().WriteString(op.String()).Blank()
	case op == OpIsNull || op == OpNotNull:
		b.Blank().WriteString(op.String())
	default:
		panic(fmt.Sprintf("invalid op %d", op))
	}
	return b
}

// Comma 添加一个逗号到查询中。
func (b *Builder) Comma() *Builder {
	b.WriteString(", ")
	return b
}

// Blank 添加一个空格到查询中。
func (b *Builder) Blank() *Builder {
	b.WriteString(" ")
	return b
}

// Join 添加多个查询到生成器中。
func (b *Builder) Join(qs ...Querier) *Builder {
	return b.join(qs, "")
}

// join 添加多个查询到生成器中，用分隔符分隔。
func (b *Builder) join(qs []Querier, sep string) *Builder {
	for i, q := range qs {
		if i > 0 {
			b.WriteString(sep)
		}
		st, ok := q.(state)
		if ok {
			st.SetDialect(b.dialect)
			st.SetTotal(b.total)
		}
		query, args := q.Query()
		b.WriteString(query)
		b.args = append(b.args, args...)
		b.total += len(args)
	}
	return b
}

// clone 克隆查询构建器。
func (b Builder) clone() Builder {
	c := Builder{dialect: b.dialect, total: b.total, sb: &strings.Builder{}}
	if len(b.args) > 0 {
		c.args = append(c.args, b.args...)
	}
	if b.sb != nil {
		c.sb.WriteString(b.sb.String())
	}
	return c
}

// postgres 检查是否是PostgreSQL。
func (b *Builder) postgres() bool {
	return b.dialect == dialect.PostgreSQL
}

// mysql 检查是否是MySQL。
func (b *Builder) mysql() bool {
	return b.dialect == dialect.MySQL
}

// isIdent 检查字符串是否包含标识符。标识符：["]、[`]
func (b *Builder) isIdent(s string) bool {
	switch {
	case b.postgres():
		return strings.Contains(s, `"`)
	default:
		return strings.Contains(s, "`")
	}
}

// joinReturning 添加RETURNING子句到查询中，MySQL不支持。
func joinReturning(b *Builder, columns []FieldName) {
	if len(columns) == 0 || (!b.postgres()) {
		return
	}

	s := []string{}
	for _, c := range columns {
		s = append(s, c.String())
	}

	b.WriteString(" RETURNING ")
	b.IdentComma(s...)
}

// isAlias 检查字符串是否包含别名。别名：[ AS ]、[ as ]
func isAlias(s string) bool {
	return strings.Contains(s, " AS ") || strings.Contains(s, " as ")
}

// isFunc 检查字符串是否包含函数。函数：[(]、[)]
func isFunc(s string) bool {
	return strings.Contains(s, "(") && strings.Contains(s, ")")
}

// isModifier 检查字符串是否包含修饰符。修饰符：[DISTINCT]、[ALL]、[WITH ROLLUP]
func isModifier(s string) bool {
	for _, m := range [...]string{"DISTINCT", "ALL", "WITH ROLLUP"} {
		if strings.HasPrefix(s, m) {
			return true
		}
	}
	return false
}
