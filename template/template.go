package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"
)

func NewTemplate(name string) *Template {
	t := &Template{Template: template.New(name)}
	return t.Funcs(Funcs)
}

// 一个辅助函数，用于封装对返回(*Template, error)的调用，如果err不为空，则panic
func MustParse(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

type (
	Template struct {
		*template.Template
		FuncMap template.FuncMap
	}

	FuncMap = template.FuncMap
)

// Funcs 把funcMap中的函数添加到模版中
func (t *Template) Funcs(funcMap FuncMap) *Template {
	t.Template.Funcs(funcMap)
	if t.FuncMap == nil {
		t.FuncMap = make(template.FuncMap)
	}
	for name, f := range funcMap {
		if _, ok := t.FuncMap[name]; !ok {
			t.FuncMap[name] = f
		}
	}
	return t
}

// ParseDir 解析目录下的所有文件
func (t *Template) ParseDir(path string) (*Template, error) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk path %s: %w", path, err)
		}
		if info.IsDir() || strings.HasSuffix(path, ".go") {
			return nil
		}
		_, err = t.ParseFiles(path)
		return err
	})
	return t, err
}

// ParseFiles 解析文件
func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	if _, err := t.Template.ParseFiles(filenames...); err != nil {
		return nil, err
	}
	return t, nil
}

// ParseGlob 解析glob模式下的文件
func (t *Template) ParseGlob(pattern string) (*Template, error) {
	if _, err := t.Template.ParseGlob(pattern); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Template) ParseFS(fsys fs.FS, patterns ...string) (*Template, error) {
	if _, err := t.Template.ParseFS(fsys, patterns...); err != nil {
		return nil, err
	}
	return t, nil
}

// AddParseTree adds the given parse tree to the template.
func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error) {
	if _, err := t.Template.AddParseTree(name, tree); err != nil {
		return nil, err
	}
	return t, nil
}

type (
	FileTemplate[T any] struct {
		Name   string
		Format func(TemplatePathFormat) string
		Skip   func(*T) bool
	}

	TemplatePathFormat interface {
		Dir() string
	}
)

func Dir(s string) func(t TemplatePathFormat) string {
	return func(t TemplatePathFormat) string {
		return fmt.Sprintf(s, t.Dir())
	}
}
