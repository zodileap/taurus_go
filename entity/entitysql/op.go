package entitysql

type Op int

// Predicate and arithmetic operators.
const (
	OpEQ      Op = iota // =
	OpNEQ               // <>
	OpGT                // >
	OpGTE               // >=
	OpLT                // <
	OpLTE               // <=
	OpIn                // IN
	OpNotIn             // NOT IN
	OpLike              // LIKE
	OpIsNull            // IS NULL
	OpNotNull           // IS NOT NULL
	OpAdd               // +
	OpSub               // -
	OpMul               // *
	OpDiv               // / (Quotient)
	OpMod               // % (Reminder)
)

var ops = [...]string{
	OpEQ:      "=",
	OpNEQ:     "<>",
	OpGT:      ">",
	OpGTE:     ">=",
	OpLT:      "<",
	OpLTE:     "<=",
	OpIn:      "IN",
	OpNotIn:   "NOT IN",
	OpLike:    "LIKE",
	OpIsNull:  "IS NULL",
	OpNotNull: "IS NOT NULL",
	OpAdd:     "+",
	OpSub:     "-",
	OpMul:     "*",
	OpDiv:     "/",
	OpMod:     "%",
}

func (op Op) String() string {
	return ops[op]
}

var OpAnd = func(p *Predicate) {
	p.And()
}

var OpOr = func(p *Predicate) {
	p.Or()
}

var OpNot = func(p *Predicate) {
	p.Not()
}
