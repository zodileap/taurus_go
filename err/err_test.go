package err

import (
	"strings"
	"testing"
)

func TestErrCode(t *testing.T) {
	errCode := New("1234567890", "hello %s", "reason").Sprintf("world")

	if errCode.Code() != "1234567890" {
		t.Fatalf("Code 返回值不正确: %s", errCode.Code())
	}

	message := errCode.Error()
	if !strings.Contains(message, "code:1234567890") {
		t.Fatalf("错误信息缺少 code: %s", message)
	}
	if !strings.Contains(message, "hello world") {
		t.Fatalf("错误信息缺少格式化后的消息: %s", message)
	}
	if !strings.Contains(message, "reason:reason") {
		t.Fatalf("错误信息缺少 reason: %s", message)
	}
}

func TestValidFormat(t *testing.T) {
	if !ValidFormat("Err_0100020015") {
		t.Fatal("合法错误码格式被误判为非法")
	}
	if !ValidFormat("Err_030001000x") {
		t.Fatal("带 x 的合法错误码格式被误判为非法")
	}
	if ValidFormat("Err_9100020015") {
		t.Fatal("非法错误码格式被误判为合法")
	}
}
