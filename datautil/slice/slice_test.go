package sliceutil

import (
	"fmt"
	"testing"
)

func TestMoveElementToEndAndRemovePrevious(t *testing.T) {
	data := []struct {
		input        []int
		elementIndex int
		want         []int
	}{
		{
			input:        []int{1, 2, 3, 4, 5},
			elementIndex: 2,
			want:         []int{4, 5, 3},
		},
		{
			input:        []int{1, 2},
			elementIndex: 2,
			want:         []int{},
		},
		{
			input:        []int{1, 2, 3, 4, 5},
			elementIndex: -1,
			want:         []int{},
		},
	}
	for _, d := range data {
		got, err := MoveElementToEndAndRemovePrevious(d.input, d.elementIndex)
		if err != nil {
			fmt.Printf("错误:%v", err)
		}
		fmt.Printf("结果:%v", got)
		fmt.Printf("预期:%v", d.want)
	}
}

func TestFindMedianSortedArrays(t *testing.T) {
	data := []struct {
		nums1 []int
		nums2 []int
		want  float64
	}{
		// {
		// 	nums1: []int{1, 3},
		// 	nums2: []int{2},
		// 	want:  2.0,
		// },
		// {
		// 	nums1: []int{1, 2},
		// 	nums2: []int{3, 4},
		// 	want:  2.5,
		// },
		// {
		// 	nums1: []int{1, 2, 3, 4, 5, 6, 7, 8},
		// 	nums2: []int{9, 10},
		// 	want:  5.5,
		// },
		// {
		// 	nums1: []int{1, 2, 3, 4, 6, 7, 8, 10},
		// 	nums2: []int{5, 9},
		// 	want:  5.5,
		// },
		// {
		// 	nums1: []int{1},
		// 	nums2: []int{},
		// 	want:  1.0,
		// },
		// {
		// 	nums1: []int{},
		// 	nums2: []int{1},
		// 	want:  1.0,
		// },
		// {
		// 	nums1: []int{},
		// 	nums2: []int{2, 3},
		// 	want:  2.5,
		// },
		{
			nums1: []int{1},
			nums2: []int{2, 3},
			want:  2,
		},
	}
	for _, d := range data {
		got, err := findMedianSortedArrays(d.nums1, d.nums2)
		if err != nil {
			fmt.Printf("错误:%v\n", err)
		}
		if got != d.want {
			fmt.Printf("结果:%v", got)
			fmt.Printf("预期:%v\n", d.want)
		}
	}
}
