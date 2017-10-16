package squirrel

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/lann/builder"
)

type updateData struct {
	PlaceholderFormat PlaceholderFormat
	Prefixes          exprs
	Table             string
	SetClauses        []setClause
	WhereParts        []Sqlizer
	GroupBys          []string
	HavingParts       []Sqlizer
	OrderBys          []string
	Limit             string
	Offset            string
	Suffixes          exprs
}

type setClause struct {
	column string
	value  interface{}
}

func (d *updateData) ToSql() (sqlStr string, args []interface{}, err error) {
	if len(d.Table) == 0 {
		err = fmt.Errorf("update statements must specify a table")
		return
	}
	if len(d.SetClauses) == 0 {
		err = fmt.Errorf("update statements must have at least one Set clause")
		return
	}

	sql := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, _ = d.Prefixes.AppendToSql(sql, " ", args)
		sql.WriteString(" ")
	}

	sql.WriteString("UPDATE ")
	sql.WriteString(d.Table)

	sql.WriteString(" SET ")
	setSqls := make([]string, len(d.SetClauses))
	for i, setClause := range d.SetClauses {
		var valSql string
		e, isExpr := setClause.value.(expr)
		if isExpr {
			valSql = e.sql
			args = append(args, e.args...)
		} else {
			valSql = "?"
			args = append(args, setClause.value)
		}
		if ok, column := getSetColumn(setClause.column); ok {
			setSqls[i] = fmt.Sprintf("%s%s", column, valSql)
		} else {
			setSqls[i] = fmt.Sprintf("%s = %s", column, valSql)
		}
	}
	sql.WriteString(strings.Join(setSqls, ", "))

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
func getSetColumn(column string) (bool, string) {
	r := []rune(column)
	str := string(r[:3])
	if str == "=-=" || str == "=+=" {
		return true, string(r[3:])
	}
	return false, column
}

// Builder

// UpdateBuilder builds SQL UPDATE statements.
type UpdateBuilder builder.Builder

func init() {
	builder.Register(UpdateBuilder{}, updateData{})
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b UpdateBuilder) PlaceholderFormat(f PlaceholderFormat) WhereConditions {
	return builder.Set(b, "PlaceholderFormat", f).(UpdateBuilder)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b UpdateBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(updateData)
	return data.ToSql()
}

// Prefix adds an expression to the beginning of the query
func (b UpdateBuilder) Prefix(sql string, args ...interface{}) UpdateCondition {
	return builder.Append(b, "Prefixes", Expr(sql, args...)).(UpdateBuilder)
}

// Table sets the table to be updated.
func (b UpdateBuilder) Table(table string) UpdateCondition {
	return builder.Set(b, "Table", table).(UpdateBuilder)
}

// Set adds SET clauses to the query.
func (b UpdateBuilder) Set(column string, value interface{}) UpdateCondition {
	return builder.Append(b, "SetClauses", setClause{column: column, value: value}).(UpdateBuilder)
}

func (b UpdateBuilder) IncrBy(column string, num int) UpdateCondition {
	column = fmt.Sprintf("=+=%s = %s+", column, column)
	return builder.Append(b, "SetClauses", setClause{column: column, value: num}).(UpdateBuilder)
}

func (b UpdateBuilder) DecrBy(column string, num int) UpdateCondition {
	column = fmt.Sprintf("=-=%s = %s+", column, column)
	return builder.Append(b, "SetClauses", setClause{column: column, value: num}).(UpdateBuilder)
}

// SetMap is a convenience method which calls .Set for each key/value pair in clauses.
func (b UpdateBuilder) SetMap(clauses map[string]interface{}) UpdateCondition {
	keys := make([]string, len(clauses))
	i := 0
	for key := range clauses {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		val, _ := clauses[key]
		b = b.Set(key, val).(UpdateBuilder)
	}
	return b
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b UpdateBuilder) Where(pred interface{}, args ...interface{}) WhereConditions {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(UpdateBuilder)
}

//Condition
func (b UpdateBuilder) Condition() WhereConditions {
	return builder.Append(b, "WhereParts", newWherePart("")).(UpdateBuilder)
}

//expr
func (b UpdateBuilder) Expr(sql string, args ...interface{}) WhereConditions {
	return builder.Append(b, "WhereParts", newWherePart(expr{sql: sql, args: args})).(UpdateBuilder)
}

//eq
func (b UpdateBuilder) Eq(column string, arg interface{}) WhereConditions {
	return b.Where(Eq{column: arg})
}

func (b UpdateBuilder) NotEq(column string, arg interface{}) WhereConditions {
	return b.Where(NotEq{column: arg})
}

//gt
func (b UpdateBuilder) Gt(column string, arg interface{}) WhereConditions {
	return b.Where(Gt{column: arg})
}

//gtOrEq
func (b UpdateBuilder) GtOrEq(column string, arg interface{}) WhereConditions {
	return b.Where(GtOrEq{column: arg})
}

//lt
func (b UpdateBuilder) Lt(column string, arg interface{}) WhereConditions {
	return b.Where(Lt{column: arg})
}

//ltOrEq
func (b UpdateBuilder) LtOrEq(column string, arg interface{}) WhereConditions {
	return b.Where(LtOrEq{column: arg})
}

// OrderBy adds ORDER BY expressions to the query.
func (b UpdateBuilder) OrderBy(orderBys ...string) WhereConditions {
	return builder.Extend(b, "OrderBys", orderBys).(UpdateBuilder)
}

// GroupBy adds GROUP BY expressions to the query.
func (b UpdateBuilder) GroupBy(groupBys ...string) WhereConditions {
	return builder.Extend(b, "GroupBys", groupBys).(UpdateBuilder)
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b UpdateBuilder) Having(pred interface{}, rest ...interface{}) WhereConditions {
	return builder.Append(b, "HavingParts", newWherePart(pred, rest...)).(UpdateBuilder)
}

// Limit sets a LIMIT clause on the update.
func (b UpdateBuilder) Limit(limit int) WhereConditions {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(UpdateBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b UpdateBuilder) Offset(offset int) WhereConditions {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(UpdateBuilder)
}

// Suffix adds an expression to the end of the query
func (b UpdateBuilder) Suffix(sql string, args ...interface{}) WhereConditions {
	return builder.Append(b, "Suffixes", Expr(sql, args...)).(UpdateBuilder)
}
