package entity

import (
	"errors"
	"testing"

	"github.com/zodileap/taurus_go/entity/dialect"

	terr "github.com/zodileap/taurus_go/err"
)

func cloneConfig(c Config) Config {
	cloned := Config{}
	if c.BatchSize != nil {
		value := *c.BatchSize
		cloned.BatchSize = &value
	}
	if c.SqlConsole != nil {
		value := *c.SqlConsole
		cloned.SqlConsole = &value
	}
	if c.SqlLogger != nil {
		value := *c.SqlLogger
		cloned.SqlLogger = &value
	}
	return cloned
}

func requireEntityErrCode(t *testing.T, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("期望错误码 %s，实际无错误", want)
	}

	var errCode terr.ErrCode
	if !errors.As(got, &errCode) {
		t.Fatalf("期望 ErrCode，实际为 %T: %v", got, got)
	}
	if errCode.Code() != want {
		t.Fatalf("错误码不匹配，期望 %s，实际 %s", want, errCode.Code())
	}
}

func resetConnections() {
	mu.Lock()
	defer mu.Unlock()
	clients = make(map[string]ConnectionConfig)
}

func TestConfigDefaultsAndSetConfig(t *testing.T) {
	original := cloneConfig(GetConfig())
	t.Cleanup(func() {
		SetConfig(original)
	})

	current := GetConfig()
	if current.BatchSize == nil || *current.BatchSize <= 0 {
		t.Fatalf("默认 BatchSize 不正确: %+v", current.BatchSize)
	}
	if current.SqlConsole == nil || *current.SqlConsole {
		t.Fatalf("默认 SqlConsole 不正确: %+v", current.SqlConsole)
	}
	if current.SqlLogger == nil || *current.SqlLogger == "" {
		t.Fatalf("默认 SqlLogger 不正确: %+v", current.SqlLogger)
	}

	newBatchSize := 128
	SetConfig(Config{BatchSize: &newBatchSize})

	updated := GetConfig()
	if updated.BatchSize == nil || *updated.BatchSize != 128 {
		t.Fatalf("BatchSize 更新失败: %+v", updated.BatchSize)
	}
	if updated.SqlConsole == nil || *updated.SqlConsole != *original.SqlConsole {
		t.Fatalf("SqlConsole 不应被覆盖: %+v", updated.SqlConsole)
	}
	if updated.SqlLogger == nil || *updated.SqlLogger != *original.SqlLogger {
		t.Fatalf("SqlLogger 不应被覆盖: %+v", updated.SqlLogger)
	}
}

func TestAddConnectionValidation(t *testing.T) {
	resetConnections()
	t.Cleanup(resetConnections)

	requireEntityErrCode(t, AddConnection(ConnectionConfig{}), Err_0100010001.Code())

	requireEntityErrCode(t, AddConnection(ConnectionConfig{
		Tag:    "invalid",
		Driver: dialect.DbDriver("sqlite"),
	}), Err_0100010002.Code())

	if err := AddConnection(ConnectionConfig{
		Tag:    "primary",
		Driver: dialect.PostgreSQL,
	}); err != nil {
		t.Fatalf("添加合法连接失败: %v", err)
	}

	requireEntityErrCode(t, AddConnection(ConnectionConfig{
		Tag:    "primary",
		Driver: dialect.PostgreSQL,
	}), Err_0100010004.Code())
}

func TestGetConnectionValidation(t *testing.T) {
	resetConnections()
	t.Cleanup(resetConnections)

	_, err := GetConnection("missing")
	requireEntityErrCode(t, err, Err_0100010003.Code())

	if err := AddConnection(ConnectionConfig{
		Tag:    "postgres",
		Driver: dialect.PostgreSQL,
	}); err != nil {
		t.Fatalf("添加连接失败: %v", err)
	}

	_, err = GetConnection("postgres")
	requireEntityErrCode(t, err, Err_0100010005.Code())
}
