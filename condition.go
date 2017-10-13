package squirrel

type Conditions interface {
	ToSql() (string, []interface{}, error)
	Where(interface{}, ...interface{}) Conditions
	Eq(string, interface{}) Conditions
	Gt(string, interface{}) Conditions
	GtOrEq(string, interface{}) Conditions
	Lt(string, interface{}) Conditions
	LtOrEq(string, interface{}) Conditions
	OrderBy(...string) Conditions
	Limit(int) Conditions
	Offset(int) Conditions
	Suffix(string, ...interface{}) Conditions
	//GroupBy(...string) Conditions
	//Having(interface{},...interface{}) Conditions
}

type JoinCondition interface {
	Join(string, ...interface{}) JoinCondition
	LeftJoin(string, ...interface{}) JoinCondition
	RightJoin(string, ...interface{}) JoinCondition
	ToSql() (string, []interface{}, error)
	Where(interface{}, ...interface{}) Conditions
	Eq(string, interface{}) Conditions
	Gt(string, interface{}) Conditions
	GtOrEq(string, interface{}) Conditions
	Lt(string, interface{}) Conditions
	LtOrEq(string, interface{}) Conditions
	OrderBy(...string) Conditions
	Limit(int) Conditions
	Offset(int) Conditions
	Suffix(string, ...interface{}) Conditions
}
