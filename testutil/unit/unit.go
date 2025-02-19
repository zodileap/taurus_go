package unit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	terr "github.com/yohobala/taurus_go/err"
)

type TestCase[I any, ER any] struct {
	Name        string
	MockReturns TestCaseMockReturns
	Input       I
	ExpectedRes ER
	ExpectedErr error
}

type TestCaseMockReturns = map[string]TestCaseMockReturn

type TestCaseMockReturn struct {
	Result interface{}
	Err    error
}

type TestAPIResult struct {
	StatusCode int
	Body       interface{}
}

func ValidErr(err error, expectedErr error, t *testing.T) {
	te, ok := err.(*terr.ErrCode)
	te2, ok2 := expectedErr.(*terr.ErrCode)
	if ok && ok2 {
		if te.Code() != te2.Code() {
			t.Errorf("错误信息不一致，\n期望值: %v,\n实际值: %v", expectedErr, err)
		}
		return
	}

	if err != expectedErr {
		if err == nil ||
			expectedErr == nil ||
			(err.Error() != expectedErr.Error() &&
				!strings.Contains(err.Error(), expectedErr.Error())) {
			t.Errorf("错误信息不一致，\n期望值: %v,\n实际值: %v", expectedErr, err)
		}
	}

}

func ValidRes(result interface{}, expectedResult interface{}, t *testing.T) {

	if diff := cmp.Diff(expectedResult, result); diff != "" {

		t.Errorf("结果不一致\n期望值: %v,\n实际值: %v", expectedResult, result)
	}
}

// Must 是一个工具函数，用于检查错误，如果有错误，会直接输出错误信息。
func Must(t *testing.T, err error) {
	if err != nil {
		t.Errorf(err.Error())
	}
}
