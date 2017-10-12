package squirrel

import (
	"fmt"
	"github.com/lann/builder"
)

// StatementBuilderType is the type of StatementBuilder.
type StatementBuilderType builder.Builder

// Select returns a SelectBuilder for this StatementBuilderType.
func (b StatementBuilderType) Select(columns ...string) SelectBuilder {
	return SelectBuilder(b).Columns(columns...)
}

func (b StatementBuilderType) Count(columns string) SelectBuilder {
	str := fmt.Sprintf("COUNT(%s)", columns)
	return SelectBuilder(b).Columns(str)
}

// Insert returns a InsertBuilder for this StatementBuilderType.
func (b StatementBuilderType) Insert(into string) InsertBuilder {
	return InsertBuilder(b).Into(into)
}

// Update returns a UpdateBuilder for this StatementBuilderType.
func (b StatementBuilderType) Update(table string) UpdateBuilder {
	return UpdateBuilder(b).Table(table)
}

// Delete returns a DeleteBuilder for this StatementBuilderType.
func (b StatementBuilderType) Delete(from string) DeleteBuilder {
	return DeleteBuilder(b).From(from)
}

func (b StatementBuilderType) Where(pred interface{}, args ...interface{}) WhereBuilder {
	return WhereBuilder(b).Where(pred, args...)
}

func (b StatementBuilderType) Condition() WhereBuilder {
	return WhereBuilder(b).Where("")
}

func (b StatementBuilderType) Join(join string, rest ...interface{}) JoinBuilder {
	return JoinBuilder(b).Join(join, rest...)
}

func (b StatementBuilderType) JoinClause(pred interface{}, args ...interface{}) JoinBuilder {
	return JoinBuilder(b).JoinClause(pred, args...)
}

func (b StatementBuilderType) LeftJoin(join string, rest ...interface{}) JoinBuilder {
	return JoinBuilder(b).LeftJoin(join, rest...)
}

func (b StatementBuilderType) RightJoin(join string, rest ...interface{}) JoinBuilder {
	return JoinBuilder(b).RightJoin(join, rest...)
}

// PlaceholderFormat sets the PlaceholderFormat field for any child builders.
func (b StatementBuilderType) PlaceholderFormat(f PlaceholderFormat) StatementBuilderType {
	return builder.Set(b, "PlaceholderFormat", f).(StatementBuilderType)
}

// RunWith sets the RunWith field for any child builders.
func (b StatementBuilderType) RunWith(runner BaseRunner) StatementBuilderType {
	return setRunWith(b, runner).(StatementBuilderType)
}

// StatementBuilder is a parent builder for other builders, e.g. SelectBuilder.
var StatementBuilder = StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(Question)

// Select returns a new SelectBuilder, optionally setting some result columns.
//
// See SelectBuilder.Columns.
func Select(columns ...string) SelectBuilder {
	return StatementBuilder.Select(columns...)
}

// Insert returns a new InsertBuilder with the given table name.
//
// See InsertBuilder.Into.
func Insert(into string) InsertBuilder {
	return StatementBuilder.Insert(into)
}

// Update returns a new UpdateBuilder with the given table name.
//
// See UpdateBuilder.Table.
func Update(table string) UpdateBuilder {
	return StatementBuilder.Update(table)
}

// Delete returns a new DeleteBuilder with the given table name.
//
// See DeleteBuilder.Table.
func Delete(from string) DeleteBuilder {
	return StatementBuilder.Delete(from)
}

//新增的where方法
func Where(pred interface{}, args ...interface{}) WhereBuilder {
	return StatementBuilder.Where(pred, args...)
}

func Condition() WhereBuilder {
	return StatementBuilder.Condition()
}

//新增的join方法
func Join(join string, rest ...interface{}) JoinBuilder {
	return StatementBuilder.Join(join, rest...)
}

func JoinClause(pred interface{}, args ...interface{}) JoinBuilder {
	return StatementBuilder.JoinClause(pred, args...)
}

func LeftJoin(join string, rest ...interface{}) JoinBuilder {
	return StatementBuilder.LeftJoin(join, rest...)
}

func RightJoin(join string, rest ...interface{}) JoinBuilder {
	return StatementBuilder.RightJoin(join, rest...)
}

// Case returns a new CaseBuilder
// "what" represents case value
func Case(what ...interface{}) CaseBuilder {
	b := CaseBuilder(builder.EmptyBuilder)

	switch len(what) {
	case 0:
	case 1:
		b = b.what(what[0])
	default:
		b = b.what(newPart(what[0], what[1:]...))

	}
	return b
}
