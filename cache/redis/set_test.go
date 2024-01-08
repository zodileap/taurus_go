package redis

import (
	"testing"
	"time"

	"github.com/yohobala/taurus_go/testutil/unit"

	"github.com/redis/go-redis/v9"
)

var testClientName = "test"
var testPptions *redis.Options = &redis.Options{
	Addr:     "localhost:6379",
	Username: "",
	Password: "", // no password set
	DB:       1,  // use default DB
}

type (
	stringData struct {
		Key   string
		Val   string
		Field string
		Exp   time.Duration
	}

	mstringData struct {
		Key   string
		Pairs map[string]string
		Exp   time.Duration
	}

	mstringDataOpt struct {
		Data mstringData
		Op   string
	}
)

func TestSet(t *testing.T) {
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
	testCases := []unit.TestCase[T, SetRes]{
		{
			Name:        "测试数据1",
			MockReturns: nil,
			Input: T{
				Key: "Key1",
				Items: []I{
					{
						Val: "1",
						Exp: 0,
						Op:  "add",
					},
					{
						Val: "2",
						Exp: 0,
						Op:  "add",
					},
					{
						Val: "3",
						Exp: 0,
						Op:  "add",
					},
					{
						Val: "3",
						Exp: 0,
						Op:  "del",
					},
					{
						Exp: 0,
						Op:  "get",
					},
				},
			},
			ExpectedRes: SetRes{
				Key:    "Key1",
				Value:  []string{"1", "2"},
				AddNum: 3,
				DelNum: 1,
				Oper:   SMembers,
			},
			ExpectedErr: nil,
		},
		{
			Name: "测试全部删除",
			Input: T{
				Key: "Key2",
				Items: []I{
					{
						Val: "1",
						Exp: 0,
						Op:  "add",
					},
					{
						Val: "2",
						Exp: 0,
						Op:  "add",
					},
					{
						Op: "delall",
					},
				},
			},
			ExpectedRes: SetRes{
				Key:    "Key2",
				Value:  []string{},
				AddNum: 0,
				DelNum: 0,
				Oper:   Del,
			},
			ExpectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := c.Set()
			if err != nil {
				t.Errorf(err.Error())
			}
			for _, v := range tc.Input.Items {
				switch v.Op {
				case "add":
					s.Add(tc.Input.Key, v.Exp, v.Val)
				case "del":
					s.Del(tc.Input.Key, v.Val)
				case "get":
					s.Get(tc.Input.Key)
				case "delall":
					s.DelAll(tc.Input.Key)
				}
			}
			r, e := c.Save()
			unit.ValidErr(e, tc.ExpectedErr, t)
			data := *r.GetSet(tc.Input.Key)
			unit.ValidRes(data, tc.ExpectedRes, t)
		})
	}
}

func TestSetAddR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()

	testCases := []unit.TestCase[stringData, int64]{
		{
			Name: "测试数据1",
			Input: stringData{
				Key: "set_key1",
				Val: "1",
				Exp: 0,
			},
			ExpectedRes: 1,
		},
		{
			Name: "测试数据2",
			Input: stringData{
				Key: "set_key2",
				Val: "2",
				Exp: 0,
			},
			ExpectedRes: 1,
		},
		{
			Name: "测试数据3",
			Input: stringData{
				Key: "set_key3",
				Val: "3",
				Exp: 0,
			},
			ExpectedRes: 1,
		},
		{
			Name: "测试数据4",
			Input: stringData{
				Key: "set_key4",
				Val: "4",
				Exp: 1 * time.Second,
			},
			ExpectedRes: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := c.Set()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := s.AddR(tc.Input.Key, tc.Input.Exp, tc.Input.Val)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestSetGetR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()

	testCases := []unit.TestCase[string, []string]{
		{
			Name:        "测试数据1",
			Input:       "set_key1",
			ExpectedRes: []string{"1"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := c.Set()
			if err != nil {
				t.Errorf(err.Error())
			}

			r, e := s.GetR(tc.Input)
			unit.ValidErr(e, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestSetDelR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()

	testCases := []unit.TestCase[[]string, int64]{
		{
			Name:        "测试数据1",
			Input:       []string{"set_key1", "1"},
			ExpectedRes: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := c.Set()
			if err != nil {
				t.Errorf(err.Error())
			}

			r, e := s.DelR(tc.Input[0], tc.Input[1])
			unit.ValidErr(e, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestSetDelAllR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()

	testCases := []unit.TestCase[string, int64]{
		{
			Name:        "测试数据1",
			Input:       "set_key2",
			ExpectedRes: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := c.Set()
			if err != nil {
				t.Errorf(err.Error())
			}

			r, e := s.DelAllR(tc.Input)
			unit.ValidErr(e, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}
