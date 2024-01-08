package redis

import (
	"errors"
	"testing"
	"time"

	"github.com/yohobala/taurus_go/testutil/unit"

	"github.com/redis/go-redis/v9"
)

var testClientName = "test"
var testPptions *redis.Options = &redis.Options{
	Addr:     "localhost:30001",
	Username: "root",
	Password: "root", // no password set
	DB:       15,     // use default DB
}

type StringData struct {
	Key        string
	Value      string
	Expiration time.Duration
}

func TestSet(t *testing.T) {
	if err := StringSet(testClientName, "", "", 0); err == nil {
		t.Errorf("未初始化时Set()应该返回错误")
	}

	SetClient(testClientName, testPptions)
	defer ClearClient()

	testCases := []unit.TestCase[StringData, any]{
		{
			Name:        "测试数据1",
			MockReturns: nil,
			Input: StringData{
				Key:        "",
				Value:      "2131231",
				Expiration: 1,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:        "测试数据2",
			MockReturns: nil,
			Input: StringData{
				Key:        "a",
				Value:      "213121",
				Expiration: 0,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:        "测试数据3",
			MockReturns: nil,
			Input: StringData{
				Key:        "detartrated",
				Value:      "231231",
				Expiration: 0,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:        "测试数据4",
			MockReturns: nil,
			Input: StringData{
				Key:        "Evil I did dwell; lewd did : I live.",
				Value:      "231232131",
				Expiration: 0,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:        "测试数据5",
			MockReturns: nil,
			Input: StringData{
				Key:        "Able was I ere I saw Elba",
				Value:      "54535141",
				Expiration: 0,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:        "测试数据6",
			MockReturns: nil,
			Input: StringData{
				Key:        "été",
				Value:      "132",
				Expiration: 0,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := StringSet(testClientName, tc.Input.Key, tc.Input.Value, tc.Input.Expiration)
			unit.ValidErr(err, tc.ExpectedErr, t)
		})
	}
}

func TestGet(t *testing.T) {
	if _, err := StringGet(testClientName, ""); err == nil {
		t.Errorf("未初始化时Set()应该返回错误")
	}

	SetClient(testClientName, testPptions)
	defer ClearClient()

	testCases := []unit.TestCase[string, string]{
		{
			Name:           "测试数据1",
			MockReturns:    nil,
			Input:          "",
			ExpectedResult: "2131231",
			ExpectedErr:    nil,
		},
		{
			Name:           "测试数据2",
			MockReturns:    nil,
			Input:          "a",
			ExpectedResult: "213121",
			ExpectedErr:    nil,
		},
		{
			Name:           "测试数据3",
			MockReturns:    nil,
			Input:          "detartrated",
			ExpectedResult: "231231",
			ExpectedErr:    nil,
		},
		{
			Name:           "测试数据4",
			MockReturns:    nil,
			Input:          "Evil I did dwell; lewd did I live.",
			ExpectedResult: "231232131",
			ExpectedErr:    nil,
		},
		{
			Name:           "测试数据5",
			MockReturns:    nil,
			Input:          "Able was I ere I saw Elba",
			ExpectedResult: "54535141",
			ExpectedErr:    nil,
		},
		{
			Name:           "测试数据6",
			MockReturns:    nil,
			Input:          "été",
			ExpectedResult: "132",
			ExpectedErr:    nil,
		},
		{
			Name:           "测试数据7",
			MockReturns:    nil,
			Input:          "nil",
			ExpectedResult: "",
			ExpectedErr:    errors.New("redis: nil"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := StringGet(testClientName, tc.Input)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidExpectedResult(res, tc.ExpectedResult, t)
		})
	}
}

func TestHashFieldSet(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	testCases := []unit.TestCase[struct {
		ClientName string
		Key        string
		Field      string
		Value      string
		Exp        time.Duration
	}, any]{
		{
			Name:        "测试有效的哈希字段设置",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
				Field      string
				Value      string
				Exp        time.Duration
			}{
				ClientName: testClientName,
				Key:        "hashKey",
				Field:      "field1",
				Value:      "value1",
				Exp:        10,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:        "测试无效的客户端名称",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
				Field      string
				Value      string
				Exp        time.Duration
			}{
				ClientName: "invalidClient",
				Key:        "hashKey",
				Field:      "field1",
				Value:      "value1",
				Exp:        10,
			},
			ExpectedResult: nil,
			ExpectedErr:    Err_nil_options,
		},
		// 可以添加更多测试用例
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := HashFieldSet(tc.Input.ClientName, tc.Input.Key, tc.Input.Field, tc.Input.Value, tc.Input.Exp)
			unit.ValidErr(err, tc.ExpectedErr, t)
		})
	}
}

func TestHashSet(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	testCases := []unit.TestCase[struct {
		ClientName string
		Key        string
		Fields     map[string]interface{}
		Exp        time.Duration
	}, any]{
		{
			Name:        "测试有效的哈希设置",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
				Fields     map[string]interface{}
				Exp        time.Duration
			}{
				ClientName: testClientName,
				Key:        "hashKey",
				Fields: map[string]interface{}{
					"field1": "value1",
					"field2": "value2",
				},
				Exp: 10,
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:        "测试无效的客户端名称",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
				Fields     map[string]interface{}
				Exp        time.Duration
			}{
				ClientName: "invalidClient",
				Key:        "hashKey",
				Fields: map[string]interface{}{
					"field1": "value1",
					"field2": "value2",
				},
				Exp: 10,
			},
			ExpectedResult: nil,
			ExpectedErr:    Err_nil_options,
		},
		// 可以添加更多测试用例
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := HashSet(tc.Input.ClientName, tc.Input.Key, tc.Input.Fields, tc.Input.Exp)
			unit.ValidErr(err, tc.ExpectedErr, t)
		})
	}
}

func TestGetHashField(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	testCases := []unit.TestCase[struct {
		ClientName string
		Key        string
		Field      string
	}, string]{
		{
			Name:        "测试有效的哈希字段获取",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
				Field      string
			}{
				ClientName: testClientName,
				Key:        "hashKey",
				Field:      "field1",
			},
			ExpectedResult: "value1",
			ExpectedErr:    nil,
		},
		{
			Name:        "测试无效的客户端名称",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
				Field      string
			}{
				ClientName: "invalidClient",
				Key:        "hashKey",
				Field:      "field1",
			},
			ExpectedResult: "",
			ExpectedErr:    Err_nil_options,
		},
		{
			Name:        "测试不存在的key",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
				Field      string
			}{
				ClientName: testClientName,
				Key:        "notExistKey",
				Field:      "field1",
			},
			ExpectedResult: "",
			ExpectedErr:    Err_not_key,
		},
		// 可以添加更多测试用例
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := HashFieldGet(tc.Input.ClientName, tc.Input.Key, tc.Input.Field)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidExpectedResult(result, tc.ExpectedResult, t)
		})
	}
}

func TestGetHash(t *testing.T) {
	SetClient(testClientName, testPptions)
	defer ClearClient()
	testCases := []unit.TestCase[struct {
		ClientName string
		Key        string
	}, map[string]string]{
		{
			Name:        "测试有效的哈希对象获取",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
			}{
				ClientName: testClientName,
				Key:        "hashKey",
			},
			ExpectedResult: map[string]string{
				"field1": "value1",
				"field2": "value2",
			},
			ExpectedErr: nil,
		},
		{
			Name:        "测试无效的客户端名称",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
			}{
				ClientName: "invalidClient",
				Key:        "hashKey",
			},
			ExpectedResult: nil,
			ExpectedErr:    Err_nil_options,
		},
		{
			Name:        "测试不存在的key",
			MockReturns: nil,
			Input: struct {
				ClientName string
				Key        string
			}{
				ClientName: testClientName,
				Key:        "notExistKey",
			},
			ExpectedResult: nil,
			ExpectedErr:    Err_not_key,
		},
		// 可以添加更多测试用例
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := HashGet(tc.Input.ClientName, tc.Input.Key)
			unit.ValidErr(err, tc.ExpectedErr, t)
			unit.ValidExpectedResult(result, tc.ExpectedResult, t)
		})
	}
}
