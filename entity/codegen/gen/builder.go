package gen

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template/parse"

	"github.com/yohobala/taurus_go/asset"
	"github.com/yohobala/taurus_go/entity/codegen/load"
	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/template"
)

type (
	Builder struct {
		*Config
		// Databases 包含了从entity package中加载的所有database。
		Databases []*load.Database
		// Nodes是Schema中的database info的集合。
		Nodes []*Info
	}

	Generator interface {
		// Generate 为提供的table生成代码
		Generate(*Builder) error
	}

	GenerateFunc func(*Builder) error

	Hook func(Generator) Generator

	tableError struct {
		msg string
	}
)

// NewBuilder 根据提供的Entity信息，生成数据库的表。
//
// 参数：
//   - c: 代码生成的配置。
//   - entities: 从entity package中加载的所有entity。
func NewBuilder(c *Config, databases ...*load.Database) (s *Builder, err error) {
	defer catch(&err)
	s = &Builder{Config: c, Databases: databases}
	for i := range databases {
		s.addNode(databases[i])
	}
	return
}

func (b *Builder) Gen() error {
	var gen Generator = GenerateFunc(generate)
	for i := len(b.Hooks) - 1; i >= 0; i-- {
		gen = b.Hooks[i](gen)
	}
	return gen.Generate(b)
}

func (b *Builder) templates(dbType dialect.DbDriver) (*template.Template, []InstanceTemplate) {
	initTemplates(b, dbType)
	var (
		external = make([]InstanceTemplate, 0, len(b.Templates))
	)
	for _, extTmpl := range b.Templates {
		templates.Funcs(extTmpl.FuncMap)
		for _, tmpl := range extTmpl.Templates() {
			if parse.IsEmptyTree(tmpl.Root) {
				continue
			}
			name := tmpl.Name()
			ext := filepath.Ext(name)
			if ext == "" {
				ext = ".go"
			} else {
				name = strings.TrimSuffix(name, ext)
			}
			// 如果模版不是已经定义的模版或扩展，则生成一个新的文件。
			if templates.Lookup(name) == nil {
				external = append(external, InstanceTemplate{
					Name: name,
					Format: func(t template.TemplatePathFormat) string {
						lastSlashIndex := strings.LastIndex(name, "/")
						if lastSlashIndex == -1 {
							// 没有斜杠，直接添加前缀
							return "db_" + t.Dir() + "_" + name + ext
						}
						// 在最后一个斜杠之后添加前缀
						return name[:lastSlashIndex+1] + t.Dir() + "_" + name[lastSlashIndex+1:] + ext
					},
				})

				templates = template.MustParse(templates.AddParseTree(name, tmpl.Tree))
			}
		}
	}
	return templates, external
}

func (b *Builder) addNode(database *load.Database) {
	t, err := NewInfo(b.Config, database)
	check(err, "create info %s", database.Name)
	b.Nodes = append(b.Nodes, t)
}

func (f GenerateFunc) Generate(t *Builder) error {
	return f(t)
}

func (t tableError) Error() string { return fmt.Sprintf("taurus_go/entity: %s", t.msg) }

// 默认的代码生成器实现。
func generate(t *Builder) error {
	var (
		assets asset.Assets
	)
	// 获取模版。
	// 为每个节点生成代码：
	for _, n := range t.Nodes {
		templates, extend := t.templates(n.Database.Type)
		for _, tmpl := range append(Templates, extend...) {
			if dir := filepath.Dir(tmpl.Format(n)); dir != "." {
				assets.AddDir(filepath.Join(t.Config.Target, dir))
			}
			b := bytes.NewBuffer(nil)
			if err := templates.ExecuteTemplate(b, tmpl.Name, n); err != nil {
				return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
			}
			assets.Add(filepath.Join(t.Config.Target, tmpl.Format(n)), b.Bytes())
		}
		// 为节点的每个entity生成代码
		for _, e := range n.Database.Entities {
			ei, err := NewEntityInfo(t.Config, e)
			if err != nil {
				return err
			}
			assets.AddDir(filepath.Join(t.Config.Target, ei.Dir()))
			for _, tmpl := range EntityTemplates {
				b := bytes.NewBuffer(nil)
				if err := templates.ExecuteTemplate(b, tmpl.Name, ei); err != nil {
					return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
				}
				assets.Add(filepath.Join(t.Config.Target, tmpl.Format(ei)), b.Bytes())
			}
		}
	}
	templates, _ := t.templates("")
	// 通用的核心功能的模板:
	for _, tmpl := range InstanceTemplates {
		if tmpl.Skip != nil && tmpl.Skip(t) {
			continue
		}
		if dir := filepath.Dir(tmpl.Format); dir != "." {
			assets.AddDir(filepath.Join(t.Config.Target, dir))
		}
		b := bytes.NewBuffer(nil)
		if err := templates.ExecuteTemplate(b, tmpl.Name, t); err != nil {
			return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
		}
		assets.Add(filepath.Join(t.Config.Target, tmpl.Format), b.Bytes())
	}

	// 清理功能
	// for _, f := range AllFeatures {
	// 	if f.cleanup == nil || t.featureEnabled(f) {
	// 		continue
	// 	}
	// 	if err := f.cleanup(t.Config); err != nil {
	// 		return fmt.Errorf("cleanup %q feature assets: %w", f.Name, err)
	// 	}
	// }
	// 写入和格式化生成的代码。
	if err := assets.Write(); err != nil {
		return err
	}
	// 清理旧的节点和模板文件。
	cleanOldNodes(assets, t.Config.Target)
	for _, n := range deletedTemplates {
		if err := os.Remove(filepath.Join(t.Target, n)); err != nil && !os.IsNotExist(err) {
			log.Printf("remove old file %s: %s\n", filepath.Join(t.Target, n), err)
		}
	}
	return assets.Format()
}

// catch 错误处理，如果是tableError类型错误，会抛出panic
func catch(err *error) {
	if e := recover(); e != nil {
		terr, ok := e.(tableError)
		if !ok {
			panic(e)
		}
		*err = terr
	}
}

// check 如果err不是nil在抛出panic。
func check(err error, msg string, args ...any) {
	if err != nil {
		args = append(args, err)
		panic(tableError{fmt.Sprintf(msg+": %s", args...)})
	}
}

// 清理在 entity 中已被删除但相关文件仍存在于文件系统中的节点（Node）的生成文件。
func cleanOldNodes(assets asset.Assets, target string) {
	// 读取目标目录
	d, err := os.ReadDir(target)
	if err != nil {
		return
	}
	// 查找已删除的节点，
	// 如果一个文件以 _query.go 结尾，它可能是一个节点相关的文件。
	// 函数通过文件名推断出节点的类型（Type），
	// 并检查是否这个节点对应的目录仍存在于 assets.Dirs 中。如果不存在，这意味着节点可能已被删除。
	var deleted []*Info
	for _, f := range d {
		if !strings.HasSuffix(f.Name(), "_query.go") {
			continue
		}
		typ := &Info{Database: &load.Database{Name: strings.TrimSuffix(f.Name(), ".go")}}
		// 获取文件路径，并判断是否存在于assets.Dirs。
		path := filepath.Join(target, typ.Dir())
		if _, ok := assets.Dirs[path]; ok {
			continue
		}
		// 如果这个是一个节点，那它应该存在一个模型文件和一个目录（例如 ent/t.go, ent/t）。
		_, err1 := os.Stat(path + ".go")
		f2, err2 := os.Stat(path)
		if err1 == nil && err2 == nil && f2.IsDir() {
			deleted = append(deleted, typ)
		}
	}
	// 确认节点是否被删除。
	for _, typ := range deleted {
		for _, t := range Templates {
			err := os.Remove(filepath.Join(target, t.Format(typ)))
			if err != nil && !os.IsNotExist(err) {
				log.Printf("remove old file %s: %s\n", filepath.Join(target, t.Format(typ)), err)
			}
		}
		err := os.Remove(filepath.Join(target, typ.Dir()))
		if err != nil && !os.IsNotExist(err) {
			log.Printf("remove old dir %s: %s\n", filepath.Join(target, typ.Dir()), err)
		}
	}
}
