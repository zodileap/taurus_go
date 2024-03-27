package entitysql

// Field 字段信息。
type Field struct {
	// Name 字段名称。
	Name FieldName
	// Primary 主键值。0表示不是主键，大于0表示主键。
	Primary int
	// Default 是否默认值。
	Default bool
	// Required 是否必填。
	Required bool
}
