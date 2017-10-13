package squirrel

import "testing"

func TestCondition(t *testing.T) {
	s:=Select("a")
	a(s)
}

func a(c Conditions)  {
	var b interface{}
	c=c.Eq("a",b)
}