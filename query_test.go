package query

import (
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

type test struct {
	Id        string
	Name      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time
}

func TestSelect(t *testing.T) {
	statement, args := Select(test{}).Compile()
	must.Eq(t, "SELECT id,name,created_at,updated_at FROM tests ", statement)
	must.SliceEmpty(t, args)
}
