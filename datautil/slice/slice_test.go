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
