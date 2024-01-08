package asset

import (
	"fmt"
	"testing"

	"github.com/yohobala/taurus_go/testutil/unit"
)

// TestFile 测试文件
func TestFile(t *testing.T) {
	type _TestFile struct {
		Path string
		Body []byte
	}

	tcs := []unit.TestCase[_TestFile, any]{
		{
			Name: "测试文件1",
			Input: _TestFile{
				Path: "test.txt",
				Body: []byte("Hellow, World!"),
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			var assets Assets
			assets.Add(tc.Input.Path, tc.Input.Body)
			if err := assets.Write(); err != nil {
				unit.ValidErr(err, tc.ExpectedErr, t)
			}
		})
	}
}

func TestDir(t *testing.T) {
	tcs := []unit.TestCase[string, any]{
		{
			Name:  "测试文件夹1",
			Input: "test",
		},
		{
			Name:        "测试文件夹2",
			Input:       "",
			ExpectedErr: fmt.Errorf("code: 02-0001-0001"),
		},
		{
			Name:  "测试文件夹3",
			Input: "./-*test.go",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			var assets Assets
			assets.AddDir(tc.Input)
			if err := assets.Write(); err != nil {
				unit.ValidErr(err, tc.ExpectedErr, t)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type _TestFile struct {
		Path string
		Body []byte
	}

	tcs := []unit.TestCase[_TestFile, any]{
		{
			Name: "测试文件1",
			Input: _TestFile{
				Path: "test.go",
				Body: []byte(`
				package main
				func Main() {fmt.Println("Hello, World!")}
				`),
			},
		},
		{
			Name: "测试文件2",
			Input: _TestFile{
				Path: "test.go",
				Body: []byte(`
				func Main() {fmt.Println("Hello, World!")}
				`),
			},
			ExpectedErr: fmt.Errorf("code: 02-0001-0003"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			var assets Assets
			assets.Add(tc.Input.Path, tc.Input.Body)
			if err := assets.Write(); err != nil {
				unit.ValidErr(err, tc.ExpectedErr, t)
			}
			if err := assets.Format(); err != nil {
				unit.ValidErr(err, tc.ExpectedErr, t)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	type _CopyFile struct {
		Src string
		Dst string
	}

	tcs := []unit.TestCase[_CopyFile, any]{
		{
			Name: "测试文件1",
			Input: _CopyFile{
				Src: "test.go",
				Dst: "test2.go",
			},
		},
		{
			Name: "测试文件2",
			Input: _CopyFile{
				Src: "test.go",
				Dst: "./test2.go",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			if err := CopyFile(tc.Input.Src, tc.Input.Dst); err != nil {
				unit.ValidErr(err, tc.ExpectedErr, t)
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	type _CopyDir struct {
		Src string
		Dst string
	}

	tcs := []unit.TestCase[_CopyDir, any]{
		{
			Name: "测试文件夹1",
			Input: _CopyDir{
				Src: "test",
				Dst: "test2",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			if err := CopyDir(tc.Input.Src, tc.Input.Dst); err != nil {
				unit.ValidErr(err, tc.ExpectedErr, t)
			}
		})
	}
}
