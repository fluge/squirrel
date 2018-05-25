package squirrel

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestWhere(t *testing.T) {
	sql, args, err := Where(Eq{"username": []string{"moe", "larry", "curly", "shemp"}}).ToSql()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, sql, " WHERE username IN (?,?,?,?)")
	assert.Equal(t, args, []interface{}{"moe", "larry", "curly", "shemp"})
}
