package squirrel

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lann/builder"
)

type deleteData struct {
	PlaceholderFormat PlaceholderFormat
	Prefixes          exprs
	From              string
	WhereParts        []Sqlizer
	OrderBys          []string
	GroupBys          []string
	HavingParts       []Sqlizer
	Limit             string
	Offset            string
	Suffixes          exprs
}

func (d *deleteData) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(d.From) == 0 {
		err = fmt.Errorf("delete statements must specify a From table")
		return
	}

	sql := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, _ = d.Prefixes.AppendToSql(sql, " ", args)
		sql.WriteString(" ")
	}

	sql.WriteString("DELETE FROM ")
	sql.WriteString(d.From)

	if len(d.WhereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(d.WhereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.OrderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(d.OrderBys, ", "))
	}

	if len(d.Limit) > 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(d.Limit)
	}

	if len(d.Offset) > 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(d.Offset)
	}

	if len(d.Suffixes) > 0 {
		sql.WriteString(" ")
		args, _ = d.Suffixes.AppendToSql(sql, " ", args)
	}

	sqlStr, err = d.PlaceholderFormat.ReplacePlaceholders(sql.String())
	return
}

// Builder

// DeleteBuilder builds SQL DELETE statements.
type DeleteBuilder builder.Builder

func init() {
	builder.Register(DeleteBuilder{}, deleteData{})
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b DeleteBuilder) PlaceholderFormat(f PlaceholderFormat) WhereConditions {
	return builder.Set(b, "PlaceholderFormat", f).(DeleteBuilder)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b DeleteBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(deleteData)
	return data.ToSql()
}

// Prefix adds an expression to the beginning of the query
func (b DeleteBuilder) Prefix(sql string, args ...interface{}) DeleteCondition {
	return builder.Append(b, "Prefixes", Expr(sql, args...)).(DeleteBuilder)
}

// From sets the table to be deleted from.
func (b DeleteBuilder) From(from string) DeleteCondition {
	return builder.Set(b, "From", from).(DeleteBuilder)
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b DeleteBuilder) Where(pred interface{}, args ...interface{}) WhereConditions {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(DeleteBuilder)
}

//expr
func (b DeleteBuilder) Expr(sql string, args ...interface{}) WhereConditions {
	return builder.Append(b, "WhereParts", newWherePart(expr{sql: sql, args: args})).(DeleteBuilder)
}

//eq
func (b DeleteBuilder) Eq(column string, arg interface{}) WhereConditions {
	return b.Where(Eq{column: arg})
}

//gt
func (b DeleteBuilder) Gt(column string, arg interface{}) WhereConditions {
	return b.Where(Gt{column: arg})
}

//gtOrEq
func (b DeleteBuilder) GtOrEq(column string, arg interface{}) WhereConditions {
	return b.Where(GtOrEq{column: arg})
}

//lt
func (b DeleteBuilder) Lt(column string, arg interface{}) WhereConditions {
	return b.Where(Lt{column: arg})
}

//ltOrEq
func (b DeleteBuilder) LtOrEq(column string, arg interface{}) WhereConditions {
	return b.Where(LtOrEq{column: arg})
}

// OrderBy adds ORDER BY expressions to the query.
func (b DeleteBuilder) GroupBy(groupBys ...string) WhereConditions {
	return builder.Extend(b, "GroupBys", groupBys).(DeleteBuilder)
}

func (b DeleteBuilder) Having(pred interface{}, rest ...interface{}) WhereConditions {
	return builder.Append(b, "HavingParts", newWherePart(pred, rest...)).(DeleteBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b DeleteBuilder) OrderBy(orderBys ...string) WhereConditions {
	return builder.Extend(b, "OrderBys", orderBys).(DeleteBuilder)
}

// Limit sets a LIMIT clause on the query.
func (b DeleteBuilder) Limit(limit int) WhereConditions {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(DeleteBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b DeleteBuilder) Offset(offset int) WhereConditions {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(DeleteBuilder)
}

// Suffix adds an expression to the end of the query
func (b DeleteBuilder) Suffix(sql string, args ...interface{}) WhereConditions {
	return builder.Append(b, "Suffixes", Expr(sql, args...)).(DeleteBuilder)
}
