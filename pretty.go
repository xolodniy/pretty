package pretty

import (
	"fmt"
	"reflect"
	"strings"
)

func Print(input interface{}) string {
	v := reflect.ValueOf(input)

	switch v.Kind() {
	case reflect.Ptr:
		return "*" + Print(v.Elem().Interface())
	case reflect.Struct:
		typeOfS := v.Type()
		var fields string
		for i := 0; i < v.NumField(); i++ {
			f := getFieldValue(typeOfS.Field(i).Name, v.Field(i))
			if f != "" {
				fields += f + ", "
			}
		}
		fields = strings.TrimSuffix(fields, ", ")
		output := fmt.Sprintf("%s{%s}", getType(input), fields)
		return output
	case reflect.Invalid:
		return "nil"
	default:
		return fmt.Sprintf("%s{%v}", v.Kind().String(), v.Interface())
	}
}

func getFieldValue(fieldName string, v reflect.Value) string {
	var s interface{}
	var invalid bool
	switch {
	case v.Kind() == reflect.Ptr:
		return "*" + getFieldValue(fieldName, v.Elem())
	case v.Kind() == reflect.String && v.String() != "":
		s = v.String()
	case v.Kind() == reflect.Int && v.Int() != 0:
		s = v.Int()
	case v.Kind() == reflect.Struct:
		s = printChildStruct(v)
	default:
		invalid = true
	}
	if invalid || s == "" {
		return ""
	}
	return fmt.Sprintf("%s: %v", fieldName, s)
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func printChildStruct(v reflect.Value) string {
	typeOfS := v.Type()
	var fields string
	for i := 0; i < v.NumField(); i++ {
		f := getFieldValue(typeOfS.Field(i).Name, v.Field(i))
		if f != "" {
			fields += f + ", "
		}
	}
	if fields == "" {
		return ""
	}
	fields = strings.TrimSuffix(fields, ", ")
	return fmt.Sprintf("%s{%s}", v.Type(), fields)
}
