package bsql

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type wherePart part

func newWherePart(pred interface{}, args ...interface{}) Sqlizer {
	return &wherePart{pred: pred, args: args}
}

func (p wherePart) ToSql() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case Sqlizer:
		return pred.ToSql()
	case map[string]interface{}:
		return Eq(pred).ToSql()
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string-keyed map or string, not %T", pred)
	}
	return
}

// WhereBuilder builds SQL WHERE statements.
type WhereBuilder struct {
	StatementBuilderType

	whereParts  []Sqlizer
	groupBys    []string
	havingParts []Sqlizer
	orderBys    []string

	limit       uint64
	limitValid  bool
	offset      uint64
	offsetValid bool
}

// ToSql builds the query into a SQL string and bound args.
func (b *WhereBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}

	if len(b.whereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(b.whereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(b.groupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(b.groupBys, ", "))
	}

	if len(b.havingParts) > 0 {
		sql.WriteString(" HAVING ")
		args, err = appendToSql(b.havingParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(b.orderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(b.orderBys, ", "))
	}

	// TODO: limit == 0 and offswt == 0 are valid. Need to go dbr way and implement offsetValid and limitValid
	if b.limitValid {
		sql.WriteString(" LIMIT ")
		sql.WriteString(strconv.FormatUint(b.limit, 10))
	}

	if b.offsetValid {
		sql.WriteString(" OFFSET ")
		sql.WriteString(strconv.FormatUint(b.offset, 10))
	}

	sqlStr, err = b.placeholderFormat.ReplacePlaceholders(sql.String())
	return

}

// Where will panic if pred isn't any of the above types.
func (b *WhereBuilder) Where(pred interface{}, args ...interface{}) *WhereBuilder {
	b.whereParts = append(b.whereParts, newWherePart(pred, args...))
	return b
}

// GroupBy adds GROUP BY expressions to the query.
func (b *WhereBuilder) GroupBy(groupBys ...string) *WhereBuilder {
	b.groupBys = append(b.groupBys, groupBys...)
	return b
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b *WhereBuilder) Having(pred interface{}, rest ...interface{}) *WhereBuilder {
	b.havingParts = append(b.havingParts, newWherePart(pred, rest...))
	return b
}

// OrderBy adds ORDER BY expressions to the query.
func (b *WhereBuilder) OrderBy(orderBys ...string) *WhereBuilder {
	b.orderBys = append(b.orderBys, orderBys...)
	return b
}

// Limit sets a LIMIT clause on the query.
func (b *WhereBuilder) Limit(limit uint64) *WhereBuilder {
	b.limit = limit
	b.limitValid = true
	return b
}

// Offset sets a OFFSET clause on the query.
func (b *WhereBuilder) Offset(offset uint64) *WhereBuilder {
	b.offset = offset
	b.offsetValid = true
	return b
}
