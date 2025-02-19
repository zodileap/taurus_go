package sliceutil

import (
	"fmt"
	"reflect"
)

// 从数组中移除指定元素之前的全部元素，并将该元素放到数组末尾
//
// 例如：数组[1,2,3,4,5]，移除元素3，结果为[4,5,1,2,3]
//
// 参数：
//   - arr: 需要修改的数组
//   - elementIndex: 需要移除的元素的索引
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

// 寻找两个有序数组的中位数
//
// 例如：数组1[1,3]，数组2[2]，中位数为2.0
//
// 参数：
//   - nums1: 数组1
//   - nums2: 数组2
func findMedianSortedArrays(nums1 []int, nums2 []int) (float64, error) {
	length := len(nums1) + len(nums2)
	index1 := 0
	index2 := 0
	num := 0
	for i := 0; i < length/2; i++ {
		if index1 == len(nums1) {
			num = nums2[index2]
			index2++
			continue
		}
		if index2 == len(nums2) {
			num = nums1[index1]
			index1++
			continue
		}
		if nums1[index1] < nums2[index2] {
			num = nums1[index1]
			index1++
		} else {
			num = nums2[index2]
			index2++
		}
	}
	var num2 int
	if index2 == len(nums2) {
		num2 = nums1[index1]
	} else if index1 == len(nums1) || nums1[index1] > nums2[index2] {
		num2 = nums2[index2]
	} else {
		num2 = nums1[index1]
	}
	if length%2 == 0 {

		return float64(float64(num2+num) / 2.0), nil
	} else {
		return float64(num2), nil
	}
}

// 判断数组中是否包含指定元素
func ContainByInt(arr []int, element int) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

// 判断是否是切片或者数组
func IsSliceOrArray(x interface{}) bool {
	kind := reflect.TypeOf(x).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func Filter[T any](slice []T, f func(T) bool) []T {
	filtered := make([]T, 0)
	for _, v := range slice {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Find[T any](slice []T, f func(T) bool) (*T, error) {
	for _, v := range slice {
		if f(v) {
			return &v, nil
		}
	}
	return nil, nil
}
