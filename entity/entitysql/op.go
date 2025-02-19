package entitysql

type Op int

// Predicate and arithmetic operators.
const (
	// OpEQ 等于操作符。
	OpEQ Op = 0 // =
	// OpNEQ 不等于操作符。
	OpNEQ Op = 1 // <>
	// OpGT 大于操作符。
	OpGT Op = 2 // >
	// OpGTE 大于等于操作符。
	OpGTE Op = 3 // >=
	// OpLT 小于操作符。
	OpLT Op = 4 // <
	// OpLTE 小于等于操作符。
	OpLTE Op = 5 // <=
	// OpIn IN操作符。
	OpIn Op = 6 // IN
	// OpNotIn NOT IN操作符。
	OpNotIn Op = 7 // NOT IN
	// OpLike LIKE操作符。
	OpLike Op = 8 // LIKE
	// OpIsNull IS NULL操作符。
	OpIsNull Op = 9 // IS NULL
	// OpNotNull IS NOT NULL操作符。
	OpNotNull Op = 10 // IS NOT NULL
	// OpAdd 加操作符。
	OpAdd Op = 11 // +
	// OpSub 减操作符。
	OpSub Op = 12 // -
	// OpMul 乘操作符。
	OpMul Op = 13 // *
	// OpDiv 除操作符。
	OpDiv Op = 14 // / (Quotient)
	// OpMod 取模操作符。
	OpMod Op = 15 // % (Reminder)
	// OpContains 包含操作符。
	OpContains Op = 16 // @> (Contains)
)

var ops [17]string = [17]string{
	OpEQ:       "=",
	OpNEQ:      "<>",
	OpGT:       ">",
	OpGTE:      ">=",
	OpLT:       "<",
	OpLTE:      "<=",
	OpIn:       "IN",
	OpNotIn:    "NOT IN",
	OpLike:     "LIKE",
	OpIsNull:   "IS NULL",
	OpNotNull:  "IS NOT NULL",
	OpAdd:      "+",
	OpSub:      "-",
	OpMul:      "*",
	OpDiv:      "/",
	OpMod:      "%",
	OpContains: "@>",
}

func (op Op) String() string {
	return ops[op]
}

var And PredicateFunc = func(p *Predicate) {
	p.And()
}

var Or PredicateFunc = func(p *Predicate) {
	p.Or()
}

var Not PredicateFunc = func(p *Predicate) {
	p.Not()
}

var Add PredicateFunc = func(p *Predicate) {
	p.Add()
}

var Sub PredicateFunc = func(p *Predicate) {
	p.Sub()
}

var Mul PredicateFunc = func(p *Predicate) {
	p.Mul()
}

var Div PredicateFunc = func(p *Predicate) {
	p.Div()
}

var Mod PredicateFunc = func(p *Predicate) {
	p.Mod()
}
