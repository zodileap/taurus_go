package unit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestCase[I any, ER any] struct {
	Name           string
	MockReturns    TestCaseMockReturns
	Input          I
	ExpectedResult ER
	ExpectedErr    error
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
	if err != expectedErr {
		if err == nil ||
			expectedErr == nil ||
			(err.Error() != expectedErr.Error() &&
				!strings.Contains(err.Error(), expectedErr.Error())) {
			t.Errorf("错误信息不一致，\n期望值: %v,\n实际值: %v", expectedErr, err)
		}
	}

}

func ValidExpectedResult(result interface{}, expectedResult interface{}, t *testing.T) {

	if diff := cmp.Diff(expectedResult, result); diff != "" {

		t.Errorf("结果不一致\n期望值: %v,\n实际值: %v", expectedResult, result)
	}

}
