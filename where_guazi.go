package squirrel

import (
	"bytes"
	//"database/sql"
	"fmt"
	"strings"

	"github.com/lann/builder"
)

type whereData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          exprs
	Options           []string
	Joins             []Sqlizer
	WhereParts        []Sqlizer
	GroupBys          []string
	HavingParts       []Sqlizer
	OrderBys          []string
	Limit             string
	Offset            string
	Suffixes          exprs
}

//func (d *whereData) Exec() (sql.Result, error) {
//	if d.RunWith == nil {
//		return nil, RunnerNotSet
//	}
//	return ExecWith(d.RunWith, d)
//}
//
//func (d *whereData) Query() (*sql.Rows, error) {
//	if d.RunWith == nil {
//		return nil, RunnerNotSet
//	}
//	return QueryWith(d.RunWith, d)
//}
//
//func (d *whereData) QueryRow() RowScanner {
//	if d.RunWith == nil {
//		return &Row{err: RunnerNotSet}
//	}
//	queryRower, ok := d.RunWith.(QueryRower)
//	if !ok {
//		return &Row{err: RunnerNotQueryRunner}
//	}
//	return QueryRowWith(queryRower, d)
//}

func (d *whereData) ToSql() (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, _ = d.Prefixes.AppendToSql(sql, " ", args)
		sql.WriteString(" ")
	}

	if len(d.Options) > 0 {
		sql.WriteString(strings.Join(d.Options, " "))
		sql.WriteString(" ")
	}

	if len(d.Joins) > 0 {
		sql.WriteString(" ")
		args, err = appendToSql(d.Joins, sql, " ", args)
		if err != nil {
			return
		}
	}

	if len(d.WhereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(d.WhereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.GroupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(d.GroupBys, ", "))
	}

	if len(d.HavingParts) > 0 {
		sql.WriteString(" HAVING ")
		args, err = appendToSql(d.HavingParts, sql, " AND ", args)
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

// WhereBuilder builds SQL SELECT statements.
type WhereBuilder builder.Builder

func init() {
	builder.Register(WhereBuilder{}, whereData{})
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
//func (b WhereBuilder) PlaceholderFormat(f PlaceholderFormat) WhereBuilder {
//	return builder.Set(b, "PlaceholderFormat", f).(WhereBuilder)
//}
//
//// Runner methods
//
//// RunWith sets a Runner (like database/sql.DB) to be used with e.g. Exec.
//func (b WhereBuilder) RunWith(runner BaseRunner) WhereBuilder {
//	return setRunWith(b, runner).(WhereBuilder)
//}

//// Exec builds and Execs the query with the Runner set by RunWith.
//func (b WhereBuilder) Exec() (sql.Result, error) {
//	data := builder.GetStruct(b).(whereData)
//	return data.Exec()
//}
//
//// Query builds and Querys the query with the Runner set by RunWith.
//func (b WhereBuilder) Query() (*sql.Rows, error) {
//	data := builder.GetStruct(b).(whereData)
//	return data.Query()
//}
//
//// QueryRow builds and QueryRows the query with the Runner set by RunWith.
//func (b WhereBuilder) QueryRow() RowScanner {
//	data := builder.GetStruct(b).(whereData)
//	return data.QueryRow()
//}

//// Scan is a shortcut for QueryRow().Scan.
//func (b WhereBuilder) Scan(dest ...interface{}) error {
//	return b.QueryRow().Scan(dest...)
//}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b WhereBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(whereData)
	return data.ToSql()
}

// Prefix adds an expression to the beginning of the query
func (b WhereBuilder) Prefix(sql string, args ...interface{}) WhereBuilder {
	return builder.Append(b, "Prefixes", Expr(sql, args...)).(WhereBuilder)
}

// Distinct adds a DISTINCT clause to the query.
func (b WhereBuilder) Distinct() WhereBuilder {
	return b.Options("DISTINCT")
}

// Options adds select option to the query
func (b WhereBuilder) Options(options ...string) WhereBuilder {
	return builder.Extend(b, "Options", options).(WhereBuilder)
}

// JoinClause adds a join clause to the query.
func (b WhereBuilder) JoinClause(pred interface{}, args ...interface{}) WhereBuilder {
	return builder.Append(b, "Joins", newPart(pred, args...)).(WhereBuilder)
}

// Join adds a JOIN clause to the query.
func (b WhereBuilder) Join(join string, rest ...interface{}) WhereBuilder {
	return b.JoinClause("JOIN "+join, rest...)
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b WhereBuilder) LeftJoin(join string, rest ...interface{}) WhereBuilder {
	return b.JoinClause("LEFT JOIN "+join, rest...)
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b WhereBuilder) RightJoin(join string, rest ...interface{}) WhereBuilder {
	return b.JoinClause("RIGHT JOIN "+join, rest...)
}

// Where adds an expression to the WHERE clause of the query.
//
// Expressions are ANDed together in the generated SQL.
//
// Where accepts several types for its pred argument:
//
// nil OR "" - ignored.
//
// string - SQL expression.
// If the expression has SQL placeholders then a set of arguments must be passed
// as well, one for each placeholder.
//
// map[string]interface{} OR Eq - map of SQL expressions to values. Each key is
// transformed into an expression like "<key> = ?", with the corresponding value
// bound to the placeholder. If the value is nil, the expression will be "<key>
// IS NULL". If the value is an array or slice, the expression will be "<key> IN
// (?,?,...)", with one placeholder for each item in the value. These expressions
// are ANDed together.
//
// Where will panic if pred isn't any of the above types.
func (b WhereBuilder) Where(pred interface{}, args ...interface{}) WhereBuilder {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(WhereBuilder)
}

// GroupBy adds GROUP BY expressions to the query.
func (b WhereBuilder) GroupBy(groupBys ...string) WhereBuilder {
	return builder.Extend(b, "GroupBys", groupBys).(WhereBuilder)
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b WhereBuilder) Having(pred interface{}, rest ...interface{}) WhereBuilder {
	return builder.Append(b, "HavingParts", newWherePart(pred, rest...)).(WhereBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b WhereBuilder) OrderBy(orderBys ...string) WhereBuilder {
	return builder.Extend(b, "OrderBys", orderBys).(WhereBuilder)
}

// Limit sets a LIMIT clause on the query.
func (b WhereBuilder) Limit(limit uint64) WhereBuilder {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(WhereBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b WhereBuilder) Offset(offset uint64) WhereBuilder {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(WhereBuilder)
}

// Suffix adds an expression to the end of the query
func (b WhereBuilder) Suffix(sql string, args ...interface{}) WhereBuilder {
	return builder.Append(b, "Suffixes", Expr(sql, args...)).(WhereBuilder)
}
