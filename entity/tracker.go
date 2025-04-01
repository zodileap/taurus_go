package entity

import (
	"context"

	stringutil "github.com/zodileap/taurus_go/datautil/string"
	"github.com/zodileap/taurus_go/entity/dialect"
)

type (
	// Tracker 是一个接口，用于跟踪实体类的状态。
	Tracker interface {
		Add(...Mutator)
		Mutators() []Mutator
		Clear()
	}
	// Tracking 用于跟踪实体类的状态。
	Tracking struct {
		// mutators 用于存储目前正在追踪的实体。
		mutators []Mutator
	}
)

// Add 用于添加一个实体到追踪器中。
//
// Params:
//
//   - m: 需要追踪的实体。
func (t *Tracking) Add(m ...Mutator) {
	t.mutators = append(t.mutators, m...)
}

// Mutators 用于获取追踪器中的实体。
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
		// Exec 根据实体的状态执行操作。
		Exec(ctx context.Context, tx dialect.Tx) error
	}

	// Mutation 实体的修改器，用于存储实体的状态以及需要改变的字段。
	Mutation struct {
		key    string
		state  EntityState
		fields []string
	}
)

// NewMutation 创建一个新的实体的修改器。
//
// Params:
//
//   - state: 实体当前的状态。
//
// Returns:
//
//	0: 新的实体修改器。
func NewMutation(state EntityState) *Mutation {
	key, err := stringutil.GenerateKey()
	if err != nil {
		panic(err)
	}
	return &Mutation{state: state, key: key}
}

// Key 获取修改器的实体的键。
//
// Returns:
//
//	0: 实体的键。
func (m Mutation) Key() string {
	return m.key
}

// SetState 设置实体的状态。
//
// Params:
//
//   - state: 实体的状态。
func (m *Mutation) SetState(state EntityState) {
	m.state = state
}

// State 获取实体的状态。
//
// Returns:
//
//	0: 实体的状态。
func (m Mutation) State() EntityState {
	return m.state
}

// SetFields 设置需要改变的字段。
//
// Params:
//
//   - fields: 需要改变的字段。
func (m *Mutation) SetFields(fields ...string) {
	m.fields = append(m.fields, fields...)
}

// Fields 获取改变值的字段。主要是用于更新操作。
//
// Returns:
//
//	0: 需要改变的字段。
func (m Mutation) Fields() []string {
	return m.fields
}

// EntityState 实体类状态，用于标识实体类的状态。
type EntityState = int16

const (
	// NotSet 未设置，这个不作为实体类的状态，是用于设置实体类的状态的操作
	NotSet EntityState = -1
	// Detached 未追踪,不存在于数据库中、属性未修改、调用Save()方法时不会执行执行操作。
	Detached EntityState = 0
	// Unchanged 未修改，存在于数据库中、属性未修改、调用Save()方法时，不会执行执行操作。
	Unchanged EntityState = 1
	// Deleted 已删除，存在于数据库中、调用Save()方法时，会执行删除操作。
	Deleted EntityState = 2
	// Modified 已修改，存在于数据库中、属性已修改、调用Save()方法时，会执行更新操作。
	Modified EntityState = 3
	// Added 已添加，不存在于数据库中、属性已修改、调用Save()方法时，会执行插入操作。
	Added EntityState = 4
)
