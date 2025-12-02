package query

import (
	"fmt"
	"strings"
)

type part struct {
	text   string
	params []any
}
type Query struct {
	parts []part
}

func Select(model any) *Query {
	q := &Query{}
	tableName, fields := reflectDetails(model)
	columns := []string{}
	for _, field := range fields {
		if !field.excluded {
			columns = append(columns, field.name)
		}
	}
	q.parts = append(q.parts, part{
		text: fmt.Sprintf("SELECT %s FROM %s ", strings.Join(columns, ","), tableName),
	})
	return q

}

func (q *Query) Compile() (string, []any) {
	query := ""
	args := []any{}
	for _, part := range q.parts {
		query += part.text
		// fmt.Println(part.params)
		args = append(args, part.params...)
	}
	return query, args
}

type Opt func(string) (string, []any)

func (q *Query) Where(ops ...Opt) *Query {
	clause := "WHERE "
	attrs := []any{}
	for _, op := range ops {
		newClause, ats := op(clause)
		clause = newClause
		if ats != nil {
			attrs = append(attrs, ats...)
		}
	}
	part := part{
		text:   clause,
		params: attrs,
	}
	q.parts = append(q.parts, part)
	return q
}

func And(opts ...Opt) Opt {
	return group("AND", opts...)
}
func Or(opts ...Opt) Opt {
	return group("OR", opts...)
}

func group(operator string, opts ...Opt) Opt {
	start := "("
	end := ")"
	middle := ""
	attrs := []any{}
	for _, opt := range opts {
		res, ats := opt("")
		if middle == "" {
			middle = res
		} else {
			middle = middle + fmt.Sprintf(" %s ", operator) + res
		}
		if ats != nil {
			attrs = append(attrs, ats...)
		}
	}
	whole := start + middle + end
	return func(c string) (string, []any) {
		return c + whole, attrs
	}
}

func Eq(key string, val any) Opt {
	return func(clause string) (string, []any) {
		clause += fmt.Sprintf("%s = ?", key)
		return clause, []any{val}
	}
}
func Neq(key string, val any) Opt {
	return func(clause string) (string, []any) {
		clause += fmt.Sprintf("%s != ?", key)
		return clause, []any{val}
	}
}
func GreaterEq(key string, val any) Opt {
	return func(s string) (string, []any) {
		s += fmt.Sprintf("%s >= ?", key)
		return s, []any{val}
	}
}
func Greater(key string, val any) Opt {
	return func(s string) (string, []any) {
		s += fmt.Sprintf("%s > ?", key)
		return s, []any{val}
	}

}
func Between(key string, val1, val2 any) Opt {
	return func(s string) (string, []any) {
		s += fmt.Sprintf("(%s >= ? AND %s <= ?)", key, key)
		return s, []any{val1, val2}
	}
}

func In(key string, vals ...any) Opt {
	return func(s string) (string, []any) {
		if len(vals) < 1 {
			return s, nil
		}
		partial := fmt.Sprintf("%s IN (", key)
		l := len(vals)
		for i := range vals {
			partial += " ?"
			if i < l-1 {
				partial += ","
			}
		}
		partial += ")"
		return partial, vals
	}
}
