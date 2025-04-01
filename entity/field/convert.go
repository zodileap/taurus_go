package field

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	byteutil "github.com/zodileap/taurus_go/datautil/byte"
)

// errNilPtr 新建一个错误，表示目标指针为空。
var errNilPtr = errors.New("destination pointer is nil")

// convertAssign 将 src 中的值复制到 dest，如果可能的话进行转换。
// 如果复制会导致信息丢失，则会返回错误。
//
// Params:
//
//   - dest: 目标值。
//   - src: 源值。
func convertAssign(dest, src any) error {
	return convertAssignRows(dest, src)
}

// convertAssignRows 将 src 中的值复制到 dest，如果可能的话进行转换。
// 如果复制会导致信息丢失，则会返回错误。
// dest 应该是指针类型，否则会返回错误。
//
// Params:
//
//   - dest: 目标值。
//   - src: 源值。
func convertAssignRows(dest, src any) error {
	// 类型断言和赋值来实现.
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s)
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = append((*d)[:0], s...)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
			return nil
		case *any:
			if d == nil {
				return errNilPtr
			}
			*d = bytes.Clone(s)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = bytes.Clone(s)
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		default:
			return BytesToSlice(d, s)
		}
	case time.Time:
		switch d := dest.(type) {
		case *time.Time:
			*d = s
			return nil
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s.AppendFormat((*d)[:0], time.RFC3339Nano)
			return nil
		}
	case decimalDecompose:
		switch d := dest.(type) {
		case decimalCompose:
			return d.Compose(s.Decompose(nil))
		}
	case nil:
		switch d := dest.(type) {
		case *any:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		}
	// The driver is returning a cursor the client may iterate over.
	case driver.Rows:
		errors.New("invalid context to convert cursor rows, missing parent *Rows")
	}

	var sv reflect.Value

	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *RawBytes:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes([]byte(*d)[:0], sv); ok {
			*d = RawBytes(b)
			return nil
		}
	case *bool:
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool)
		}
		return err
	case *any:
		*d = src
		return nil
	}

	if scanner, ok := dest.(Scanner); ok {
		return scanner.Scan(src)
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Pointer {
		return errors.New("destination not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(bytes.Clone(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	// 以下转换使用字符串值作为中间表示在各种数字类型之间进行转换。
	// 这也允许扫描到用户定义的类型，如 "type Int int64"。为对称起见，也检查字符串目标类型。
	switch dv.Kind() {
	case reflect.Pointer:
		if src == nil {
			dv.SetZero()
			return nil
		}
		dv.Set(reflect.New(dv.Type().Elem()))
		return convertAssignRows(dv.Interface(), src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	case reflect.String:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", dv.Kind())
		}
		switch v := src.(type) {
		case string:
			dv.SetString(v)
			return nil
		case []byte:
			dv.SetString(string(v))
			return nil
		}
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

// strconvErr 返回 strconv.ParseInt 或 strconv.ParseUint 的错误。
func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}

func asString(src any) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}

type decimalDecompose interface {
	// Decompose returns the internal decimal state in parts.
	// If the provided buf has sufficient capacity, buf may be returned as the coefficient with
	// the value set and length set as appropriate.
	Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32)
}

type decimalCompose interface {
	// Compose sets the internal decimal value from parts. If the value cannot be
	// represented then an error should be returned.
	Compose(form byte, negative bool, coefficient []byte, exponent int32) error
}

type arrayToPGStringCallBack func(any) (string, error)

// arrayToPGString 将任何维度的数组转换为 PostgreSQL 数组格式的字符串
func arrayToPGString(value interface{}, callback arrayToPGStringCallBack) (string, error) {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return callback(value)
	}

	var result strings.Builder
	result.WriteString("{")
	for i := 0; i < val.Len(); i++ {
		if i > 0 {
			result.WriteString(", ")
		}
		element, err := arrayToPGString(val.Index(i).Interface(), callback)
		if err != nil {
			return "", err
		}
		result.WriteString(element)
	}
	result.WriteString("}")

	return result.String(), nil
}

type sliceType int

const (
	sliceTypeCommon sliceType = 0
	sliceTypeBool   sliceType = 1
	sliceTypeTime   sliceType = 2
	sliceTypeString sliceType = 3
)

// BytesToSlice 处理通用切片类型的转换
func BytesToSlice(dest any, src []byte) error {
	src = byteutil.ReplaceAll(src, 123, []byte{91})
	src = byteutil.ReplaceAll(src, 125, []byte{93})
	// 获取dest的反射值对象
	v := reflect.ValueOf(dest)

	// 检查传入的dest是否是指针
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	// 获取指针指向的元素
	e := v.Elem()
	if e.Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a slice pointer")
	}
	t := validSliceType(e)
	if t == sliceTypeBool {
		src = byteutil.ReplaceAll(src, 116, []byte{116, 114, 117, 101})
		src = byteutil.ReplaceAll(src, 102, []byte{102, 97, 108, 115, 101})
	} else if t == sliceTypeTime {
		var data interface{}
		if err := json.Unmarshal(src, &data); err != nil {
			return fmt.Errorf("slice parse error: %v", err)
		}
		if err := parseTimeRecursive(data, e.Addr().Interface().(*[]time.Time)); err != nil {
			return fmt.Errorf("slice parse error: %v", err)
		}
		return nil
	} else if t == sliceTypeString {
		stingSlice, err := bytesToStrings(src)
		if err != nil {
			return fmt.Errorf("slice parse error: %v", err)
		}
		e.Set(reflect.ValueOf(stingSlice))
		return nil
	}

	if err := json.Unmarshal(src, e.Addr().Interface()); err != nil {
		return fmt.Errorf("slice parse error: %v", err)
	}

	return nil
}

// validSliceType 检查切片的类型
func validSliceType(v reflect.Value) sliceType {
	for v.Kind() == reflect.Slice {
		v = reflect.New(v.Type().Elem()).Elem()
		// 当不再是切片时，检查是否为布尔类型
		if v.Kind() != reflect.Slice {
			if v.Kind() == reflect.Bool {
				return sliceTypeBool
			} else if v.Kind() == reflect.TypeOf(time.Time{}).Kind() {
				return sliceTypeTime
			} else if v.Kind() == reflect.String {
				return sliceTypeString
			} else {
				return sliceTypeCommon
			}
		}
	}
	return sliceTypeCommon
}

// parseTimeRecursive 递归解析时间字符串
func parseTimeRecursive(data interface{}, result *[]time.Time) error {
	// 检查数据类型
	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)

	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		// 迭代数组或切片中的每个元素
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()
			if err := parseTimeRecursive(item, result); err != nil {
				return err
			}
		}
	case reflect.String:
		// 解析字符串为 time.Time
		timeStr := data.(string)
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999-07", timeStr)
		if err != nil {
			return fmt.Errorf("error parsing time '%s': %w", timeStr, err)
		}
		*result = append(*result, parsedTime)
	default:
		return fmt.Errorf("unsupported data type: %s", rt.Kind())
	}

	return nil
}

// bytesToStrings 将形如 [91 97 98 44 98 99 93] 的字节数组转换为字符串数组
// 过程：去除方括号 -> 按逗号分割 -> 转换为字符串数组
func bytesToStrings(input []byte) ([]string, error) {
	// 检查输入是否有效
	if len(input) < 2 {
		return nil, fmt.Errorf("输入字节数组长度过短")
	}

	// 检查首尾是否为方括号
	if input[0] != '[' || input[len(input)-1] != ']' {
		return nil, fmt.Errorf("输入格式错误：需要以 [ 开头，] 结尾")
	}

	// 去除首尾的方括号
	content := input[1 : len(input)-1]

	// 如果去除方括号后为空，返回空数组
	if len(content) == 0 {
		return []string{}, nil
	}

	// 按逗号分割
	parts := bytes.Split(content, []byte{44}) // 44 是逗号的 ASCII 码

	// 转换每个部分为字符串
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		// 去除可能存在的空白字符
		part = bytes.TrimSpace(part)
		if len(part) > 0 {
			result = append(result, string(part))
		}
	}

	return result, nil
}
