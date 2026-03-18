package entitysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/dialect"
)

type createTestTx struct {
	dialect     dialect.DbDriver
	execCalled  bool
	queryCalled bool
}

func (t *createTestTx) Commit() error {
	return nil
}

func (t *createTestTx) Rollback() error {
	return nil
}

func (t *createTestTx) Dialect() dialect.DbDriver {
	return t.dialect
}

func (t *createTestTx) Exec(ctx context.Context, query string, args []any, v any) error {
	t.execCalled = true
	return nil
}

func (t *createTestTx) Query(ctx context.Context, query string, args []any, v *dialect.Rows) error {
	t.queryCalled = true
	v.RowsScanner = &createTestRows{remaining: 1}
	return nil
}

type createTestRows struct {
	remaining int
}

func (r *createTestRows) Close() error {
	return nil
}

func (r *createTestRows) ColumnTypes() ([]*sql.ColumnType, error) {
	return nil, nil
}

func (r *createTestRows) Columns() ([]string, error) {
	return nil, nil
}

func (r *createTestRows) Err() error {
	return nil
}

func (r *createTestRows) Next() bool {
	if r.remaining == 0 {
		return false
	}
	r.remaining--
	return true
}

func (r *createTestRows) NextResultSet() bool {
	return false
}

func (r *createTestRows) Scan(dest ...any) error {
	return nil
}

func TestNewCreateReturnsErrorForUnsupportedReturningDialect(t *testing.T) {
	tx := &createTestTx{dialect: dialect.MySQL}

	err := NewCreate(context.Background(), tx, newCreateSpecWithReturning(func(row dialect.Rows, selects []ScannerField) error {
		return nil
	}))
	if !errors.Is(err, ErrReturningUnsupported) {
		t.Fatalf("期望 ErrReturningUnsupported，实际 %v", err)
	}
	if tx.execCalled || tx.queryCalled {
		t.Fatalf("返回不支持错误时不应执行 SQL，exec=%v query=%v", tx.execCalled, tx.queryCalled)
	}
}

func TestNewCreateWithReturningQueriesAndScans(t *testing.T) {
	tx := &createTestTx{dialect: dialect.PostgreSQL}
	scanCalled := false

	err := NewCreate(context.Background(), tx, newCreateSpecWithReturning(func(row dialect.Rows, selects []ScannerField) error {
		scanCalled = true
		if len(selects) != 1 || selects[0].String() != "id" {
			t.Fatalf("返回字段不正确: %#v", selects)
		}
		return nil
	}))
	if err != nil {
		t.Fatalf("NewCreate 返回了意外错误: %v", err)
	}
	if !tx.queryCalled {
		t.Fatal("RETURNING 查询应走 Query")
	}
	if tx.execCalled {
		t.Fatal("RETURNING 查询不应走 Exec")
	}
	if !scanCalled {
		t.Fatal("RETURNING 查询应调用 Scan")
	}
}

func newCreateSpecWithReturning(scan Scanner) *CreateSpec {
	idField := NewFieldSpec("id")
	idField.Param = entity.FieldValue(1)
	idField.ParamFormat = func(dbType dialect.DbDriver, param string) string {
		return param
	}

	return &CreateSpec{
		Entity: &EntitySpec{
			Name:    "users",
			Columns: NewFieldSpecs("id"),
		},
		Fields:    [][]*FieldSpec{{&idField}},
		Returning: []FieldName{"id"},
		Scan:      scan,
	}
}
