package database

import (
	"reflect"
	"database/sql"
	//"database/sql/driver"
	"github.com/gobuffalo/flect"
)

type Loader struct {
	model reflect.Value
}

type ColumnType struct {
	*sql.ColumnType
	Index int
}
	
func GetStructFields(rows sql.Rows, v interface{}) {
	v1 := reflect.ValueOf(v)
	
	if v1.Kind() != reflect.Ptr {
		panic("Expecting a Struct Pointer")
	}
	
	v1 = reflect.Indirect(v1)
	t1 := v1.Type()
	
	_columns, err := rows.ColumnTypes()
	if err != nil {
		panic(err)
	}
	
	var columns = make(map[string]ColumnType)
	var pointers = make([]interface{}, len(_columns))
	for i, column := range _columns {
		columns[column.Name()] = ColumnType{ column, i } 
	}
	
	for i := 0; i < t1.NumField(); i++ {
		fieldName := t1.Field(i).Tag.Get("json")
		if(fieldName == "") {
			fieldName = flect.Camelize(t1.Field(i).Name)
		}
		
		if(t1.Field(i).Type.Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem())) {
			pointers[columns[fieldName].Index] = v1.Field(i).Interface()
		}
		
		if(columns[fieldName].ScanType().AssignableTo(t1.Field(i).Type)) {
			pointers[columns[fieldName].Index] = v1.Field(i).Interface()
		}
	}
}
