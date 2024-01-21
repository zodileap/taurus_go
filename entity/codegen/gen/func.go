package gen

import (
	"fmt"
	"sort"
	"strings"

	stringutil "github.com/yohobala/taurus_go/encoding/string"
	"github.com/yohobala/taurus_go/entity/codegen/load"
	"github.com/yohobala/taurus_go/template"
)

var funcMap = template.FuncMap{
	"joinFieldAttrNames": joinFieldAttrNames,
	"joinFieldPrimaies":  joinFieldPrimaies,
	"joinRequiredFields": joinRequiredFields,
	"joinFieldsString":   joinFieldsString,
	"getRequiredFields":  getRequiredFields,
}

// joinFieldAttrNames 把字段的AttrName连接起来。
func joinFieldAttrNames(fs []*load.Field) string {
	var ss []string
	for _, f := range fs {
		ss = append(ss, fmt.Sprintf(`'%s'`, f.AttrName))
	}
	return strings.Join(ss, ",")
}

// joinFieldPrimaies 把主键字段的AttrName连接起来。
func joinFieldPrimaies(fs []*load.Field) string {
	var fields []*load.Field
	for _, f := range fs {
		if f.Primary > 0 {
			fields = append(fields, f)
		}
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Primary < fields[j].Primary
	})
	var ss []string
	for _, f := range fields {
		ss = append(ss, fmt.Sprintf(`"%s"`, f.AttrName))
	}
	return strings.Join(ss, ",")
}

// joinRequiredFields 把没有默认值但是是必填的字段拼接成方法的参数。用于New
func joinRequiredFields(fs []*load.Field, param bool) string {
	params := []string{}
	for _, f := range fs {
		if !f.Default {
			if f.Required {
				var s string
				if param {
					s = fmt.Sprintf(`%s`, stringutil.ToSnakeCase(f.AttrName))
				} else {
					s = fmt.Sprintf(`%s %s`, stringutil.ToSnakeCase(f.AttrName), f.ValueType)
				}
				params = append(params, s)
			}
		}
	}
	s := strings.Join(params, ", ")
	if len(params) > 0 {
		s += ","
	}
	return s
}

// joinFieldsString 把全部字段拼接成format,用于String()
func joinFieldsString(fs []*load.Field) string {
	var ss []string
	for _, f := range fs {
		ss = append(ss, fmt.Sprintf(`%s: %%v`, f.Name))
	}
	return strings.Join(ss, ", ")
}

func getRequiredFields(fs []*load.Field) []*load.Field {
	var fields []*load.Field
	for _, f := range fs {
		if !f.Default {
			if f.Required {
				fields = append(fields, f)
			}
		}
	}
	return fields
}
