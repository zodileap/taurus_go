package internal

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/zodileap/taurus_go/asset"
	stringutil "github.com/zodileap/taurus_go/datautil/string"
	entity "github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/codegen/gen"
	"github.com/zodileap/taurus_go/template"
)

const defaultEntity = "entity"
const defaultSchema = "schema"
const defaultDatabase = "db.go"
const defaultGenerate = "generate.go"

//go:embed template/*
var templateDir embed.FS

// NewCmd 新建Schema命令, 通过运行`github.com/zodileap/taurus_go/entity/cmd new`调用。
//
// Returns:
//
//	0: "github.com/spf13/cobra"的Command对象。
func NewCmd() *cobra.Command {
	var target string
	var schema string
	var entities []string
	cmd := &cobra.Command{
		Use:     "new [entities]",
		Short:   "initialize a new environment with zero or more entities",
		Example: "go run -mod=mod github.com/zodileap/taurus_go/entity/cmd new User Group",
		Args: func(_ *cobra.Command, names []string) error {
			for _, name := range names {
				if !unicode.IsUpper(rune(name[0])) {
					log.Fatalln(entity.Err_0100020002.Sprintf(name))
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, names []string) {
			var (
				err  error
				tmpl *template.Template
			)
			tmpl = template.NewTemplate("entity")
			tmpl, err = tmpl.ParseFS(templateDir, "template/new.tmpl")
			if err != nil {
				log.Fatalln(err)
			}
			if err := newDB(target, schema, names, tmpl); err != nil {
				log.Fatalln(err)
			}
			if err := newEntity(target, schema, entities, tmpl); err != nil {
				log.Fatalln(err)
			}
		},
	}
	cmd.Flags().StringVarP(&target, "target", "t", defaultEntity, "target directory for schemas, defaults to entity")
	cmd.Flags().StringVar(&schema, "schema", defaultSchema, "schema package name")
	cmd.Flags().StringSliceVarP(&entities, "entities", "e", []string{}, "entities to create")
	return cmd
}

// newDB 在Schema中创建数据库。
//
// Params:
//
//   - target: 目标目录。
//   - schema: Schema的文件夹名称。
//   - names: 数据库名称。
//   - tmpl: 模板文件。
//
// ErrCodes:
//   - Err_0200010001。
//   - Err_0200010003。
func newDB(target string, schema string, names []string, tmpl *template.Template) error {
	var (
		assets    asset.Assets
		b         *bytes.Buffer
		existFile bool = false
	)
	assets.AddDir(target)
	assets.AddDir(filepath.Join(target, schema))
	if err := assets.Write(); err != nil {
		return err
	}
	dbFile := filepath.Join(target, schema, defaultDatabase)
	if !asset.FileExists(dbFile) {
		b = bytes.NewBuffer(nil)
		assets.Add(filepath.Join(target, schema, defaultDatabase), nil)
		if err := tmpl.ExecuteTemplate(b, "schema/package", schema); err != nil {
			return entity.Err_0100020003.Sprintf("new.tmpl", err)
		}
	} else {
		existFile = true
		var err error
		b, err = asset.ReadFileToBuffer(dbFile)
		if err != nil {
			return err
		}
	}
	for _, name := range names {
		if existFile {
			exist, err := hasStructInFile(dbFile, name)
			if err != nil {
				return err
			}
			if exist {
				fmt.Printf("skip ==> database %s already exists.\n", name)
				continue
			}
		}
		if err := gen.ValidSchemaName(name); err != nil {
			return entity.Err_0100020001.Sprintf(name, err)
		}
		if err := tmpl.ExecuteTemplate(b, "database", name); err != nil {
			return entity.Err_0100020003.Sprintf("database", err)
		}
	}
	assets.Add(filepath.Join(target, schema, defaultDatabase), b.Bytes())
	genByte := bytes.NewBuffer(nil)
	if err := tmpl.ExecuteTemplate(genByte, "generate", defaultEntity); err != nil {
		return entity.Err_0100020003.Sprintf("generate", err)
	}
	assets.Add(filepath.Join(target, defaultGenerate), genByte.Bytes())
	if err := assets.Write(); err != nil {
		return err
	}
	return assets.Format()
}

// newEntity 在Schema中创建实体。
//
// Params:
//
//   - target: 目标目录。
//   - schema: Schema的文件夹名称。
//   - names: 实体名称。
//   - tmpl: 模板文件。
//
// ErrCodes:
//   - Err_0200010001。
//   - Err_0200010003。
func newEntity(target string, schema string, names []string, tmpl *template.Template) error {
	var assets asset.Assets
	for _, name := range names {
		var b *bytes.Buffer
		entityFile := filepath.Join(target, schema, stringutil.ToSnakeCase(name)+".go")
		if !asset.FileExists(entityFile) {
			b = bytes.NewBuffer(nil)
			if err := tmpl.ExecuteTemplate(b, "schema/package", schema); err != nil {
				return entity.Err_0100020003.Sprintf("schema/package", err)
			}
			if err := tmpl.ExecuteTemplate(b, "entity", name); err != nil {
				return entity.Err_0100020003.Sprintf("entity", err)
			}
			assets.Add(entityFile, b.Bytes())
		} else {
			fmt.Printf("skip ==> entity %s already exists.\n", name)
			continue
		}
	}
	if err := assets.Write(); err != nil {
		return err
	}
	return assets.Format()
}

// hasStructInFile 检查给定的 Go 文件中是否存在指定的结构体名称。
//
// Params:
//
//   - filename: Go 文件名。
//   - structName: 结构体名称。
//
// Returns:
//
//   - 是否存在。
//   - 错误信息。
func hasStructInFile(filename, structName string) (bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if typeSpec.Name.Name == structName {
				return true, nil
			}
		}
	}

	return false, nil
}
