package redis

import (
	"testing"
	"time"

	terr "github.com/zodileap/taurus_go/err"
	"github.com/zodileap/taurus_go/tlog"

	"github.com/zodileap/taurus_go/testutil/unit"
)

func TestString(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	type I struct {
		Val string
		Exp time.Duration
		Op  string
	}
	type T struct {
		Key   string
		Items []I
	}
	testCases := []unit.TestCase[T, StringRes]{
		{
			Name:        "测试数据1",
			MockReturns: nil,
			Input: T{
				Key: "string_key1",
				Items: []I{
					{
						Val: "1",
						Exp: 1000,
						Op:  "add",
					},
					{
						Op: "get",
					},
					{
						Op: "del",
					},
				},
			},
			ExpectedRes: StringRes{},
			ExpectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := c.String()
			if err != nil {
				t.Errorf(err.Error())
			}
			for _, item := range tc.Input.Items {
				switch item.Op {
				case "add":
					s.Add(item.Val, item.Exp, tc.Input.Key)
				case "get":
					s.Get(tc.Input.Key)
				case "del":
					s.Del(tc.Input.Key)
				}
			}
			_, err = c.Save()
			unit.ValidErr(err, tc.ExpectedErr, t)
			tlog.Print(err)
			// unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestStringAtomicity(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	type I struct {
		Val string
		Exp time.Duration
		Op  string
	}
	type T struct {
		Key   string
		Items []I
	}
	testCases := []unit.TestCase[[]T, StringRes]{
		{
			Name:        "测试数据1",
			MockReturns: nil,
			Input: []T{
				{
					Key: "string_key1",
					Items: []I{
						{
							Val: "1",
							Exp: 0,
							Op:  "add",
						},
						{
							Op: "get",
						},
						{
							Op: "delete",
						},
					},
				},
				{
					Key: "string_key2",
					Items: []I{
						{
							Val: "1",
							Exp: 0,
							Op:  "add",
						},
					},
				},
			},
			ExpectedRes: StringRes{},
			ExpectedErr: terr.New("0300010003", "", ""),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := c.String()
			if err != nil {
				t.Errorf(err.Error())
			}
			for _, g := range tc.Input {
				for _, item := range g.Items {
					switch item.Op {
					case "add":
						s.Add(item.Val, item.Exp, g.Key)
					case "get":
						s.Get(g.Key)
					case "del":
						s.Del(g.Key)
					}
				}
			}
			_, err = c.Save()
			unit.ValidErr(err, tc.ExpectedErr, t)
			// unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestStringAddR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	testCases := []unit.TestCase[stringData, string]{
		{
			Name: "测试数据1",
			Input: stringData{
				Key: "string_key1",
				Val: "1",
				Exp: 0,
			},
			ExpectedRes: "OK",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.String()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = h.AddR(tc.Input.Key, tc.Input.Exp, tc.Input.Val)
			unit.ValidErr(err, tc.ExpectedErr, t)
		})
	}
}

func TestStringGetR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	testCases := []unit.TestCase[string, string]{
		{
			Name:        "测试数据1",
			Input:       "string_key1",
			ExpectedRes: "1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.String()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.GetR(tc.Input)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestStringDelR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	testCases := []unit.TestCase[string, int64]{
		{
			Name:  "测试数据1",
			Input: "string_key1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.String()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = h.DelR(tc.Input)
			unit.ValidErr(err, tc.ExpectedErr, t)
		})
	}
}
