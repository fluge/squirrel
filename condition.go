package squirrel


// DeleteBuilder builds SQL DELETE statements.
type Conditions interface {
	ToSql()(interface{},interface{},error)
	Where(interface{},...interface{})interface{}
	Eq(string,interface{})interface{}
	Gt(string,interface{})interface{}
	GtOrEq(string,interface{})interface{}
	Lt(string,interface{})interface{}
	LtOrEq(string,interface{})interface{}
	OrderBy(...string)interface{}
	Limit(int)interface{}
	Offset(int)interface{}
	GroupBy(...string)interface{}
	Having(interface{},...interface{})interface{}
}



