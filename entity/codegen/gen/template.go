package gen

import (
	"embed"
	"fmt"

	"github.com/zodileap/taurus_go/entity/dialect"
	"github.com/zodileap/taurus_go/template"
)

type (
	// InstanceTemplate 模版实例
	InstanceTemplate = template.FileTemplate[any]

	// GenericTemplate 是生成代码中给每个InstanceTemplate使用的通用部分
	GenericTemplate struct {
		Name   string
		Format string
		Skip   func(*Builder) bool
	}
)

var (
	// DatabaseTemplates 保存database将要生成的文件的模版信息
	DatabaseTemplates []InstanceTemplate = []InstanceTemplate{
		{
			Name:   "database",
			Format: pkgf("d_%s.go"),
		},
		{
			Name:   "sql/table",
			Format: pkgf("sql/%s.sql"),
		},
	}
	// EntityTemplates entity的相关模版
	EntityTemplates []InstanceTemplate = []InstanceTemplate{
		{
			Name:   "entity/builder",
			Format: pkgf("e_%s_builder.go"),
		},
		{
			Name:   "entity/entity",
			Format: pkgf("e_%s_entity.go"),
		},
		{
			Name:   "entity/fields",
			Format: pkgf("e_%s_fields.go"),
		},
		{
			Name:   "entity/create",
			Format: pkgf("e_%s_create.go"),
		},
		{
			Name:   "entity/delete",
			Format: pkgf("e_%s_delete.go"),
		},
		{
			Name:   "entity/query",
			Format: pkgf("e_%s_query.go"),
		},
		{
			Name:   "entity/update",
			Format: pkgf("e_%s_update.go"),
		},
		{
			Name: "entity/meta",
			Format: func(t template.TemplatePathFormat) string {
				return fmt.Sprintf("%[1]s/%[1]s.go", t.Dir())
			},
		},
		{
			Name:   "entity/where",
			Format: pkgf("%s/where.go"),
		},
		{
			Name:   "entity/order",
			Format: pkgf("%s/order.go"),
		},
		{
			Name:   "rel/entity",
			Format: pkgf("e_%s_rel.go"),
		},
	}
	// InstanceTemplates 内部使用的模版
	InstanceTemplates []GenericTemplate = []GenericTemplate{
		{
			Name:   "internal/core",
			Format: "internal/core.go",
		},
		{
			Name:   "rel/rels",
			Format: "rels.go",
		},
	}
	// ExtraCodesTemplates 额外的代码模版
	ExtraCodesTemplates []GenericTemplate = []GenericTemplate{
		{
			Name:   "extraCode",
			Format: "extra_codes.go",
		},
	}

	// SqlTemplates sql文件的模版
	SqlTemplates []GenericTemplate = []GenericTemplate{}
	templates    *template.Template
	//go:embed template/*
	templateDir embed.FS

	deletedTemplates = []string{"config.go", "context.go"}
	importPkg        = map[string]string{}
)

// initTemplates 初始化模版。
//
// Params:
//
//   - builder: 生成资源文件的构建器。
//   - dbType: 数据库类型。
func initTemplates(builder *Builder, dbType dialect.DbDriver) {
	// 根据数据库类型选择模版
	if dbType == dialect.PostgreSQL {
		templates = template.MustParse(template.NewTemplate("templates").
			Funcs(FuncMap).
			ParseFS(templateDir,
				"template/*.tmpl",
				"template/internal/*.tmpl",
				"template/postgresql/*.tmpl",
				"template/postgresql/sql/*.tmpl",
				"template/postgresql/entity/*.tmpl",
				"template/rel/*.tmpl",
			))
	} else {
		templates = template.MustParse(template.NewTemplate("templates").
			Funcs(FuncMap).
			ParseFS(templateDir,
				"template/*.tmpl",
				"template/internal/*.tmpl",
				"template/rel/*.tmpl",
			))
	}

}

// pkgf 返回一个格式化的路径
//
// Params:
//
//   - s: 格式化字符串。
//
// Returns:
//
//	0: 实现了TemplatePathFormat接口的字符串。
func pkgf(s string) func(t template.TemplatePathFormat) string {
	return func(t template.TemplatePathFormat) string {
		return fmt.Sprintf(s, t.Dir())
	}
}
