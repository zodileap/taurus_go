package redis

import (
	"testing"
	"time"

	"github.com/yohobala/taurus_go/testutil/unit"
	"github.com/yohobala/taurus_go/tlog"
)

func TestHash(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	type I struct {
		Filed  string
		Val    string
		Exp    time.Duration
		Op     string
		Pairs  map[string]string
		Fields []string
	}
	type T struct {
		Key   string
		Items []I
	}
	testCases := []unit.TestCase[T, HashRes]{
		{
			Name: "测试Add",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Filed: "field1",
						Val:   "1",
						Exp:   0,
						Op:    HSet,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:      "han_key1",
				Value:    []string{},
				MapValue: map[string]string{},
				AddNum:   1,
				DelNum:   0,
				Oper:     HSet,
			},
		},
		{
			Name: "测试AddMR",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Pairs: map[string]string{
							"field2": "2",
							"field3": "3",
						},
						Op: HMSet,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:      "han_key1",
				Value:    []string{},
				MapValue: map[string]string{},
				AddNum:   2,
				DelNum:   0,
				Oper:     HSet,
			},
		},
		{
			Name: "测试Get",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Filed: "field1",
						Op:    HGet,
					},
				},
			},
			ExpectedRes: HashRes{
				Key: "han_key1",
				Value: []string{
					"1",
				},
				MapValue: map[string]string{},
				AddNum:   0,
				DelNum:   0,
				Oper:     HGet,
			},
		},
		{
			Name: "测试GetMR",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Fields: []string{"field1", "field2", "field3", "field4"},
						Op:     HMGet,
					},
				},
			},
			ExpectedRes: HashRes{
				Key: "han_key1",
				Value: []string{
					"1", "2", "3", "",
				},
				MapValue: map[string]string{},
				AddNum:   0,
				DelNum:   0,
				Oper:     HMGet,
			},
		},
		{
			Name: "测试GetVals",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Op: HVals,
					},
				},
			},
			ExpectedRes: HashRes{
				Key: "han_key1",
				Value: []string{
					"1", "2", "3",
				},
				MapValue: map[string]string{},
				AddNum:   0,
				DelNum:   0,
				Oper:     HVals,
			},
		},
		{
			Name: "测试GetAll",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Op: HGetAll,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:   "han_key1",
				Value: []string{},
				MapValue: map[string]string{
					"field1": "1",
					"field2": "2",
					"field3": "3",
				},
				AddNum: 0,
				DelNum: 0,
				Oper:   HGetAll,
			},
		},
		{
			Name: "测试Del",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Filed: "field1",
						Op:    HDel,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:      "han_key1",
				Value:    []string{},
				MapValue: map[string]string{},
				AddNum:   0,
				DelNum:   1,
				Oper:     HDel,
			},
		},
		{
			Name: "测试DelAll",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Filed: "field1",
						Op:    Del,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:      "han_key1",
				Value:    []string{},
				MapValue: map[string]string{},
				AddNum:   0,
				DelNum:   1,
				Oper:     Del,
			},
		},
		{
			Name: "测试最终删除",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Filed: "field2",
						Val:   "2",
						Exp:   0,
						Op:    HSet,
					},
					{
						Filed: "field3",
						Val:   "3",
						Exp:   0,
						Op:    HSet,
					},
					{
						Op: Del,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:      "han_key1",
				Value:    []string{},
				MapValue: map[string]string{},
				AddNum:   0,
				DelNum:   0, // 因为提前删除了，所以是0
				Oper:     Del,
			},
		},
		{
			Name: "测试多次添加后得到",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Filed: "field2",
						Val:   "2",
						Exp:   0,
						Op:    HSet,
					},
					{
						Filed: "field3",
						Val:   "3",
						Exp:   0,
						Op:    HSet,
					},
					{
						Op: HGetAll,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:   "han_key1",
				Value: []string{},
				MapValue: map[string]string{
					"field2": "2",
					"field3": "3",
				},
				AddNum: 2,
				DelNum: 0,
				Oper:   HGetAll,
			},
		}, {
			Name: "测试多个得到",
			Input: T{
				Key: "han_key1",
				Items: []I{
					{
						Filed: "field2",
						Op:    HGet,
					},
					{
						Pairs: map[string]string{
							"field2": "2",
							"field3": "3",
						},
						Op: HMSet,
					},
					{
						Op: HGetAll,
					},
				},
			},
			ExpectedRes: HashRes{
				Key:   "han_key1",
				Value: []string{},
				MapValue: map[string]string{
					"field2": "2",
					"field3": "3",
				},
				AddNum: 0,
				DelNum: 0,
				Oper:   HGetAll,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()

			if err != nil {
				tlog.Print(err.Error())
			}
			for _, item := range tc.Input.Items {
				switch item.Op {
				case HSet:
					h.Add(tc.Input.Key, item.Exp, item.Filed, item.Val)
				case HMSet:
					h.AddM(tc.Input.Key, item.Exp, item.Pairs)
				case HGet:
					h.Get(tc.Input.Key, item.Filed)
				case HMGet:
					h.GetM(tc.Input.Key, item.Fields)
				case HVals:
					h.GetVals(tc.Input.Key)
				case HGetAll:
					h.GetAll(tc.Input.Key)
				case HDel:
					h.Del(tc.Input.Key, item.Filed)
				case Del:
					h.DelAll(tc.Input.Key)
				}
			}
			r, err := c.Save()
			unit.ValidErr(err, tc.ExpectedErr, t)
			data := *r.GetHash(tc.Input.Key)
			unit.ValidRes(data, tc.ExpectedRes, t)
		})
	}
}

func TestHashAddR(t *testing.T) {
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
				Key:   "han_key1",
				Field: "field1",
				Val:   "1",
				Exp:   0,
			},
			ExpectedRes: 1,
		},
		{
			Name: "测试数据2",
			Input: stringData{
				Key:   "han_key1",
				Field: "field2",
				Val:   "2",
				Exp:   0,
			},
			ExpectedRes: 1,
		},
		{
			Name: "测试数据3",
			Input: stringData{
				Key:   "han_key1",
				Field: "field3",
				Val:   "3",
				Exp:   0,
			},
			ExpectedRes: 1,
		},
		{
			Name: "测试数据4",
			Input: stringData{
				Key:   "han_key1",
				Field: "field4",
				Val:   "4",
				Exp:   0,
			},
			ExpectedRes: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.AddR(tc.Input.Key, tc.Input.Exp, tc.Input.Field, tc.Input.Val)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestHashAddMR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	testCases := []unit.TestCase[mstringData, int64]{
		{
			Name: "测试数据1",
			Input: mstringData{
				Key: "han_key1",
				Pairs: map[string]string{
					"field5": "1",
					"field6": "2",
					"field7": "3",
					"field8": "4",
				},
			},
			ExpectedRes: 4,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			l, err := h.AddMR(tc.Input.Key, tc.Input.Exp, tc.Input.Pairs)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(l, tc.ExpectedRes, t)
		})
	}
}

func TestHashGetR(t *testing.T) {
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
				Key:   "han_key1",
				Field: "field1",
			},
			ExpectedRes: "1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.GetR(tc.Input.Key, tc.Input.Field)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestHashGetMR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	type T struct {
		Key   string
		Field []string
	}
	testCases := []unit.TestCase[T, []interface{}]{
		{
			Name: "测试数据1",
			Input: T{
				Key:   "han_key1",
				Field: []string{"field1", "field2", "field3", "field4"},
			},
			ExpectedRes: []interface{}{"1", "2", "3", "4"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.GetMR(tc.Input.Key, tc.Input.Field)
			tlog.Print(r)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestHashGetValsR(t *testing.T) {
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
			Input:       "han_key1",
			ExpectedRes: []string{"1", "2", "3", "4", "1", "2", "3", "4"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.GetValsR(tc.Input)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestHashGetAllR(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	c, err := GetClient(testClientName)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer c.Close()
	testCases := []unit.TestCase[string, map[string]string]{
		{
			Name:  "测试数据1",
			Input: "han_key1",
			ExpectedRes: map[string]string{
				"field1": "1",
				"field2": "2",
				"field3": "3",
				"field4": "4",
				"field5": "1",
				"field6": "2",
				"field7": "3",
				"field8": "4",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.GetAllR(tc.Input)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestHashDelR(t *testing.T) {
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
				Key:   "han_key1",
				Field: "field1",
				Val:   "1",
				Exp:   0,
			},
			ExpectedRes: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.DelR(tc.Input.Key, tc.Input.Field)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}

func TestHashDelAllR(t *testing.T) {
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
			Input:       "han_key1",
			ExpectedRes: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			h, err := c.Hash()
			if err != nil {
				t.Errorf(err.Error())
			}
			r, err := h.DelAllR(tc.Input)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidRes(r, tc.ExpectedRes, t)
		})
	}
}
