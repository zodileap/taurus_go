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
)

type (
	ExeTmplConfig struct {
		// 模版文件会匹配Config每个字段的值，例如{{.Names}}}
		*Config
		// 模版文件会匹配{{.Package}}的值
		Package string
	}
)

var (
	//go:embed template/main.tmpl tmpl_entity.go
	files     embed.FS
	buildTmpl = templates()
)

func templates() *template.Template {
	tmpls, err := entityTemplates()
	if err != nil {
		panic(err)
	}
	// 从embed.FS中加载模板文件
	tmpl := template.Must(template.New("templates").
		ParseFS(files, "template/main.tmpl"))
	for _, t := range tmpls {
		tmpl = template.Must(tmpl.Parse(t))
	}
	return tmpl
}

// 从entity.go中解析出模板。
// 把entity.go中的代码分成两部分，一部分是导入路径，一部分是代码。
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
		return nil, fmt.Errorf("parse entity file: %w", err)
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
			return nil, fmt.Errorf("format node: %w", err)
		}
		code.WriteByte('\n')
	}
	return []string{
		fmt.Sprintf(`{{ define "entity" }} %s {{ end }}`, code.String()),
		fmt.Sprintf(`{{ define "imports" }} %s {{ end }}`, strings.Join(imports, "\n")),
	}, nil
}
