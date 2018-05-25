package squirrel

import (
	"fmt"
	"github.com/lann/builder"
)

// StatementBuilderType is the type of StatementBuilder.
type StatementBuilderType builder.Builder

// Select returns a SelectBuilder for this StatementBuilderType.
func (b StatementBuilderType) Select(columns ...string) SelectCondition {
	return SelectBuilder(b).Columns(columns...)
}

func (b StatementBuilderType) Count(columns string) SelectCondition {
	str := fmt.Sprintf("COUNT(%s)", columns)
	return SelectBuilder(b).Columns(str)
}

// Insert returns a InsertBuilder for this StatementBuilderType.
func (b StatementBuilderType) Insert(into string) InsertCondition {
	return InsertBuilder(b).Into(into)
}

// Update returns a UpdateBuilder for this StatementBuilderType.
func (b StatementBuilderType) Update(table string) UpdateCondition {
	return UpdateBuilder(b).Table(table)
}

// Delete returns a DeleteBuilder for this StatementBuilderType.
func (b StatementBuilderType) Delete(from string) DeleteCondition {
	return DeleteBuilder(b).From(from)
}

func (b StatementBuilderType) Where(pred interface{}, args ...interface{}) WhereConditions {
	return WhereBuilder(b).Where(pred, args...)
}

func (b StatementBuilderType) Condition() WhereConditions {
	return WhereBuilder(b).Where("")
}

func (b StatementBuilderType) Join(join string, rest ...interface{}) JoinCondition {
	return JoinBuilder(b).Join(join, rest...)
}

func (b StatementBuilderType) JoinClause(pred interface{}, args ...interface{}) JoinCondition {
	return JoinBuilder(b).JoinClause(pred, args...)
}

func (b StatementBuilderType) LeftJoin(join string, rest ...interface{}) JoinCondition {
	return JoinBuilder(b).LeftJoin(join, rest...)
}

func (b StatementBuilderType) RightJoin(join string, rest ...interface{}) JoinCondition {
	return JoinBuilder(b).RightJoin(join, rest...)
}

// PlaceholderFormat sets the PlaceholderFormat field for any child builders.
func (b StatementBuilderType) PlaceholderFormat(f PlaceholderFormat) StatementBuilderType {
	return builder.Set(b, "PlaceholderFormat", f).(StatementBuilderType)
}

// StatementBuilder is a parent builder for other builders, e.g. SelectBuilder.
var StatementBuilder = StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(Question)

// Select returns a new SelectBuilder, optionally setting some result columns.
//
// See SelectBuilder.Columns.
func Select(columns ...string) SelectCondition {
	return StatementBuilder.Select(columns...)
}

// Insert returns a new InsertBuilder with the given table name.
//
// See InsertBuilder.Into.
func Insert(into string) InsertCondition {
	return StatementBuilder.Insert(into)
}

// Update returns a new UpdateBuilder with the given table name.
//
// See UpdateBuilder.Table.
func Update(table string) UpdateCondition {
	return StatementBuilder.Update(table)
}

// Delete returns a new DeleteBuilder with the given table name.
//
// See DeleteBuilder.Table.
func Delete(from string) DeleteCondition {
	return StatementBuilder.Delete(from)
}

//新增的where方法
func Where(pred interface{}, args ...interface{}) WhereConditions {
	return StatementBuilder.Where(pred, args...)
}

func Condition() WhereConditions {
	return StatementBuilder.Condition()
}

//新增的join方法
func Join(join string, rest ...interface{}) JoinCondition {
	return StatementBuilder.Join(join, rest...)
}

func JoinClause(pred interface{}, args ...interface{}) JoinCondition {
	return StatementBuilder.JoinClause(pred, args...)
}

func LeftJoin(join string, rest ...interface{}) JoinCondition {
	return StatementBuilder.LeftJoin(join, rest...)
}

func RightJoin(join string, rest ...interface{}) JoinCondition {
	return StatementBuilder.RightJoin(join, rest...)
}

//// Case returns a new CaseBuilder
//// "what" represents case value
//func Case(what ...interface{}) CaseBuilder {
//	b := CaseBuilder(builder.EmptyBuilder)
//
//	switch len(what) {
//	case 0:
//	case 1:
//		b = b.what(what[0])
//	default:
//		b = b.what(newPart(what[0], what[1:]...))
//
//	}
//	return b
//}
