package sliceutil

import (
	"fmt"
	"reflect"
)

// 将指定元素及其之前的所有元素移到数组末尾。这个函数会修改原始数组的顺序，
// 将目标元素之前的所有元素（包括目标元素）移动到数组的末尾。
//
// Example:
//
//	arr := []int{1, 2, 3, 4, 5}
//	result, _ := MoveElementToEndAndRemovePrevious(arr, 2)
//	// result = [4, 5, 1, 2, 3]
//
// ErrCodes:
//   - text_err_index_out_of_range
func MoveElementToEndAndRemovePrevious[T any](arr []T, elementIndex int) ([]T, error) {
	if elementIndex >= len(arr) || elementIndex < 0 {
		return nil, fmt.Errorf(text_err_index_out_of_range)
	}
	arr = arr[elementIndex:]
	elementIndex = 0
	element := arr[elementIndex]
	arr = append(arr[elementIndex+1:], element)

	return arr, nil
}

// 检查整数切片中是否包含指定的整数值。
//
// Example:
//
//	arr := []int{1, 2, 3, 4, 5}
//	exists := ContainByInt(arr, 3)
//	// exists = true
func ContainByInt(arr []int, element int) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

// 判断给定的接口值是否为切片或数组类型。使用反射来检查类型。
//
// Example:
//
//	arr := []int{1, 2, 3}
//	isSlice := IsSliceOrArray(arr)
//	// isSlice = true
func IsSliceOrArray(x interface{}) bool {
	kind := reflect.TypeOf(x).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

// 根据提供的过滤函数筛选切片中的元素。返回一个新的切片，
// 其中包含所有满足过滤条件的元素。
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	even := Filter(numbers, func(n int) bool {
//	    return n%2 == 0
//	})
//	// even = [2, 4]
func Filter[T any](slice []T, f func(T) bool) []T {
	filtered := make([]T, 0)
	for _, v := range slice {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// 在切片中查找第一个满足条件的元素。如果找到则返回该元素的指针，
// 如果未找到则返回 nil。
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	found, _ := Find(numbers, func(n int) bool {
//	    return n > 3
//	})
//	// found = &4
func Find[T any](slice []T, f func(T) bool) (*T, error) {
	for _, v := range slice {
		if f(v) {
			return &v, nil
		}
	}
	return nil, nil
}
