package squirrel

import "testing"

func TestCondition(t *testing.T) {
	s:=Select("a").From("b").Eq("a",1)
	a(s)
}

func a(c Conditions)  {
	var b interface{}
	c=c.Eq("a",b)
}