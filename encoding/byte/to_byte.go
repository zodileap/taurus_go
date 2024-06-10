package byteutil

import (
	"strconv"
)

// IntToBytes 将 int 转换为 []byte。
//
// Params:
//
//   - n: int 值。
func IntToBytes(n int) []byte {
	b := make([]byte, 0)
	return strconv.AppendInt(b, int64(n), 10)
}

// Int8ToBytes 将 int8 转换为 []byte。
//
// Params:
//
//   - n: int8 值。
func Int8ToBytes(n int8) []byte {
	b := make([]byte, 0)
	return strconv.AppendInt(b, int64(n), 10)
}

// Int16ToBytes 将 int16 转换为 []byte。
//
// Params:
//
//   - n: int16 值。
func Int16ToBytes(n int16) []byte {
	b := make([]byte, 0)
	return strconv.AppendInt(b, int64(n), 10)
}

// Int32ToBytes 将 int32 转换为 []byte。
//
// Params:
//
//   - n: int32 值。
func Int32ToBytes(n int32) []byte {
	b := make([]byte, 0)
	return strconv.AppendInt(b, int64(n), 10)
}

// Int64ToBytes 将 int64 转换为 []byte。
//
// Params:
//
//   - n: int64 值。
func Int64ToBytes(n int64) []byte {
	b := make([]byte, 0)
	return strconv.AppendInt(b, n, 10)
}

// IntSToBytes 将 []int 转换为 []byte。
//
// Params:
//
//   - n: []int 值。
func IntSToBytes(n []int, sep string) []byte {
	b := make([]byte, 0)
	for _, num := range n {
		b = strconv.AppendInt(b, int64(num), 10)
		b = append(b, []byte(sep)...)
	}
	return b
}

// Int8SToBytes 将 []int8 转换为 []byte。
//
// Params:
//
//   - n: []int8 值。
func Int8SToBytes(n []int8, sep string) []byte {
	b := make([]byte, 0)
	for _, num := range n {
		b = strconv.AppendInt(b, int64(num), 10)
		b = append(b, []byte(sep)...)
	}
	return b
}

// Int16SToBytes 将 []int16 转换为 []byte。
//
// Params:
//
//   - n: []int16 值。
func Int16SToBytes(n []int16, sep string) []byte {
	b := make([]byte, 0)
	for _, num := range n {
		b = strconv.AppendInt(b, int64(num), 10)
		b = append(b, []byte(sep)...)
	}
	return b
}

// Int32SToBytes 将 []int32 转换为 []byte。
//
// Params:
//
//   - n: []int32 值。
func Int32SToBytes(n []int32, sep string) []byte {
	b := make([]byte, 0)
	for _, num := range n {
		b = strconv.AppendInt(b, int64(num), 10)
		b = append(b, []byte(sep)...)
	}
	return b
}

// Int64SToBytes 将 []int64 转换为 []byte。
//
// Params:
//
//   - n: []int64 值。
func Int64SToBytes(n []int64, sep string) []byte {
	b := make([]byte, 0)
	for _, num := range n {
		b = strconv.AppendInt(b, num, 10)
		b = append(b, []byte(sep)...)
	}
	return b
}

func StringToBytes(s string) []byte {
	return []byte(s)
}
