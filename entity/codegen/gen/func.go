package gen

import (
	"fmt"
	"sort"
	"strings"

	stringutil "github.com/yohobala/taurus_go/encoding/string"
	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/codegen/load"
	"github.com/yohobala/taurus_go/template"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// funcMap gen中模版需要用到的函数的映射
var funcMap = template.FuncMap{
	"joinFieldAttrNames":        joinFieldAttrNames,
	"joinFieldPrimaies":         joinFieldPrimaies,
	"joinRequiredFields":        joinRequiredFields,
	"joinFieldsString":          joinFieldsString,
	"getPrimaryField":           getPrimaryField,
	"snakeCaseToLowerCamelCase": snakeCaseToLowerCamelCase,
	"getRequiredFields":         getRequiredFields,
	"getEntityRel":              getEntityRel,
	"getEntityRelDirection":     getEntityRelDirection,
}

// joinFieldAttrNames 把字段的AttrName连接起来。
//
// Params:
//
//   - fs: 字段列表。
//
// Returns:
//
//	0: 拼接后的字符串。
func joinFieldAttrNames(fs []*load.Field) string {
	var ss []string
	for _, f := range fs {
		ss = append(ss, fmt.Sprintf(`'%s'`, f.AttrName))
	}
	return strings.Join(ss, ",")
}

// joinFieldPrimaies 把主键字段的AttrName连接起来。
//
// Params:
//
//   - fs: 字段列表。
//
// Returns:
//
//	0: 拼接后的字符串。
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
//
// Params:
//
//   - fs: 字段列表。
//   - param: 是否只返回参数名。
//
// Returns:
//
//	0: 拼接后的字符串。
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
//
// Params:
//
//   - fs: 字段列表。
//
// Returns:
//
//	0: 拼接后的字符串。
func joinFieldsString(fs []*load.Field) string {
	var ss []string
	for _, f := range fs {
		ss = append(ss, fmt.Sprintf(`%s: %%v`, f.Name))
	}
	return strings.Join(ss, ", ")
}

func getPrimaryField(fs []*load.Field) *load.Field {
	for _, f := range fs {
		if f.Primary == 1 {
			return f
		}
	}
	return nil
}

// getLowerCamelCase 获取小驼峰命名，会清除snake_case的下划线。
func snakeCaseToLowerCamelCase(a string) string {
	// 分割字符串为单词数组
	parts := strings.Split(a, "_")
	// 创建一个 Title 使用的 caser
	caser := cases.Title(language.English, cases.NoLower)
	// 处理每个单词，除了第一个单词保持小写，其他单词首字母大写
	for i := 1; i < len(parts); i++ {
		parts[i] = caser.String(parts[i])
	}
	// 将单词数组连接为一个字符串
	return strings.Join(parts, "")
}

// getRequiredFields 获取必填字段
//
// Params:
//
//   - fs: 字段列表。
//
// Returns:
//
//	0: 必填字段列表。
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

// getEntityRelField 获取关联实体的生成的结构体属性名和属性类型的字符串个
//
//	比如Author，会返回 Authors, rel.AuthorEntityRelation, []*AuthorEntity
//
// Params:
//
//   - rel: 关联关系。
//   - e: 实体。
//
// Returns:
//
//	0: 关联实体的属性名。
//	1: 关联属性的类型。
//	2: 关联实体的实体结构体名字
//	3: 实体类型
//	4: 原始的load.RelationEntity
type getEntityRelResult struct {
	Name       string
	AttrName   string
	RelType    string
	EntityType string
	Rel        load.RelationEntity
}

func getEntityRel(rel *load.Relation, e *load.Entity) *getEntityRelResult {
	if rel.Principal.AttrName == e.AttrName {
		if rel.Dependent.Rel == entity.O {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(strings.ReplaceAll(rel.Dependent.AttrName, "_", ""), "", 1),
				AttrName:   rel.Dependent.AttrName,
				RelType:    fmt.Sprintf("%sRelation", stringutil.ToUpperFirst(rel.Dependent.Name, "", 1)),
				EntityType: fmt.Sprintf("*%s", stringutil.ToUpperFirst(rel.Dependent.Name, "", 1)),
				Rel:        rel.Dependent,
			}
		} else {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(strings.ReplaceAll(rel.Dependent.AttrName, "_", ""), "", 1) + "s",
				AttrName:   rel.Dependent.AttrName,
				RelType:    fmt.Sprintf("%sRelation", stringutil.ToUpperFirst(rel.Dependent.Name, "", 1)),
				EntityType: fmt.Sprintf("[]*%s", stringutil.ToUpperFirst(rel.Dependent.Name, "", 1)),
				Rel:        rel.Dependent,
			}
		}
	} else if rel.Dependent.AttrName == e.AttrName {
		if rel.Principal.Rel == entity.O {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(strings.ReplaceAll(rel.Principal.Name, "_", ""), "", 1),
				AttrName:   rel.Principal.AttrName,
				RelType:    fmt.Sprintf("%sRelation", stringutil.ToUpperFirst(rel.Principal.Name, "", 1)),
				EntityType: fmt.Sprintf("*%s", stringutil.ToUpperFirst(rel.Principal.Name, "", 1)),
				Rel:        rel.Principal,
			}
		} else {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(strings.ReplaceAll(rel.Principal.Name, "_", ""), "", 1) + "s",
				AttrName:   rel.Principal.AttrName,
				RelType:    fmt.Sprintf("%sRelation", stringutil.ToUpperFirst(rel.Principal.Name, "", 1)),
				EntityType: fmt.Sprintf("[]*%s", stringutil.ToUpperFirst(rel.Principal.Name, "", 1)),
				Rel:        rel.Principal,
			}
		}
	}

	return nil
}

type getEntityRelDirectionResult struct {
	To   *load.RelationEntity
	Join *load.RelationEntity
}

func getEntityRelDirection(rel *load.Relation, e *load.Entity) getEntityRelDirectionResult {
	if rel.Principal.AttrName == e.AttrName {
		return getEntityRelDirectionResult{
			To:   &rel.Principal,
			Join: &rel.Dependent,
		}
	} else if rel.Dependent.AttrName == e.AttrName {
		return getEntityRelDirectionResult{
			To:   &rel.Dependent,
			Join: &rel.Principal,
		}
	}
	return getEntityRelDirectionResult{}
}
