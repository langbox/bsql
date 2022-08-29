package bsql

import (
	"context"
	"database/sql"
)

type DBStub struct {
	res sql.Result
	err error

	LastPrepareSql string
	PrepareCount   int

	LastExecSql  string
	LastExecArgs []interface{}

	LastQuerySql  string
	LastQueryArgs []interface{}

	LastQueryRowSql  string
	LastQueryRowArgs []interface{}
}

func (s *DBStub) Prepare(query string) (*sql.Stmt, error) {
	s.LastPrepareSql = query
	s.PrepareCount++
	return nil, nil
}

func (s *DBStub) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	s.LastPrepareSql = query
	s.PrepareCount++
	return nil, nil
}

func (s *DBStub) Exec(query string, args ...interface{}) (sql.Result, error) {
	s.LastExecSql = query
	s.LastExecArgs = args
	return s.res, s.err
}

func (s *DBStub) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	s.LastExecSql = query
	s.LastExecArgs = args
	return s.res, s.err
}

func (s *DBStub) Query(query string, args ...interface{}) (*sql.Rows, error) {
	s.LastQuerySql = query
	s.LastQueryArgs = args
	return nil, nil
}

func (s *DBStub) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	s.LastQuerySql = query
	s.LastQueryArgs = args
	return nil, nil
}

type resultStub struct {
	rowsAffected int64
	lastInsertId int64
	err          error
}

func (r *resultStub) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

func (r *resultStub) LastInsertId() (int64, error) {
	return r.lastInsertId, r.err
}

var sqlizer = Select("test")
var sqlStr = "SELECT test"
