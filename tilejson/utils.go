package tilejson

import (
	"reflect"

	"github.com/go-courier/reflectx"
)

func StructToFields(v interface{}) map[string]FieldType {
	structType := reflectx.Deref(reflect.TypeOf(v))
	if structType.Kind() != reflect.Struct {
		return nil
	}
	fields := map[string]FieldType{}
	for i := 0; i < structType.NumField(); i++ {
		ft := structType.Field(i)
		name, ok := ft.Tag.Lookup("name")
		if ok {
			switch ft.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				fields[name] = FieldTypeNumber
			case reflect.String:
				fields[name] = FieldTypeString
			case reflect.Bool:
				fields[name] = FieldTypeBoolean
			}
		}
	}

	return fields
}

func StructToProperties(v interface{}) map[string]interface{} {
	s := reflectx.Indirect(reflect.ValueOf(v))
	if s.Kind() != reflect.Struct {
		return nil
	}
	typ := s.Type()
	props := map[string]interface{}{}
	for i := 0; i < s.NumField(); i++ {
		ft := typ.Field(i)
		name, ok := ft.Tag.Lookup("name")
		if ok {
			props[name] = s.Field(i).Interface()
		}
	}
	return props
}
