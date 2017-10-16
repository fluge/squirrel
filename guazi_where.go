package squirrel

import (
	"bytes"
	"fmt"
	"github.com/lann/builder"
	"strings"
)

type whereData struct {
	PlaceholderFormat PlaceholderFormat
	WhereParts        []Sqlizer
	GroupBys          []string
	HavingParts       []Sqlizer
	OrderBys          []string
	Limit             string
	Offset            string
	Suffixes          exprs
}

func (d *whereData) ToSql() (sqlStr string, args []interface{}, err error) {
	sql := &bytes.Buffer{}
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
func (b WhereBuilder) PlaceholderFormat(f PlaceholderFormat) WhereConditions {
	return builder.Set(b, "PlaceholderFormat", f).(WhereBuilder)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b WhereBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(whereData)
	return data.ToSql()
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
func (b WhereBuilder) Where(pred interface{}, args ...interface{}) WhereConditions {
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(WhereBuilder)
}

//Condition
func (b WhereBuilder) Condition() WhereConditions {
	return builder.Append(b, "WhereParts", newWherePart("")).(WhereBuilder)
}

//expr
func (b WhereBuilder) Expr(sql string, args ...interface{}) WhereConditions {
	return b.Where(Expr(sql, args))
}

//eq
func (b WhereBuilder) Eq(column string, arg interface{}) WhereConditions {
	return b.Where(Eq{column: arg})
}

//gt
func (b WhereBuilder) Gt(column string, arg interface{}) WhereConditions {
	return b.Where(Gt{column: arg})
}

//gtOrEq
func (b WhereBuilder) GtOrEq(column string, arg interface{}) WhereConditions {
	return b.Where(GtOrEq{column: arg})
}

//lt
func (b WhereBuilder) Lt(column string, arg interface{}) WhereConditions {
	return b.Where(Lt{column: arg})
}

//ltOrEq
func (b WhereBuilder) LtOrEq(column string, arg interface{}) WhereConditions {
	return b.Where(LtOrEq{column: arg})
}

//or
func (b WhereBuilder) Or(pred ...interface{}) WhereBuilder {
	or := Or{}
	for _, v := range pred {
		switch t := v.(type) {
		case expr:
			or = append(or, t)
		case Gt:
			or = append(or, t)
		case Eq:
			or = append(or, t)
		case GtOrEq:
			or = append(or, t)
		case LtOrEq:
			or = append(or, t)
		case Lt:
			or = append(or, t)
		case And:
			and := And{}
			and = append(and, t)
			or = append(or, and...)
		default:
			panic("unsport ")
		}
	}
	return builder.Append(b, "WhereParts", newWherePart(or)).(WhereBuilder)
}

// GroupBy adds GROUP BY expressions to the query.
func (b WhereBuilder) GroupBy(groupBys ...string) WhereConditions {
	return builder.Extend(b, "GroupBys", groupBys).(WhereBuilder)
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b WhereBuilder) Having(pred interface{}, rest ...interface{}) WhereConditions {
	return builder.Append(b, "HavingParts", newWherePart(pred, rest...)).(WhereBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b WhereBuilder) OrderBy(orderBys ...string) WhereConditions {
	return builder.Extend(b, "OrderBys", orderBys).(WhereBuilder)
}

// Limit sets a LIMIT clause on the query.
func (b WhereBuilder) Limit(limit int) WhereConditions {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(WhereBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b WhereBuilder) Offset(offset int) WhereConditions {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(WhereBuilder)
}

// Suffix adds an expression to the end of the query
func (b WhereBuilder) Suffix(sql string, args ...interface{}) WhereConditions {
	return builder.Append(b, "Suffixes", Expr(sql, args...)).(WhereBuilder)
}
