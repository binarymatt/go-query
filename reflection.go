package query

import (
	"reflect"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func reflectDetails(model any) (tableName string, fields []field) {
	pluralize := pluralize.NewClient()
	m := reflect.ValueOf(model)
	v := reflect.Indirect(m)
	typeOfS := v.Type()
	tableName = pluralize.Plural(strcase.ToSnake(typeOfS.Name()))
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			field := retrieveField(v, i)
			fields = append(fields, field)
		}
	}
	return
}

type field struct {
	name      string
	value     any
	omitEmpty bool
	tagged    bool
	excluded  bool
	isZero    bool
}

func retrieveField(val reflect.Value, index int) field {
	f := field{}
	field := val.Type().Field(index)
	fieldVal := val.Field(index)
	// f.typ = fieldVal.Type()
	f.value = fieldVal.Interface()
	f.name = strcase.ToSnake(field.Name)
	tag := field.Tag.Get("db")
	if tag == "-" {
		f.excluded = true
		f.tagged = true
	} else if tag != "" {
		f.tagged = true
		name, opts, _ := strings.Cut(tag, ",")
		if name != "" {
			f.name = name
			if strings.HasPrefix(opts, "omitempty") {
				f.omitEmpty = true
			}
		}
	}
	f.isZero = fieldVal.IsZero()
	return f
}
