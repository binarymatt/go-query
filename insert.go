package query

import (
	"fmt"
	"strings"
)

func Insert(model any) (string, []any) {
	// tableName, columns, values := reflectModelDetails(model)
	tableName, fields := reflectDetails(model)
	columns := []string{}
	values := []any{}
	binds := []string{}
	for _, field := range fields {
		// if field is excluded skip
		if field.excluded || (field.omitEmpty && field.isZero) {
			continue
		}
		columns = append(columns, field.name)
		binds = append(binds, "?")
		values = append(values, field.value)
	}
	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ","), strings.Join(binds, ","))
	return statement, values
}
