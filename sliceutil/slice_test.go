package sliceutil

import (
	"reflect"
	"testing"
)

func TestMoveElementToEndAndRemovePrevious(t *testing.T) {
	data := []struct {
		name         string
		input        []int
		elementIndex int
		want         []int
		wantErr      bool
	}{
		{
			name:         "move tail after cut",
			input:        []int{1, 2, 3, 4, 5},
			elementIndex: 2,
			want:         []int{4, 5, 3},
		},
		{
			name:         "index out of range",
			input:        []int{1, 2},
			elementIndex: 2,
			wantErr:      true,
		},
		{
			name:         "negative index",
			input:        []int{1, 2, 3, 4, 5},
			elementIndex: -1,
			wantErr:      true,
		},
	}
	for _, d := range data {
		got, err := MoveElementToEndAndRemovePrevious(d.input, d.elementIndex)
		if d.wantErr {
			if err == nil {
				t.Fatalf("%s: 期望返回错误，实际成功", d.name)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: 返回了意外错误: %v", d.name, err)
		}
		if !reflect.DeepEqual(got, d.want) {
			t.Fatalf("%s: 结果不正确，期望 %v，实际 %v", d.name, d.want, got)
		}
	}
}

func TestContainByInt(t *testing.T) {
	if !ContainByInt([]int{1, 2, 3}, 2) {
		t.Fatal("期望切片包含元素 2")
	}
	if ContainByInt([]int{1, 2, 3}, 4) {
		t.Fatal("期望切片不包含元素 4")
	}
}

func TestIsSliceOrArray(t *testing.T) {
	if !IsSliceOrArray([]int{1, 2, 3}) {
		t.Fatal("切片应被识别为 slice/array")
	}
	if !IsSliceOrArray([2]int{1, 2}) {
		t.Fatal("数组应被识别为 slice/array")
	}
	if IsSliceOrArray(123) {
		t.Fatal("普通整型不应被识别为 slice/array")
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool {
		return n%2 == 0
	})
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter 结果不正确，期望 %v，实际 %v", want, got)
	}
}

func TestFind(t *testing.T) {
	got, err := Find([]int{1, 2, 3, 4, 5}, func(n int) bool {
		return n > 3
	})
	if err != nil {
		t.Fatalf("Find 返回了意外错误: %v", err)
	}
	if got == nil || *got != 4 {
		t.Fatalf("Find 结果不正确，实际 %v", got)
	}

	notFound, err := Find([]int{1, 2, 3}, func(n int) bool {
		return n > 10
	})
	if err != nil {
		t.Fatalf("Find 未命中时返回了意外错误: %v", err)
	}
	if notFound != nil {
		t.Fatalf("未命中时应返回 nil，实际 %v", *notFound)
	}
}
