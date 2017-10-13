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
