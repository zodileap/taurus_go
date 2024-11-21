package template

import (
	"errors"
	"path/filepath"
	"strings"
	"text/template"

	stringutil "github.com/yohobala/taurus_go/datautil/string"
)

var (
	// Funcs are the predefined template
	// functions used by the codegen.
	Funcs = template.FuncMap{
		"filePathBase":       filepath.Base,
		"createMap":          dict,
		"stringToLower":      toLower,
		"stringToFirstCap":   toFirstCap,
		"stringToFirstLower": toFirstLower,
		"stringToSnakeCase":  stringutil.ToSnakeCase,
		"stringReplace":      strings.Replace,
		"stringHasPrefix":    strings.HasPrefix,
		"stringReplaceAll":   strings.ReplaceAll,
		"stringSub":          sub,
		"stringJoin":         StringJoin,
		"stringIndex":        strings.Index,
		"stringSplice":       splice,
		"toString":           toString,
	}
	acronyms = make(map[string]struct{})
)

// dict 创建一个map，key/value对的列表。
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// toLower 把字符串转换为小写。
func toLower(s string) string {
	if _, ok := acronyms[s]; ok {
		return s
	}
	return strings.ToLower(s)
}

func toFirstCap(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func toFirstLower(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

func sub(a, b int) int {
	return a - b
}

// StringJoin 把字符串列表连接起来。
func StringJoin(ss ...string) string {
	return strings.Join(ss, "")
}

// toString 将uint8转换为string
func toString(b uint8) string {
	return string(b)
}

// splice 从字符串指定位置截取指定长度
func splice(s string, start int, length int) string {
	runes := []rune(s)
	if start < 0 {
		start = len(runes) + start // 负数表示从末尾开始
	}
	if start >= len(runes) {
		return s
	}
	if length < 0 {
		length = len(runes) - start + length
	}
	if start+length > len(runes) {
		length = len(runes) - start
	}
	return string(runes[start : start+length])
}
