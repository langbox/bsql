package bsql

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// InsertBuilder builds SQL INSERT statements.
type InsertBuilder struct {
	StatementBuilderType

	returning

	prefixes exprs
	options  []string
	into     string
	columns  []string
	values   [][]interface{}
	suffixes exprs
	iselect  *SelectBuilder
}

// NewInsertBuilder creates new instance of InsertBuilder
func NewInsertBuilder(b StatementBuilderType) *InsertBuilder {
	return &InsertBuilder{StatementBuilderType: b}
}

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b *InsertBuilder) PlaceholderFormat(f PlaceholderFormat) *InsertBuilder {
	b.placeholderFormat = f
	return b
}

// ToSql builds the query into a SQL string and bound args.
func (b *InsertBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(b.into) == 0 {
		err = fmt.Errorf("insert statements must specify a table")
		return
	}
	if len(b.values) == 0 && b.iselect == nil {
		err = fmt.Errorf("insert statements must have at least one set of values or select clause")
		return
	}

	sql := &bytes.Buffer{}

	if len(b.prefixes) > 0 {
		args, _ = b.prefixes.AppendToSql(sql, " ", args)
		sql.WriteString(" ")
	}

	sql.WriteString("INSERT ")

	if len(b.options) > 0 {
		sql.WriteString(strings.Join(b.options, " "))
		sql.WriteString(" ")
	}

	sql.WriteString("INTO ")
	sql.WriteString(b.into)
	sql.WriteString(" ")

	if len(b.columns) > 0 {
		sql.WriteString("(")
		sql.WriteString(strings.Join(b.columns, ","))
		sql.WriteString(") ")
	}

	if b.iselect != nil {
		args, err = b.appendSelectToSQL(sql, args)
	} else {
		args, err = b.appendValuesToSQL(sql, args)
	}
	if err != nil {
		return
	}

	if len(b.returning) > 0 {
		args, err = b.returning.AppendToSql(sql, args)
		if err != nil {
			return
		}
	}

	if len(b.suffixes) > 0 {
		sql.WriteString(" ")
		args, _ = b.suffixes.AppendToSql(sql, " ", args)
	}

	sqlStr, err = b.placeholderFormat.ReplacePlaceholders(sql.String())
	return
}

func (b *InsertBuilder) appendValuesToSQL(w io.Writer, args []interface{}) ([]interface{}, error) {
	if len(b.values) == 0 {
		return args, errors.New("values for insert statements are not set")
	}

	io.WriteString(w, "VALUES ")

	valuesStrings := make([]string, len(b.values))
	for r, row := range b.values {
		valueStrings := make([]string, len(row))
		for v, val := range row {

			switch typedVal := val.(type) {
			case expr:
				valueStrings[v] = typedVal.sql
				args = append(args, typedVal.args...)
			case Sqlizer:
				var valSql string
				var valArgs []interface{}
				var err error

				valSql, valArgs, err = typedVal.ToSql()
				if err != nil {
					return nil, err
				}

				valueStrings[v] = valSql
				args = append(args, valArgs...)
			default:
				valueStrings[v] = "?"
				args = append(args, val)
			}
		}
		valuesStrings[r] = fmt.Sprintf("(%s)", strings.Join(valueStrings, ","))
	}

	io.WriteString(w, strings.Join(valuesStrings, ","))

	return args, nil
}

func (b *InsertBuilder) appendSelectToSQL(w io.Writer, args []interface{}) ([]interface{}, error) {
	if b.iselect == nil {
		return args, errors.New("select clause for insert statements are not set")
	}

	selectClause, sArgs, err := b.iselect.ToSql()
	if err != nil {
		return args, err
	}

	io.WriteString(w, selectClause)
	args = append(args, sArgs...)

	return args, nil
}

// Prefix adds an expression to the beginning of the query
func (b *InsertBuilder) Prefix(sql string, args ...interface{}) *InsertBuilder {
	b.prefixes = append(b.prefixes, Expr(sql, args...))
	return b
}

// Options adds keyword options before the INTO clause of the query.
func (b *InsertBuilder) Options(options ...string) *InsertBuilder {
	b.options = append(b.options, options...)
	return b
}

// Into sets the INTO clause of the query.
func (b *InsertBuilder) Into(into string) *InsertBuilder {
	b.into = into
	return b
}

// Columns adds insert columns to the query.
func (b *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	b.columns = append(b.columns, columns...)
	return b
}

// Values adds a single row's values to the query.
func (b *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	b.values = append(b.values, values)
	return b
}

// Returning adds columns to RETURNING clause of the query
//
// INSERT ... RETURNING is PostgreSQL specific extension
func (b *InsertBuilder) Returning(columns ...string) *InsertBuilder {
	b.returning.Returning(columns...)
	return b
}

// ReturningSelect adds subquery to RETURNING clause of the query
//
// INSERT ... RETURNING is PostgreSQL specific extension
func (b *InsertBuilder) ReturningSelect(from *SelectBuilder, alias string) *InsertBuilder {
	b.returning.ReturningSelect(from, alias)
	return b
}

// Suffix adds an expression to the end of the query
func (b *InsertBuilder) Suffix(sql string, args ...interface{}) *InsertBuilder {
	b.suffixes = append(b.suffixes, Expr(sql, args...))
	return b
}

// SetMap set columns and values for insert builder from a map of column name and value
// note that it will reset all previous columns and values was set if any
func (b *InsertBuilder) SetMap(clauses map[string]interface{}) *InsertBuilder {
	// TODO: replace resetting previous values with extending existing ones?
	cols := make([]string, 0, len(clauses))
	vals := make([]interface{}, 0, len(clauses))

	for col, val := range clauses {
		cols = append(cols, col)
		vals = append(vals, val)
	}

	b.columns = cols
	b.values = [][]interface{}{vals}
	return b
}

// Select set Select clause for insert query
// If Values and Select are used, then Select has higher priority
func (b *InsertBuilder) Select(sb *SelectBuilder) *InsertBuilder {
	b.iselect = sb
	return b
}
