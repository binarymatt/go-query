package query

import (
	"fmt"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

type Sample struct {
	Tenant    string
	Pattern   string
	Endpoint  string
	CreatedAt time.Time `db:"created_at"`
	DeletedAt time.Time `db:"deleted_at,omitempty"`
}

func TestInsert_HappyPath(t *testing.T) {
	n := time.Now()
	m := &Sample{
		Tenant:    "test",
		Pattern:   "{}",
		Endpoint:  "localhost",
		CreatedAt: n,
	}
	query, values := Insert(m)
	fmt.Println(query)
	must.Eq(t, "INSERT INTO samples (tenant,pattern,endpoint,created_at) VALUES (?,?,?,?)", query)
	must.Eq(t, []any{"test", "{}", "localhost", n}, values)
}

func TestInsert_PrivateFields(t *testing.T) {
	type test struct {
		Id    int
		name  string
		Value string
	}
	tVal := test{
		Id:    1,
		name:  "test",
		Value: "val1",
	}
	statement, values := Insert(tVal)
	must.Eq(t, "INSERT INTO tests (id,value) VALUES (?,?)", statement)
	must.Eq(t, []any{1, "val1"}, values)
}

func TestInsert_Snaked(t *testing.T) {
	type SnakeTest struct {
		Id int
	}
	st := SnakeTest{
		Id: 1,
	}
	statement, values := Insert(st)

	must.Eq(t, "INSERT INTO snake_tests (id) VALUES (?)", statement)
	must.Eq(t, []any{1}, values)
}
