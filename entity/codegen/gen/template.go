package gen

import (
	"embed"
	"fmt"

	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/template"
)

type (
	// InstanceTemplate 根据数据库结构生成一个模版
	InstanceTemplate = template.FileTemplate[any]

	// GenericTemplate 是生成代码中给每个InstanceTemplate使用的通用部分
	GenericTemplate struct {
		Name   string
		Format string
		Skip   func(*Builder) bool
	}
)

var (
	// Templates 保存database将要生成的文件的模版信息
	Templates = []InstanceTemplate{
		{
			Name:   "database",
			Format: pkgf("%s.go"),
		},
		{
			Name:   "sql/table",
			Format: pkgf("sql/%s_table.sql"),
		},
	}
	EntityTemplates = []InstanceTemplate{
		{
			Name:   "entity/builder",
			Format: pkgf("%s/builder.go"),
		},
		{
			Name:   "entity/entity",
			Format: pkgf("%s/entity.go"),
		},
		{
			Name: "entity/meta",
			Format: func(t template.TemplatePathFormat) string {
				return fmt.Sprintf("%[1]s/%[1]s.go", t.Dir())
			},
		},
		{
			Name:   "entity/fields",
			Format: pkgf("%s/fields.go"),
		},
		{
			Name:   "entity/create",
			Format: pkgf("%s/create.go"),
		},
		{
			Name:   "entity/delete",
			Format: pkgf("%s/delete.go"),
		},
		{
			Name:   "entity/query",
			Format: pkgf("%s/query.go"),
		},
		{
			Name:   "entity/update",
			Format: pkgf("%s/update.go"),
		},
		{
			Name:   "entity/where",
			Format: pkgf("%s/where.go"),
		},
	}
	InstanceTemplates = []GenericTemplate{
		{
			Name:   "internal/core",
			Format: "internal/core.go",
		},
	}
	SqlTemplates = []GenericTemplate{}
	templates    *template.Template
	//go:embed template/*
	templateDir embed.FS

	deletedTemplates = []string{"config.go", "context.go"}
	importPkg        = map[string]string{}
)

// 初始化模版
func initTemplates(builder *Builder, dbType dialect.DbDriver) {
	// 解析模板
	if dbType == dialect.PostgreSQL {
		templates = template.MustParse(template.NewTemplate("templates").
			Funcs(funcMap).
			ParseFS(templateDir,
				"template/*.tmpl",
				"template/internal/*.tmpl",
				"template/postgresql/*.tmpl",
				"template/postgresql/sql/*.tmpl",
				"template/postgresql/entity/*.tmpl",
			))
	} else {
		templates = template.MustParse(template.NewTemplate("templates").
			Funcs(funcMap).
			ParseFS(templateDir,
				"template/*.tmpl",
				"template/internal/*.tmpl",
			))
	}

}

func pkgf(s string) func(t template.TemplatePathFormat) string {
	return func(t template.TemplatePathFormat) string {
		return fmt.Sprintf(s, t.Dir())
	}
}
