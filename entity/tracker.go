package entity

type (
	Tracker interface {
		Add(...Mutator)
		Mutators() []Mutator
		Clear()
	}
	// Tracking 用于跟踪实体类的状态。
	Tracking struct {
		// Mutators 用于存储目前正在追踪的实体。
		mutators []Mutator
	}
)

// Add 用于添加一个实体到追踪器中。
func (t *Tracking) Add(m ...Mutator) {
	t.mutators = append(t.mutators, m...)
}

func (t *Tracking) Mutators() []Mutator {
	return t.mutators
}

// Clear 用于清空追踪器中的实体。
func (t *Tracking) Clear() {
	t.mutators = make([]Mutator, 0)
}

type (
	// Mutator 是一个接口，用于标记实体的状态以及根据状态执行操作。
	Mutator interface {
		// State 获取实体的状态。
		State() EntityState
		// SetState 设置实体的状态。
		SetState(state EntityState)
		// Fields 获取改变值的字段。主要是用于更新操作。
		Fields() []string
		SetFields(fields ...string)
		// Exec 根据实体的状态执行操作。
		Exec() error
	}

	Mutation struct {
		state  EntityState
		fields []string
	}
)

func NewMutation(state EntityState) *Mutation {
	return &Mutation{state: state}
}

func (m *Mutation) SetState(state EntityState) {
	m.state = state
}

func (m *Mutation) State() EntityState {
	return m.state
}

func (m *Mutation) SetFields(fields ...string) {
	m.fields = append(m.fields, fields...)
}

func (m *Mutation) Fields() []string {
	return m.fields
}

// 实体类状态，用于标识实体类的状态。
type EntityState = int16

const (
	// 未追踪,不存在于数据库中、属性未修改、调用Save()方法时不会执行执行操作。
	Detached EntityState = 0
	// 未修改，存在于数据库中、属性未修改、调用Save()方法时，不会执行执行操作。
	Unchanged EntityState = 1
	// 已删除，存在于数据库中、调用Save()方法时，会执行删除操作。
	Deleted EntityState = 2
	// 已修改，存在于数据库中、属性已修改、调用Save()方法时，会执行更新操作。
	Modified EntityState = 3
	// 已添加，不存在于数据库中、属性已修改、调用Save()方法时，会执行插入操作。
	Added EntityState = 4
)
