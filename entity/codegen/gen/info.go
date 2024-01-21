package gen

import (
	"strings"

	stringutil "github.com/yohobala/taurus_go/encoding/string"
	"github.com/yohobala/taurus_go/entity/codegen/load"
)

type (
	// Info 表示一个Builder中的一个节点的信息，它所包含的信息
	Info struct {
		*Config
		Database *load.Database
	}

	EntityInfo struct {
		*Config
		Entity *load.Entity
	}
)

// 从提供的database中创建一个Info
func NewInfo(c *Config, database *load.Database) (*Info, error) {
	typ := &Info{
		Config:   c,
		Database: database,
	}
	return typ, nil
}

func NewEntityInfo(c *Config, entity *load.Entity) (*EntityInfo, error) {
	typ := &EntityInfo{
		Config: c,
		Entity: entity,
	}
	return typ, nil
}

// PackageDir 返回包目录名称
func (t Info) Dir() string {
	return strings.ToLower(t.Database.Name)
}

// PackageDir 返回包目录名称
func (t EntityInfo) Dir() string {
	return stringutil.ToSnakeCase(t.Entity.AttrName)
}
