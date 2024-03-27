package codegen

import (
	"go/token"
	"path"
	"path/filepath"
	"strings"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/codegen/gen"
	"github.com/yohobala/taurus_go/entity/codegen/load"
	"github.com/yohobala/taurus_go/template"
)

// Extra 用于在Config中添加额外的配置的回调函数。
type Extra func(*gen.Config) error

// Generate 在entityPath目录下运行代码生成器。
// 如果 entity的路径是 `<project>/entity`,则生成的代码将放在 `<project>/entity`目录下。
//
// Params:
//
//   - entityPath: entity package的路径。例如`<project>/entity`。
//   - cfg: 代码生成的配置
//   - options: 代码生成的选项,用于回调函数
func Generate(entityPath string, cfg *gen.Config, options ...Extra) error {
	// 设置目标路径：如果 cfg.Target 为空，则计算 entityPath 的绝对路径，并将其设置为代码生成的默认目标路径。
	// 例如，如果 entityPath 是 "<project>/entity"，则代码生成的目标路径将是 "<project>/entity"。
	if cfg.Target == "" {
		abs, err := filepath.Abs(entityPath)
		if err != nil {
			return err
		}
		// 默认的生成代码路径是entityPath目录的上一级
		//
		// 修改这个代码，将影响到生成代码的路径
		cfg.Target = filepath.Dir(abs)
	}

	for _, opt := range options {
		if err := opt(cfg); err != nil {
			return err
		}
	}
	// 调用 gen.PrepareEnv
	undo, err := gen.PrepareEnv(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = undo()
		}
	}()
	return generate(entityPath, cfg)
}

// LoadBuilder 根据提供的路径加载entity package,并生成一个schema。
//
// Params:
//
//   - entityPath: entity package的路径
//   - cfg: 代码生成的配置
//
// Returns:
//
//   - 代码构建器
//   - 错误信息
func LoadBuilder(entityPath string, cfg *gen.Config) (*gen.Builder, error) {
	builder, err := (&load.Config{Path: entityPath, BuildFlags: cfg.BuildFlags}).Load()
	if err != nil {
		return nil, err
	}
	if cfg.Package == "" {
		// 如果 cfg.Package 为空（即未指定代码生成的目标包路径），
		// 则将其设置为 builder.PkgPath 的父目录。
		// builder.PkgPath如果为"<project>/entity/schema"，
		// 则 cfg.Package 为"<project>/entity"，
		//
		// 修改这个将影响生成代码的 Go package path
		cfg.Package = path.Dir(builder.PkgPath)
	}
	return gen.NewBuilder(cfg, builder.Databases...)
}

// generate 生成代码。
// 首先加载提供的entity,得到一个用来生成代码的schem，
// 然后用schema生成代码.
//
// Params:
//
//   - entityPath: entity package的路径
//   - cfg: 代码生成的配置
func generate(entityPath string, cfg *gen.Config) error {
	builder, err := LoadBuilder(entityPath, cfg)
	if err != nil {
		return err
	}
	if err := normalizePkg(cfg); err != nil {
		return err
	}
	return builder.Gen()
}

// TemplateDir 解析目录类型的模版。
//
// Params:
//
//   - path: 模版目录的路径
//
// Returns:
//
//	0: Extra函数。
func TemplateDir(path string) Extra {
	return templateExt(func(t *template.Template) (*template.Template, error) {
		return t.ParseDir(path)
	})
}

// TemplateFiles 解析文件类型的模版
//
// Params:
//
//   - filenames: 模版文件的路径
//
// Returns:
//
//	0: Extra函数。
func TemplateFiles(filenames ...string) Extra {
	return templateExt(func(t *template.Template) (*template.Template, error) {
		return t.ParseFiles(filenames...)
	})
}

// TemplateGlob 解析glob模式下的文件
//
// Params:
//
//   - pattern: glob模式
//
// Returns:
//
//	0: Extra函数。
func TemplateGlob(pattern string) Extra {
	return templateExt(func(t *template.Template) (*template.Template, error) {
		return t.ParseGlob(pattern)
	})
}

// templateExt 生成一个Extra函数，用于解析模版。
//
// Params:
//
//   - next: 解析模版的函数
//
// Returns:
//
//	0: 模版。
//	1: 错误信息
func templateExt(next func(t *template.Template) (*template.Template, error)) Extra {
	return func(cfg *gen.Config) (err error) {
		tmpl, err := next(template.NewTemplate("external"))
		if err != nil {
			return err
		}
		cfg.Templates = append(cfg.Templates, tmpl)
		return nil
	}
}

// normalizePkg 检查包名是否合法，标准化包名。
// 如果包名中包含"-"，则将其替换为"_"。
// 如果包名不是合法的标识符，则返回错误。
//
// Params:
//
//   - c: 代码生成的配置
func normalizePkg(c *gen.Config) error {
	base := path.Base(c.Package)
	if strings.ContainsRune(base, '-') {
		base = strings.ReplaceAll(base, "-", "_")
		c.Package = path.Join(path.Dir(c.Package), base)
	}
	if !token.IsIdentifier(base) {
		return entity.Err_0100020015.Sprintf(base)
	}
	return nil
}
