package sql

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/zodileap/taurus_go/entity/dialect"
)

type fakeExecQuerier struct {
	execCalled  bool
	queryCalled bool
	execErr     error
	queryErr    error
}

func (f *fakeExecQuerier) ExecContext(_ context.Context, _ string, _ ...any) (sql.Result, error) {
	f.execCalled = true
	return nil, f.execErr
}

func (f *fakeExecQuerier) QueryContext(_ context.Context, _ string, _ ...any) (*sql.Rows, error) {
	f.queryCalled = true
	return nil, f.queryErr
}

func TestNewDriverAndDialect(t *testing.T) {
	driver := NewDriver(dialect.PostgreSQL, Conn{ExecQuerier: &fakeExecQuerier{}})
	if driver.Dialect() != dialect.PostgreSQL {
		t.Fatalf("Dialect 返回值不正确: %s", driver.Dialect())
	}
}

func TestConnExec(t *testing.T) {
	execQuerier := &fakeExecQuerier{}
	conn := Conn{ExecQuerier: execQuerier}

	if err := conn.Exec(context.Background(), "SELECT 1", nil, nil); err != nil {
		t.Fatalf("Exec(nil result) 失败: %v", err)
	}
	if !execQuerier.execCalled {
		t.Fatal("Exec 未调用底层 ExecContext")
	}

	if err := conn.Exec(context.Background(), "SELECT 1", nil, new(int)); err == nil {
		t.Fatal("非法结果类型应返回错误")
	}
}

func TestConnQueryPassThroughError(t *testing.T) {
	expectedErr := errors.New("query failed")
	execQuerier := &fakeExecQuerier{queryErr: expectedErr}
	conn := Conn{ExecQuerier: execQuerier}

	var rows dialect.Rows
	err := conn.Query(context.Background(), "SELECT 1", nil, &rows)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Query 未透传底层错误: %v", err)
	}
	if !execQuerier.queryCalled {
		t.Fatal("Query 未调用底层 QueryContext")
	}
}
