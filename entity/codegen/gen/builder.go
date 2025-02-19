package gen

import (
	"bytes"
	"fmt"
	"go/token"
	"go/types"
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
	// Builder 用于生成资源文件的构建器。
	Builder struct {
		*Config
		// Nodes是Schema中的database info的集合。
		Nodes     []*DatabaseInfo
		EntityMap load.EntityMap
	}

	// Generator 代码生成器接口。
	Generator interface {
		// Generate 为提供的Builder生成代码
		Generate(*Builder) error
	}

	// GenerateFunc 符合Generator接口的函数。
	GenerateFunc func(*Builder) error

	// Hook 代码生成的钩子。
	Hook func(Generator) Generator

	// tableError 生成代码时的错误。
	tableError struct {
		msg string
	}

	extInstanceTemplate struct {
		Paths []string
		Tmpl  InstanceTemplate
	}
)

// NewBuilder 根据提供的Schema，初始化一个生成器。
//
// Params:
//
//   - c: 代码生成的配置。
//   - entities: 从entity package中加载的所有entity。
//
// Returns:
//
//	0: 生成器。
//	1: 错误信息。
func NewBuilder(c *Config, databases ...*load.Database) (s *Builder, err error) {
	defer catch(&err)
	s = &Builder{Config: c, EntityMap: make(load.EntityMap)}
	for i := range databases {
		s.addNode(databases[i])
		for k, v := range databases[i].EntityMap {
			s.EntityMap[k] = v
		}
	}
	return s, nil
}

// Gen 调用符合Generator接口的函数生成代码。
func (b *Builder) Gen() error {
	var gen Generator = GenerateFunc(generate)
	for i := len(b.Hooks) - 1; i >= 0; i-- {
		gen = b.Hooks[i](gen)
	}
	return gen.Generate(b)
}

// templates 返回模版和外部模版。
//
// Params:
//
//   - dbType: 数据库类型。
//
// Returns:
//
//	0: 模版。
//	1: 外部模版。
func (b *Builder) templates(dbType dialect.DbDriver) (*template.Template, []extInstanceTemplate, map[string][]InstanceTemplate) {
	initTemplates(b, dbType)
	var (
		external       = make([]extInstanceTemplate, 0)
		field_external = make(map[string][]InstanceTemplate)
	)
	for _, node := range b.Nodes {
		entities := node.Database.Entities
		for _, entity := range entities {
			entityName := entity.AttrName
			for _, field := range entity.Fields {
				if field.Templates != nil && len(field.Templates) > 0 {
					attrName := field.AttrName
					fieldTmpls := template.NewTemplate("field_external")
					fieldTmpls.ParseFiles(field.Templates...)
					field_external[entityName+attrName] = []InstanceTemplate{}
					for _, t := range fieldTmpls.Templates() {
						name := t.Name()
						ext := filepath.Ext(name)
						if ext == "" {
							ext = ".go"
						} else {
							continue
						}
						field_external[entityName+attrName] = append(field_external[entityName+attrName], InstanceTemplate{
							Name: name,
							Format: func(t template.TemplatePathFormat) string {
								lastSlashIndex := strings.LastIndex(name, "/")
								if lastSlashIndex == -1 {
									return t.Dir() + "_" + attrName + "_" + name + ext
								}
								return filepath.Join(t.Dir(), attrName+"_"+name[lastSlashIndex+1:]+ext)
							},
						})
						if templates.Lookup(name) == nil {
							templates = template.MustParse(templates.AddParseTree(name, t.Tree))
						}
					}

				}
			}
		}
	}
	for _, extTmpl := range b.Templates {
		templates.Funcs(extTmpl.Tmpl.FuncMap)
		for _, tmpl := range extTmpl.Tmpl.Templates() {
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
			targetPaths := extTmpl.TargetPaths
			if len(targetPaths) == 0 {
				// 如果模版不是已经定义的模版或扩展，则生成一个新的文件。
				if templates.Lookup(name) == nil {
					external = append(external, extInstanceTemplate{
						Paths: targetPaths,
						Tmpl: InstanceTemplate{
							Name: name,
							Format: func(t template.TemplatePathFormat) string {
								lastSlashIndex := strings.LastIndex(name, "/")
								if lastSlashIndex == -1 {
									// 没有斜杠，直接添加前缀
									return "db_" + t.Dir() + "_" + name + ext
								}
								// 在最后一个斜杠之后添加前缀
								return name[:lastSlashIndex+1] + name[lastSlashIndex+1:] + ext
							},
						},
					},
					)
					templates = template.MustParse(templates.AddParseTree(name, tmpl.Tree))
				}
			} else {
				// 如果模版不是已经定义的模版或扩展，则生成一个新的文件。
				if templates.Lookup(name) == nil {
					external = append(external, extInstanceTemplate{
						Paths: targetPaths,
						Tmpl: InstanceTemplate{
							Name: name,
							Format: func(t template.TemplatePathFormat) string {
								// 在最后一个斜杠之后添加前缀
								return name + ext
							},
						},
					},
					)
					templates = template.MustParse(templates.AddParseTree(name, tmpl.Tree))
				}
			}
		}
	}
	return templates, external, field_external
}

// addNode 添加一个数据库节点到Builder中。
//
// Params:
//
//   - database: 数据库。
func (b *Builder) addNode(database *load.Database) {
	t, err := NewDatabaseInfo(b.Config, database)
	check(err, "create info %s", database.Name)
	b.Nodes = append(b.Nodes, t)
}

// Generate 生成代码。
//
// Params:
//
//   - f: 生成函数。
func (f GenerateFunc) Generate(t *Builder) error {
	return f(t)
}

// Error 返回错误信息。
// 实现了error接口。
func (t tableError) Error() string { return fmt.Sprintf("taurus_go/entity: %s", t.msg) }

// ValidSchemaName 确定一个名字是否会与任何预定义的名字冲突。
//
// Params:
//
//   - name: 名字。
func ValidSchemaName(name string) error {
	// Schema package is lower-cased (see Type.Package).
	pkg := strings.ToLower(name)
	if token.Lookup(pkg).IsKeyword() {
		return fmt.Errorf("schema lowercase name conflicts with Go keyword %q", pkg)
	}
	if types.Universe.Lookup(pkg) != nil {
		return fmt.Errorf("schema lowercase name conflicts with Go predeclared identifier %q", pkg)
	}
	return nil
}

// generate 默认的代码生成器实现。
//
// Params:
//
//   - t: 生成器。
func generate(t *Builder) error {
	var (
		assets asset.Assets
	)
	// 获取模版。
	// 为每个节点生成代码：
	for _, n := range t.Nodes {
		templates, extend, field_external := t.templates(n.Database.Type)
		for _, tmpl := range append(DatabaseTemplates) {
			if dir := filepath.Dir(tmpl.Format(n)); dir != "." {
				assets.AddDir(filepath.Join(t.Config.Target, dir))
			}
			b := bytes.NewBuffer(nil)
			if err := templates.ExecuteTemplate(b, tmpl.Name, n); err != nil {
				return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
			}
			assets.Add(filepath.Join(t.Config.Target, tmpl.Format(n)), b.Bytes())
		}
		for _, ext := range extend {
			b := bytes.NewBuffer(nil)
			if err := templates.ExecuteTemplate(b, ext.Tmpl.Name, n); err != nil {
				return fmt.Errorf("execute template %q: %w", ext.Tmpl.Name, err)
			}
			if len(ext.Paths) == 0 {
				if dir := filepath.Dir(ext.Tmpl.Format(n)); dir != "." {
					assets.AddDir(filepath.Join(t.Config.Target, dir))
				}
				assets.Add(filepath.Join(t.Config.Target, ext.Tmpl.Format(n)), b.Bytes())
			} else {
				for _, path := range ext.Paths {
					assets.AddDir(path)
					targetDir := filepath.Dir(filepath.Join(path, ext.Tmpl.Format(n)))
					assets.AddDir(targetDir)
					assets.Add(filepath.Join(path, ext.Tmpl.Format(n)), b.Bytes())
				}
			}
		}
		// 为节点的每个entity生成代码
		for _, e := range n.Database.Entities {
			entityName := e.AttrName
			ei, err := NewEntityInfo(t.Config, e, n.Database)
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
			for _, field := range e.Fields {
				fieldName := field.AttrName
				if field_external[entityName+fieldName] != nil {
					for _, tmpl := range field_external[entityName+fieldName] {
						b := bytes.NewBuffer(nil)
						fi, err := NewFieldInfo(t.Config, field, e, n.Database.Type)
						if err != nil {
							return err
						}
						if err := templates.ExecuteTemplate(b, tmpl.Name, fi); err != nil {
							return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
						}
						assets.Add(filepath.Join(t.Config.Target, tmpl.Format(ei)), b.Bytes())
					}
				}
			}
		}
	}
	templates, _, _ := t.templates("")
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
	if t.Config.ExtraCodes != nil {
		for _, tmpl := range ExtraCodesTemplates {
			if tmpl.Skip != nil && tmpl.Skip(t) {
				continue
			}
			b := bytes.NewBuffer(nil)
			if err := templates.ExecuteTemplate(b, tmpl.Name, t); err != nil {
				return fmt.Errorf("execute template %q: %w", tmpl.Name, err)
			}
			assets.Add(filepath.Join(t.Config.Target, tmpl.Format), b.Bytes())
		}
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
//
// Params:
//
//   - err: 错误。
//   - msg: 错误信息。
//   - args: 错误信息的参数。
func check(err error, msg string, args ...any) {
	if err != nil {
		args = append(args, err)
		panic(tableError{fmt.Sprintf(msg+": %s", args...)})
	}
}

// cleanOldNodes 清理在 entity 中已被删除但相关文件仍存在于文件系统中的节点（Node）的生成文件。
//
// Params:
//
//   - assets: 资源文件。
//   - target: 目标目录。
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
	var deleted []*DatabaseInfo
	for _, f := range d {
		if !strings.HasSuffix(f.Name(), "_query.go") {
			continue
		}
		typ := &DatabaseInfo{Database: &load.Database{Name: strings.TrimSuffix(f.Name(), ".go")}}
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
		for _, t := range DatabaseTemplates {
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
