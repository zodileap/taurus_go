package stringutil

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"unicode"
)

func MoveElementToEndAndRemovePrevious[T any](arr []T, elementIndex int) ([]T, error) {
	arr = arr[elementIndex:]
	elementIndex = 0
	element := arr[elementIndex]
	arr = append(arr[elementIndex+1:], element)

	return arr, nil
}

// GenerateKey 生成一个唯一的字符串key。
// 类似于"github.com/google/uuid"，
// 但不具备像UUID那样的强大的唯一性保证和标准格式。
func GenerateKey() (string, error) {
	// 当前时间的纳秒作为一部分
	currentTime := time.Now().UnixNano()

	// 生成一些随机字节作为key的一部分
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// 将时间和随机字节转换为字符串
	key := fmt.Sprintf("%x-%x", currentTime, randomBytes)
	return key, nil
}

// LengthOfLongestSubstring 无重复字符的最长子串
// 采用滑动窗口方式时间复杂度O(n),如果暴力解法时间复杂度为O(n^2)
// 这里想象一个会伸缩的窗口在字符串中，然后一个个移动过去，
// 如果窗口中有重复的字符，就把窗口的左边界移动到重复字符的下一个位置
// 否则窗口右边界向右移动一格
// 这段代码增加了对UTF-8字符的支持，
// 源代码如下，只支持ASCII字符,但是效率更高
// start := 0
// end := 0
//
//	for i, v := range s {
//		index := strings.Index(string(s[start:i]), string(v))
//		if index == -1 {
//			if i+1 > end {
//				end = i + 1
//			}
//		} else {
//			start += index + 1
//			end += index + 1
//		}
//	}
//
// return end - start
//
// 参数：
// - s: 字符串
func LengthOfLongestSubstring(s string) int {
	runes := []rune(s)
	start := 0
	end := 0
	count := -1
	for _, v := range s {
		count++
		fmt.Printf("%v , %v \n", string(runes[start:count]), string(v))
		index := strings.Index(string(runes[start:count]), string(v))
		if index == -1 {
			if count+1 > end {
				end = count + 1
			}

		} else {
			start += index + 1
			end += index + 1
		}

		fmt.Printf("index: %v \n", index)
		fmt.Printf("start: %v \n", start)
		fmt.Printf("end: %v \n", end)
		fmt.Printf("count: %v \n", count)
	}

	return end - start
}

// ToSnakeCase 函数将字符串转换为蛇形命名法（snake_case）。
// 它将所有字符转换为小写，并在大写字母前添加下划线，第一个字符除外。
// 例如，"FooBar" -> "foo_bar"。
// "APIDoc" -> "api_doc"
// "SimpleXMLParser" -> "simple_xml_parser"
func ToSnakeCase(str string) string {
	var result strings.Builder
	runes := []rune(str)
	for i, r := range runes {
		// 如果是第一个字符，直接转小写
		if i == 0 {
			result.WriteRune(unicode.ToLower(r))
			continue
		}

		// 当前字符是大写时
		if unicode.IsUpper(r) {
			// 判断是否需要添加下划线：
			// 1. 前一个字符是小写
			// 2. 不是最后一个字符，且后一个字符是小写
			if unicode.IsLower(runes[i-1]) ||
				(i+1 < len(runes) && unicode.IsLower(runes[i+1])) {
				result.WriteRune('_')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// 函数将字符串转换为驼峰命名法（camelCase）。
// 支持将蛇形命名法和大写驼峰命名法转换为小写驼峰命名法。
// 例如，"foo_bar" -> "fooBar"。
// "FooBar" -> "fooBar"
func ToCamelCase(str string) string {
	var result strings.Builder
	nextUpper := false

	// 确定首字母是否应该大写（大写驼峰）
	firstChar := true

	for _, r := range str {
		// 处理分隔符：星号和下划线
		if r == '_' {
			nextUpper = true
		} else {
			if firstChar {
				// 首字母小写（小写驼峰）
				result.WriteRune(unicode.ToLower(r))
				firstChar = false
			} else if nextUpper {
				// 分隔符后的字符大写
				result.WriteRune(unicode.ToUpper(r))
				nextUpper = false
			} else {
				result.WriteRune(unicode.ToLower(r))
			}
		}
	}
	return result.String()
}

// ToUpperFirst 函数将字符串转换为大写字母开头的字符串。
//
// Params:
// - s: 字符串
// - sep: 分隔符，如果为空，则不分割
// - num: 需要转换的首字母大写的数量，如果为-1，则全部转换,如果为0，则不转换
func ToUpperFirst(s string, sep string, num int) string {
	if num == 0 {
		return s
	}

	var split []string
	n := 0
	var ns []string = make([]string, 0)
	if s != "" {
		split = strings.Split(s, sep)
	} else {
		split = []string{s}
	}
	for _, v := range split {
		if v == "" {
			continue
		}
		if num == -1 || (num > 0 && n < num) {
			ns = append(ns, strings.ToUpper(v[:1])+v[1:])
			n += 1
		} else {
			ns = append(ns, v)
		}
	}
	return strings.Join(ns, sep)
}

// IsFormatString 检查字符串是否包含格式化字符串。
func IsFormatString(s string) bool {
	return strings.Contains(s, "%")
}

// IsUpper 函数检查字符串的第一个字符是否为大写字母。
func IsUpper(s string) bool {
	if s == "" {
		return false
	}
	r := []rune(s)
	return unicode.IsUpper(r[0])
}

// numberToLetters 将数字转换为字母。
// 例如，0 -> "A"，1 -> "B"，...，25 -> "Z"，26 -> "AA"，27 -> "AB"，...。
func NumberToLetters(n int) string {
	result := ""
	for n >= 0 {
		// 计算当前位置的字母
		result = string(rune('A'+(n%26))) + result
		n = n/26 - 1
	}
	return result
}
