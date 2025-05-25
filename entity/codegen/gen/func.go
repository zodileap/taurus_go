package gen

import (
	"fmt"
	"sort"
	"strings"

	stringutil "github.com/zodileap/taurus_go/datautil/string"
	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/codegen/load"
	"github.com/zodileap/taurus_go/template"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FuncMap gen中模版需要用到的函数的映射
var FuncMap = template.FuncMap{
	"joinFieldAttrNames":        joinFieldAttrNames,
	"joinFieldPrimaies":         joinFieldPrimaies,
	"joinRequiredFields":        joinRequiredFields,
	"joinFieldsString":          joinFieldsString,
	"getPrimaryField":           getPrimaryField,
	"snakeCaseToLowerCamelCase": snakeCaseToLowerCamelCase,
	"getRequiredFields":         getRequiredFields,
	"getEntityRel":              getEntityRel,
	"getEntityRelDirection":     getEntityRelDirection,
	"getIndexGroups":            getIndexGroups,
	"getIndexMethod":            getIndexMethod,
	"stringFirstField":          stringFirstField,
	"stringJoinIndexFields":     stringJoinIndexFields,
	"stringJoinIndexColumns":    stringJoinIndexColumns,
	"stringJoinQuotedColumns":   stringJoinQuotedColumns,
	"getUniqueGroups":           getUniqueGroups,
	"getUniqueFieldGroups":      getUniqueFieldGroups,
	"removeArrayBrackets":       removeArrayBrackets,
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

func convertRelName(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func getEntityRel(rel *load.Relation, e *load.Entity) *getEntityRelResult {
	if rel.Principal.AttrName == e.AttrName {
		if rel.Dependent.Rel == entity.O {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(convertRelName(rel.Dependent.AttrName), "", 1),
				AttrName:   rel.Dependent.AttrName,
				RelType:    fmt.Sprintf("%sRelation", strings.ToLower(stringutil.ToUpperFirst(rel.Dependent.Name, "", 1))),
				EntityType: fmt.Sprintf("*%s", stringutil.ToUpperFirst(rel.Dependent.Name, "", 1)),
				Rel:        rel.Dependent,
			}
		} else {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(convertRelName(rel.Dependent.AttrName), "", 1) + "s",
				AttrName:   rel.Dependent.AttrName,
				RelType:    fmt.Sprintf("%sRelation", strings.ToLower(stringutil.ToUpperFirst(rel.Dependent.Name, "", 1))),
				EntityType: fmt.Sprintf("[]*%s", stringutil.ToUpperFirst(rel.Dependent.Name, "", 1)),
				Rel:        rel.Dependent,
			}
		}
	} else if rel.Dependent.AttrName == e.AttrName {
		if rel.Principal.Rel == entity.O {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(convertRelName(rel.Principal.AttrName), "", 1),
				AttrName:   rel.Principal.AttrName,
				RelType:    fmt.Sprintf("%sRelation", strings.ToLower(stringutil.ToUpperFirst(rel.Principal.Name, "", 1))),
				EntityType: fmt.Sprintf("*%s", stringutil.ToUpperFirst(rel.Principal.Name, "", 1)),
				Rel:        rel.Principal,
			}
		} else {
			return &getEntityRelResult{
				Name:       stringutil.ToUpperFirst(convertRelName(rel.Principal.AttrName), "", 1) + "s",
				AttrName:   rel.Principal.AttrName,
				RelType:    fmt.Sprintf("%sRelation", strings.ToLower(stringutil.ToUpperFirst(rel.Principal.Name, "", 1))),
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

// stringFirstField 获取切片中的第一个元素
func stringFirstField(fields []string) string {
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
}

// stringJoinIndexFields 拼接字段名
func stringJoinIndexFields(fields []string) string {
	return strings.Join(fields, "_")
}

// stringJoinIndexColumns 拼接字段为逗号分隔的列表
func stringJoinIndexColumns(fields []string) string {
	return strings.Join(fields, ", ")
}

// stringJoinQuotedColumns 将列名添加引号并用逗号连接
func stringJoinQuotedColumns(fields []string) string {
	quoted := make([]string, len(fields))
	for i, field := range fields {
		quoted[i] = fmt.Sprintf("%q", field)
	}
	return strings.Join(quoted, ", ")
}

// getIndexGroups 获取需要创建索引的字段分组
func getIndexGroups(fields []*load.Field) map[int][]string {
	groups := make(map[int][]string)
	// 先处理所有的多字段索引
	for _, field := range fields {
		for _, idx := range field.Indexes {
			if existingFields, ok := groups[idx]; !ok || len(existingFields) == 0 {
				// 如果这个索引号还没有字段，添加为单字段索引
				groups[idx] = []string{field.AttrName}
			} else {
				// 如果这个索引号已经有字段了，说明是联合索引的一部分
				groups[idx] = append(groups[idx], field.AttrName)
			}
		}
	}
	return groups
}

// getIndexMethod 获取索引方法
func getIndexMethod(fields []*load.Field, index int) string {
	for _, field := range fields {
		// 检查字段是否参与这个索引
		for _, idx := range field.Indexes {
			if idx == index && field.IndexMethod != "" && field.IndexMethod != "btree" {
				return fmt.Sprintf("USING %s", field.IndexMethod)
			}
		}
	}
	return "" // 默认btree时返回空字符串，不输出USING子句
}

// getUniqueGroups 获取唯一约束分组
func getUniqueGroups(fields []*load.Field) map[int][]string {
	groups := make(map[int][]string)
	for _, field := range fields {
		for _, idx := range field.Uniques {
			groups[idx] = append(groups[idx], field.AttrName)
		}
	}
	return groups
}

// getUniqueFieldGroups 获取实体中的唯一键字段，按组返回
func getUniqueFieldGroups(entity *load.Entity) map[int][]*load.Field {
	groups := make(map[int][]*load.Field)
	for _, field := range entity.Fields {
		for _, uniqueIdx := range field.Uniques {
			groups[uniqueIdx] = append(groups[uniqueIdx], field)
		}
	}
	return groups
}

// removeArrayBrackets 根据深度移除字符串前面的[]前缀
func removeArrayBrackets(s string, depth int) string {
	result := s
	for i := 0; i < depth; i++ {
		result = strings.TrimPrefix(result, "[]")
	}
	return result
}
