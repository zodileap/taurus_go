package load

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"text/template"

	entity "github.com/yohobala/taurus_go/entity"
)

type (
	// ExeTmplConfig 执行模板的配置。
	ExeTmplConfig struct {
		// Config 模版文件会匹配Config每个字段的值，例如[.Names]
		Config *Config
		// Package 模版文件会匹配[.Package]的值
		Package string
	}
)

var (
	//go:embed template/main.tmpl tmpl_entity.go
	files     embed.FS
	buildTmpl = templates()
)

// templates 加载所有需要的模板文件。
func templates() *template.Template {
	tmpls, err := entityTemplates()
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(template.New("templates").
		ParseFS(files, "template/main.tmpl"))
	for _, t := range tmpls {
		tmpl = template.Must(tmpl.Parse(t))
	}
	return tmpl
}

// entityTemplates 从tmpl_entity.go中解析出模板。
// 把tmpl_entity.go中的代码分成两部分，一部分是导入路径，一部分是代码。
//
// Returns:
//
//		0: 一个字符串切片，第一个元素是代码，第二个元素是导入路径。
//	 1: 错误信息。
//
// Example:
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go。
//
// ErrCodes:
//
//   - Err_0200010001。
func entityTemplates() ([]string, error) {
	var (
		// 文件中的所有导入路径
		imports []string
		// 文件中的所有代码
		code   bytes.Buffer
		fset   = token.NewFileSet()
		src, _ = files.ReadFile("tmpl_entity.go")
	)
	f, err := parser.ParseFile(fset, "tmpl_entity.go", src, parser.AllErrors)
	if err != nil {
		return nil, entity.Err_0100020016.Sprintf(err)
	}
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.IMPORT {
			for _, spec := range decl.Specs {
				imports = append(imports, spec.(*ast.ImportSpec).Path.Value)
			}
			continue
		}
		// 格式化代码，并添加到code中
		if err := format.Node(&code, fset, decl); err != nil {
			return nil, entity.Err_0100020017.Sprintf(err)
		}
		code.WriteByte('\n')
	}
	return []string{
		fmt.Sprintf(`{{ define "entity" }} %s {{ end }}`, code.String()),
		fmt.Sprintf(`{{ define "imports" }} %s {{ end }}`, strings.Join(imports, "\n")),
	}, nil
}
