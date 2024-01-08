package stringutil

import (
	"fmt"
	"strings"
	"unicode"
)

func MoveElementToEndAndRemovePrevious[T any](arr []T, elementIndex int) ([]T, error) {
	arr = arr[elementIndex:]
	elementIndex = 0
	element := arr[elementIndex]
	arr = append(arr[elementIndex+1:], element)

	return arr, nil

}

// 无重复字符的最长子串
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
func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func IsFormatString(s string) bool {
	return strings.Contains(s, "%")
}

func IsUpper(s string) bool {
	if s == "" {
		return false
	}
	r := []rune(s)
	return unicode.IsUpper(r[0])
}
