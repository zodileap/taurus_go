package internal

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yohobala/taurus_go/entity/codegen"
	"github.com/yohobala/taurus_go/entity/codegen/gen"
)

// GenerateCmd 生成Schema的资源文件，通过运行`github.com/yohobala/taurus_go/entity/cmd generate`调用。
//
// Returns:
//
//	0: "github.com/spf13/cobra"的Command对象。
func GenerateCmd() *cobra.Command {
	var (
		templates []string
		config    gen.Config
		cmd       = &cobra.Command{
			Use:     "generate [flags] path",
			Short:   "generate go code for the entity directory",
			Example: "go run -mod=mod github.com/yohobala/taurus_go/entity/cmd generate ./entity",
			// 要求至少有一个参数
			Args: cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, path []string) {
				// TODO: 目前只需要路径，没有别的flags
				exts := []codegen.Extra{}
				for _, tmpl := range templates {
					typ := "dir"
					if parts := strings.SplitN(tmpl, "=", 2); len(parts) > 1 {
						typ, tmpl = parts[0], parts[1]
					}
					switch typ {
					case "dir":
						exts = append(exts, codegen.TemplateDir(tmpl))
					case "file":
						exts = append(exts, codegen.TemplateFiles(tmpl))
					case "glob":
						exts = append(exts, codegen.TemplateGlob(tmpl))
					default:
						log.Fatalln("unsupported template type", typ)
					}
				}
				// 执行代码生成
				if err := codegen.Generate(path[0], &config, exts...); err != nil {
					log.Fatalln(err)
				}
			},
		}
	)
	cmd.Flags().StringSliceVarP(&templates, "template", "t", nil, "external templates to execute")

	return cmd
}
