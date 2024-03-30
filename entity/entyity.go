package entity

import (
	"database/sql/driver"

	"github.com/yohobala/taurus_go/entity/dialect"
)

// 单条语句，参数的最大数量，以PostgreSQL为标准。
const BatchSize int = 65535

// 和数据库连接相关定义。
type ConnectionConfig struct {
	// 数据库驱动
	Driver dialect.DbDriver

	// 用于标记当前的连接，Db通过这个tag绑定连接。
	Tag string
	// 数据库链接。
	Host string
	// 数据库端口。
	Port int
	// 数据库用户。
	User string
	// 数据库密码。
	Password string
	// 数据库名称。
	DBName string

	// 是否开启SSL的verify-ca。
	IsVerifyCa bool
	// 根证书路径。
	RootCertPath string
	// 客户端证书路径。
	ClientCertPath string
	// 客户端私钥路径。
	ClientKeyPath string
}

// ORM生成中和数据库相关定义。
type (
	// 在定义Database时，要添加这个匿名字段，用于生成代码。
	Database struct {
		DbInterface
	}
	// 数据库的配置。
	DbConfig struct {
		// 数据库名称。
		Name string
		// 连接的标签。
		Tag string
		// 数据库驱动
		Type dialect.DbDriver
	}
	DbInterface interface {
		Config() DbConfig
	}
)

func (Database) Config() DbConfig {
	return DbConfig{}
}

// ORM生成中实体表相关定义。
type (
	// 在定义Entity时，要添加这个匿名字段，用于生成代码。
	// 例如：
	// type UseEntity struct {
	//		entity.Entity
	// }
	//
	Entity struct {
		db Database
		EntityInterface
	}
	// 这个接口定义了Entity需要实现的方法。
	//
	// 这个接口在代码生成中会被调用，用于生成代码，
	// 见codengen/load/entity.go中的[Marshal]。
	EntityInterface interface {
		Config() EntityConfig
		Fields() []FieldBuilder
	}
	// 实体表配置。
	EntityConfig struct {
		// AttrName entity的数据库属性名，
		// 如果没有指定，会使用定义的结构体名称,但是会变成snake_case形式。
		//
		// 在codegen中会用于生成entity配置信息的文件和文件夹名,
		// 但是对于entity的结构体名字，还是使用定义的结构体名称，不使用AttrName，
		// 防止和别的database和entity名字冲突。
		AttrName string
		// Comment entity的注释。
		// 在生成的sql中会用于生成表的注释。
		Comment string
	}
)

func (Entity) Config() EntityConfig {
	return EntityConfig{}
}

func (Entity) Fields() []FieldBuilder {
	return nil
}

// ORM生成中实体表中的字段。
type (
	FieldValue driver.Value
	// 这个接口定义了字段在生成代码阶段需要的方法。
	FieldBuilder interface {
		// codegen中使用，用于初始化字段。
		Init(initDesc *Descriptor) error
		// codegen中使用，用于获取字段的描述。
		Descriptor() *Descriptor
		// codegen中使用，获取字段的数据库中的类型名，如果返回空字符串，会出现错误。
		// 如果dbType没有匹配的返回空字符串
		AttrType(dbType dialect.DbDriver) string
		// 用于设置字段的值的类型名称。例如entity.Int64的ValueType为"int64"。
		ValueType() string
	}
	// 这个接口定义了字段在运行时需要的方法。
	FieldStorager interface {
		// 用于扫描数据库返回的值，将值赋值给字段。
		Scan(value interface{}) error
		// 用于打印字段的值。
		String() string
		// 用于内部sql中获取字段的值。如果需要获得值，推荐通过Get()方法获得。
		Value() FieldValue
	}

	// 包含了关于字段的描述，配置信息等。
	// 这个在生成代码时会被调用。
	Descriptor struct {
		// Name 字段在结构体中的名字，这个会在codegen/load中通过Init被赋值。
		Name string `json:"name,omitempty"`
		// AttrName 字段的数据库属性名，
		// 如果为空，会使用Name的名字，,但是会变成snake_case形式
		AttrName string `json:"attr_name,omitempty"`
		// Type 字段的类型。如"entity.Int64"。
		Type string `json:"type,omitempty"`
		// AttrType 字段的数据库类型。如"entity.Int64"在PostgreSQL中对应"int8"，
		// 这AttrType的值为"int8"，这个通过AttrType()获得，所以自定义类型应该正确定义这个方法。
		AttrType string `json:"attr_type,omitempty"`
		// Size 字段的长度大小。
		Size int64 `json:"size,omitempty"`
		// Required 是否是必填字段，如果为true,在数据表中的表现就是这个字段非空。
		Required bool `json:"required,omitempty"`
		// Primary 字段是否为主键,大于等于1的才会被认为是主键。
		// 在生成的sql中Primary的值越小，越靠前，比如ID的Primary = 1，UUID的Primay = 2,
		// 则在sql中PRIMARY KEY (ID,UUID)会是这样
		Primary int `json:"primary,omitempty"`
		// Comment 字段的注释。
		Comment string `json:"comment,omitempty"`
		// Default 字段默认值。
		Default bool `json:"default,omitempty"`
		// DefaultValue 字段默认值的字符串形式。
		DefaultValue string `json:"default_value,omitempty"`
		// Locked 字段是否被锁定，如果为true,则不能被修改。
		Locked bool `json:"locked,omitempty"`
		// Sequence 字段的序列，
		// 不是所有的字段类型都可以设置序列，内置的类型中只有Int(Int16,Int32,Int64)
		// 才有Sequence()方法，自定义字段要看是否实现了设置序列的相关方法。
		Sequence Sequence `json:"validators,omitempty"`
		// Validators 字段验证函数。
		Validators []any `json:"sequence,omitempty"`
	}

	// Sequence 字段使用的序列，序列的类型默认为Int64。
	Sequence struct {
		// Name 序列的名称，不能为空字符串。
		Name *string
		// Increment 每次序列递增的值，默认1。
		Increament *int64
		// Min 序列的最小值，默认1。
		Min *int64
		// Max 序列的最大值，默认为9223372036854775807。
		Max *int64
		// Start 序列的起始值，默认1。
		Start *int64
		// Cache 指定序列中要预先分配的值的数量，默认1。
		Cache *int64
	}
)

// NewSequence 创建一个Sequence，name不能为空。
func NewSequence(name string) Sequence {
	if name == "" {
		panic("NewSequence name can't be empty")
	}
	increament := int64(1)
	min := int64(1)
	max := int64(9223372036854775807)
	start := int64(1)
	cache := int64(1)
	return Sequence{
		Name:       &name,
		Increament: &increament,
		Min:        &min,
		Max:        &max,
		Start:      &start,
		Cache:      &cache,
	}
}
