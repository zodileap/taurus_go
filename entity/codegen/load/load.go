package load

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yohobala/taurus_go/cmd"
	"github.com/yohobala/taurus_go/entity"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

type (
	// Config 用于从Schema中加载的所有database和entity的配置。
	Config struct {
		// Path 加载的Schema的路径。
		Path string
		// Entities 加载的Schema中拥有匿名字段[entity.Entity]的结构体的名称。
		Entities []string
		// BuildFlags 传递给go build的标志。
		BuildFlags []string
		// Dbs 加载的Schema中拥有匿名字段[entity.Database]的结构体的名称。
		Dbs []DbConfig
	}

	// DbConfig 是一个包含了Schema的database的信息。
	DbConfig struct {
		// Name 数据库名称。
		Name string
		// Entities 存储这个database所拥有的。
		Entities EntityMap
	}

	// BuilderInfo 是一个用于生成代码的构建器，
	// 包含了Schema的Go package路径和符合条件的database信息。
	BuilderInfo struct {
		// PkgPath 是加载的Schema包的Go package路径，之后会传给gen.Config。
		PkgPath string
		// Module 加载的entity package的模块信息。
		Module *packages.Module
		// Databases 从Schema中提取出的database的配置信息。
		// Config.Dbs中的只是获得数据库的名字和它拥有的entity。Databases获得的是完整的数据库配置信息。
		Databases []*Database
		// ExtraCodes Schema中不是数据库或者实体的代码。
		ExtraCodes []string
	}

	// EntityMap entity的key和类型。
	//
	// 和Config.Entities不同的是，
	// EntityMap是用于记录database中的entity的信息，
	// 而Config.Entities是用于记录entity结构体的名字。
	// 例如：
	// type User struct {
	// 	entity.Database
	// 	User UserEntity
	// }
	// 则这个EntityMap中的内容为：{
	// 	"User": "UserEntity"
	// }
	// 而Config.Entities中的内容为：["UserEntity"]
	EntityMap map[string]string
)

var (
	// entityInterface保存了[entity.Interface]的[reflect.Type]。
	entityInterface = reflect.TypeOf(struct{ entity.EntityInterface }{}).Field(0).Type
	// 保存了[entity.DbInterface]的[reflect.Type]。
	dbInterface = reflect.TypeOf(struct{ entity.DbInterface }{}).Field(0).Type
)

// Load 加载Schema，并且利用这些信息生成一个Builder。
//
// Returns:
//
//	0: 生成代码的构建器。
//	1: 错误信息。
//
// ErrCodes:
//
//   - Err_0100020004
//   - Err_0100020005
func (c *Config) Load() (*BuilderInfo, error) {
	// 获取传入路径下的entity信息。
	builder, err := c.load()
	if err != nil {
		return nil, entity.Err_0100020004.Sprintf(err)
	}
	if len(c.Entities) == 0 {
		return nil, entity.Err_0100020005.Sprintf(c.Path)
	}
	// 执行模版。
	var b bytes.Buffer
	err = buildTmpl.ExecuteTemplate(&b, "main", ExeTmplConfig{
		Config:  c,
		Package: builder.PkgPath,
	})
	if err != nil {
		return nil, entity.Err_0100020003.Sprintf(err)
	}
	// 格式化生成的代码，并创建目录和文件，最后写入到文件中。
	buf, err := format.Source(b.Bytes())
	if err != nil {
		return nil, entity.Err_0100020006.Sprintf(err)
	}
	if err := os.MkdirAll(".gen", os.ModePerm); err != nil {
		return nil, entity.Err_0100020007.Sprintf(err)
	}
	target := fmt.Sprintf(".gen/%s.go", filename(builder.PkgPath))
	if err := os.WriteFile(target, buf, 0644); err != nil {
		return nil, entity.Err_0100020008.Sprintf(target, err)
	}
	// 清理加载文件。
	defer os.RemoveAll(".gen")
	// 运行生成的代码，解析代码输出，得到entity。
	out, err := cmd.RunGo(target, c.BuildFlags)
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(out, "\n") {
		database, err := Unmarshal([]byte(line))
		if err != nil {
			return nil, entity.Err_0100020009.Sprintf(line, err)
		}
		builder.Databases = append(builder.Databases, database)
	}
	return builder, nil
}

// load 加载传入的路径中符合要求的Entity、database的信息，
// 通过这些信息创建一个Builder。
//
// Returns:
//
//	0: 生成代码的构建器。
//	1: 错误信息。
//
// ErrCodes:
//
//   - Err_0100020010
//   - Err_0100020011
//   - Err_0100020012
func (c *Config) load() (*BuilderInfo, error) {
	// 加载指定路径的go包
	// pkgs是一个包的切片，切片的元素数量取决于传入的路径里包含的包的数量
	// 一般来说就2个，一个是c.Path，一个是entityInterface.PkgPath()所在的包。
	pkgs, err := packages.Load(&packages.Config{
		BuildFlags: c.BuildFlags,
		// Load函数需要返回的包的信息
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule | packages.NeedSyntax,
	}, c.Path, entityInterface.PkgPath())
	if err != nil {
		return nil, entity.Err_0100020010.Sprintf(err)
	}
	if len(pkgs) < 2 {
		// 检查数量少于2是否是因为 "Go-related"引起的错误
		if err := cmd.List(c.Path, c.BuildFlags); err != nil {
			return nil, err
		}
		return nil, entity.Err_0100020011.Sprintf(c.Path)
	}
	entPkg, loadPkg := pkgs[0], pkgs[1]
	if len(loadPkg.Errors) != 0 {
		return nil, c.loadError(loadPkg.Errors[0])
	}
	if len(entPkg.Errors) != 0 {
		return nil, entPkg.Errors[0]
	}
	// 判断是否是entity接口的包，如果不是翻转。
	if pkgs[0].PkgPath != entityInterface.PkgPath() {
		entPkg, loadPkg = pkgs[1], pkgs[0]
	}
	var names []string
	// 这部分代码是检查，加载的代码中是否有实现了 entity.Entity 接口的结构体。
	// 获取 ent 接口类型：
	iface := entPkg.Types.Scope().Lookup(entityInterface.Name()).Type().Underlying().(*types.Interface)
	var dbs []DbConfig
	dbIface := entPkg.Types.Scope().Lookup(dbInterface.Name()).Type().Underlying().(*types.Interface)

	// 这个循环遍历用户定义的包（loadPkg）中的所有类型定义。
	// loadPkg.TypesInfo.Defs 包含了包中所有类型的定义，其中 k 是定义的标识符（如类型名称），v 是定义本身（如类型信息）。
	for k, v := range loadPkg.TypesInfo.Defs {
		// 这里检查定义 v 是否是一个命名类型（如结构体或接口）。
		typ, ok := v.(*types.TypeName)
		// 如果 v 不是命名类型，或者 k（标识符）不是导出的（即不是公开的），
		// 或者类型 typ没有实现entityInterface接口，则跳过当前迭代。
		if !ok || !k.IsExported() || (!types.Implements(typ.Type(), iface) && !types.Implements(typ.Type(), dbIface)) {
			continue
		}
		// 这里尝试将类型的声明（k.Obj.Decl）断言为 *ast.TypeSpec 类型，这是 Go 语言抽象语法树（AST）中表示类型声明的结构。
		spec, ok := k.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			return nil, entity.Err_0100020012.Sprintf(k.Obj.Decl, k.Name)
		}
		// 这里检查声明的类型（spec.Type）是否是结构体类型。
		structType, ok := spec.Type.(*ast.StructType)
		if !ok {
			return nil, entity.Err_0100020013.Sprintf(spec.Type, k.Name)
		}

		if types.Implements(typ.Type(), iface) {
			names = append(names, k.Name)
			continue
		}
		if types.Implements(typ.Type(), dbIface) {
			entities := EntityMap{}
			// 遍历结构体的每个字段。
			for _, field := range structType.Fields.List {
				for _, fieldName := range field.Names {
					fieldType := loadPkg.TypesInfo.TypeOf(fieldName)
					if fieldType != nil && types.Implements(fieldType, iface) {
						// 这里可以处理实现了iface接口的字段
						if ident, ok := field.Type.(*ast.Ident); ok {
							entities[fieldName.Name] = ident.Name
						} else {
							return nil, entity.Err_0100020014.Sprintf(field.Type, k.Name)
						}
					}
				}
			}
			dbs = append(dbs, DbConfig{
				Name:     k.Name,
				Entities: entities,
			})
			continue
		}
	}
	if len(c.Entities) == 0 {
		c.Entities = names
	}
	if len(c.Dbs) == 0 {
		c.Dbs = dbs
	}
	sort.Strings(c.Entities)

	// 收集Schema中额外的代码
	var extraCodes []string
	for _, file := range loadPkg.Syntax {
		// 使用文件级别的遍历来确保我们可以正确地处理 GenDecl 节点
		for _, decl := range file.Decls {
			switch n := decl.(type) {
			case *ast.GenDecl:
				var shouldAdd bool
				// 对于 GenDecl，我们需要检查它是否包含不符合条件的类型或常量
				for _, spec := range n.Specs {
					switch s := spec.(type) {
					// Type
					case *ast.TypeSpec:
						typ := loadPkg.TypesInfo.Defs[s.Name].(*types.TypeName)
						if s.Name.IsExported() && !types.Implements(typ.Type(), iface) && !types.Implements(typ.Type(), dbIface) {
							shouldAdd = true
							break
						}
					// Const 和 Var
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.IsExported() {
								shouldAdd = true
								break
							}
						}
					}
				}
				if shouldAdd {
					addNonConformingCode(n, loadPkg.Fset, &extraCodes)
				}
			case *ast.FuncDecl:
				// 检查函数是否符合条件
				if n.Recv != nil && len(n.Recv.List) > 0 {
					// 获取函数接收器的类型
					if recvType, ok := n.Recv.List[0].Type.(*ast.Ident); ok {
						recvTypeName := recvType.Name
						var recvTypeDefinition *types.TypeName
						for _, def := range loadPkg.TypesInfo.Defs {
							if typeName, ok := def.(*types.TypeName); ok && typeName.Name() == recvTypeName {
								recvTypeDefinition = typeName
								break
							}
						}
						if recvTypeDefinition != nil {
							if types.Implements(recvTypeDefinition.Type(), iface) || types.Implements(recvTypeDefinition.Type(), dbIface) {
								// 接收器类型符合条件，跳过该函数
								continue
							}
						}
						if recvType.IsExported() {
							addNonConformingCode(n, loadPkg.Fset, &extraCodes)
							continue
						}
					}
				} else if n.Name.IsExported() {
					// 函数没有接收器，直接检查函数名是否导出
					addNonConformingCode(n, loadPkg.Fset, &extraCodes)
				}

			}
		}
	}
	// 打印不符合条件的源代码
	for _, code := range extraCodes {
		fmt.Println(code)
	}

	return &BuilderInfo{PkgPath: loadPkg.PkgPath, Module: loadPkg.Module, ExtraCodes: extraCodes}, nil
}

// loadError 用于处理加载错误。
//
// Params:
//
//   - perr: 加载错误。
func (c *Config) loadError(perr packages.Error) (err error) {
	if strings.Contains(perr.Msg, "import cycle not allowed") {
		if cause := c.cycleCause(); cause != "" {
			perr.Msg += "\n" + cause
		}
	}
	err = perr
	if perr.Pos == "" {
		// Strip "-:" prefix in case of empty position.
		err = errors.New(perr.Msg)
	}
	return err
}

// cycleCause 检测在给定的 Go 代码包中是否存在可能导致循环依赖的本地类型声明。
//
// Returns:
//
//	0: 可能导致循环依赖的本地类型声明。
func (c *Config) cycleCause() (cause string) {
	// 解析代码目录。
	dir, err := parser.ParseDir(token.NewFileSet(), c.Path, nil, 0)
	// 如果出现解析 错误或无软件包可解析时，忽略报告。
	if err != nil || len(dir) == 0 {
		return
	}
	//查找包含entity的软件包，如果这个操作失败（pkg == nil），则取目录中的第一个包。
	pkg := dir[filepath.Base(c.Path)]
	if pkg == nil {
		for _, v := range dir {
			pkg = v
			break
		}
	}
	// 收集包内的本地类型声明。
	locals := make(map[string]bool)
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			g, ok := d.(*ast.GenDecl)
			if !ok || g.Tok != token.TYPE {
				continue
			}
			// 遍历包内的所有文件和声明，收集所有公开（exported）的非结构体类型声明。
			// 如果是结构体，遍历结构体的字段来检查是否嵌入了特定的类型。
			for _, s := range g.Specs {
				ts, ok := s.(*ast.TypeSpec)
				if !ok || !ts.Name.IsExported() {
					continue
				}
				// 不是结构体的类型如 "type Role int".
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					locals[ts.Name.Name] = true
					continue
				}
				var embedEntity bool
				astutil.Apply(st.Fields, func(c *astutil.Cursor) bool {
					f, ok := c.Node().(*ast.Field)
					if ok {
						switch x := f.Type.(type) {
						case *ast.SelectorExpr:
							if x.Sel.Name == "Entity" {
								embedEntity = true
							}
						case *ast.Ident:
							if name := strings.ToLower(x.Name); name == "entity" {
								embedEntity = true
							}
						}
					}
					return !embedEntity
				}, nil)
				if !embedEntity {
					locals[ts.Name.Name] = true
				}
			}
		}
	}
	if len(locals) == 0 {
		return
	}
	// 检查 entity 字段中的本地类型使用情况。
	goTypes := make(map[string]bool)
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			f, ok := d.(*ast.FuncDecl)
			if !ok || f.Name.Name != "Fields" || f.Type.Params.NumFields() != 0 || f.Type.Results.NumFields() != 1 {
				continue
			}
			astutil.Apply(f.Body, func(cursor *astutil.Cursor) bool {
				i, ok := cursor.Node().(*ast.Ident)
				if ok && locals[i.Name] {
					goTypes[i.Name] = true
				}
				return true
			}, nil)
		}
	}
	names := make([]string, 0, len(goTypes))
	for k := range goTypes {
		names = append(names, strconv.Quote(k))
	}
	sort.Strings(names)
	if len(names) > 0 {
		cause = fmt.Sprintf("To resolve this issue, move the custom types used by the generated code to a separate package: %s", strings.Join(names, ", "))
	}
	return
}

// filename 生成一个唯一的文件名。
//
// Params:
//
//   - pkg: Go package路径。
//
// Returns:
//
//	0: 文件名。
func filename(pkg string) string {
	name := strings.ReplaceAll(pkg, "/", "_")
	return fmt.Sprintf("gen_%s_%d", name, time.Now().Unix())
}

// 辅助函数，用于添加不符合条件的代码
func addNonConformingCode(node ast.Node, fset *token.FileSet, codeList *[]string) {
	startPos := fset.Position(node.Pos()).Offset
	endPos := fset.Position(node.End()).Offset
	fileContent, err := ioutil.ReadFile(fset.File(node.Pos()).Name())
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}
	sourceCode := string(fileContent[startPos:endPos])
	*codeList = append(*codeList, sourceCode)
}

// 辅助函数，用于检查切片中是否包含指定元素
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
