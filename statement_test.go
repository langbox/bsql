package bsql

import (
	"fmt"
	"testing"
)

func TestStatementBuilder(t *testing.T) {
	// db := &DBStub{}
	// sb := StatementBuilder.RunWith(db)

	// sb.Select("test").Exec()
	// assert.Equal(t, "SELECT test", db.LastExecSql)

	//
	query, args, err := Select("*").From("test").ToSql()
	fmt.Printf("query:%s, args:%v, err:%v", query, args, err)
}

func TestStatementBuilderPlaceholderFormat(t *testing.T) {
	// db := &DBStub{}
	// sb := StatementBuilder.RunWith(db).PlaceholderFormat(Dollar)

	// sb.Select("test").Where("x = ?").Exec()
	// assert.Equal(t, "SELECT test WHERE x = $1", db.LastExecSql)
}

func TestRunWithDB(t *testing.T) {
	// db := &sql.DB{}
	// assert.NotPanics(t, func() {
	// 	Select().RunWith(db)
	// 	Insert("t").RunWith(db)
	// 	Update("t").RunWith(db)
	// 	Delete("t").RunWith(db)
	// }, "RunWith(*sql.DB) should not panic")

}

func TestRunWithTx(t *testing.T) {
	// tx := &sql.Tx{}
	// assert.NotPanics(t, func() {
	// 	Select().RunWith(tx)
	// 	Insert("t").RunWith(tx)
	// 	Update("t").RunWith(tx)
	// 	Delete("t").RunWith(tx)
	// }, "RunWith(*sql.Tx) should not panic")
}
