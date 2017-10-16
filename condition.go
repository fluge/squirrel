package squirrel

type WhereConditions interface {
	ToSql() (string, []interface{}, error)
	PlaceholderFormat(PlaceholderFormat) WhereConditions
	Where(interface{}, ...interface{}) WhereConditions
	Condition() WhereConditions
	Expr(string, ...interface{}) WhereConditions
	Eq(string, interface{}) WhereConditions
	NotEq(string, interface{}) WhereConditions
	Gt(string, interface{}) WhereConditions
	GtOrEq(string, interface{}) WhereConditions
	Lt(string, interface{}) WhereConditions
	LtOrEq(string, interface{}) WhereConditions
	OrderBy(...string) WhereConditions
	GroupBy(...string) WhereConditions
	Having(interface{}, ...interface{}) WhereConditions
	Limit(int) WhereConditions
	Offset(int) WhereConditions
	Suffix(string, ...interface{}) WhereConditions
}

type SelectCondition interface {
	Prefix(string, ...interface{}) SelectCondition
	Distinct() SelectCondition
	Options(...string) SelectCondition
	Columns(...string) SelectCondition
	Column(interface{}, ...interface{}) SelectCondition
	From(string) SelectCondition
	FromSelect(SelectCondition, string) SelectCondition
	JoinCondition
}

type JoinCondition interface {
	JoinClause(interface{}, ...interface{}) JoinCondition
	Join(string, ...interface{}) JoinCondition
	LeftJoin(string, ...interface{}) JoinCondition
	RightJoin(string, ...interface{}) JoinCondition
	WhereConditions
}

type UpdateCondition interface {
	Prefix(string, ...interface{}) UpdateCondition
	Table(string) UpdateCondition
	Set(string, interface{}) UpdateCondition
	IncrBy(string, int) UpdateCondition
	DecrBy(string, int) UpdateCondition
	SetMap(map[string]interface{}) UpdateCondition
	WhereConditions
}

type DeleteCondition interface {
	Prefix(string, ...interface{}) DeleteCondition
	From(string) DeleteCondition
	WhereConditions
}
type InsertCondition interface {
	ToSql() (string, []interface{}, error)
	PlaceholderFormat(PlaceholderFormat) InsertCondition
	Prefix(string, ...interface{}) InsertCondition
	Options(...string) InsertCondition
	Into(string) InsertCondition
	Columns(...string) InsertCondition
	Values(...interface{}) InsertCondition
	Suffix(string, ...interface{}) InsertCondition
	SetMap(map[string]interface{}) InsertCondition
	Select(SelectCondition) InsertCondition
}
