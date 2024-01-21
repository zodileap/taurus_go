package entitysql

type Field struct {
	Name     FieldName
	Primary  int
	Default  bool
	Required bool
}
