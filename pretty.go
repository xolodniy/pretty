package pretty

import (
	"fmt"
	"reflect"
	"strings"
)

func Print(input interface{}) (output string) {
	defer func() {
		if r := recover(); r != nil {
			output = fmt.Sprintf("panic recovered (%v): %s", r, output)
		}
	}()
	v := reflect.ValueOf(input)

	switch v.Kind() {
	case reflect.Ptr:
		return "*" + Print(v.Elem().Interface())
	case reflect.Struct:
		typeOfS := v.Type()
		var fields string
		for i := 0; i < v.NumField(); i++ {
			f := getFieldValue(typeOfS.Field(i).Name, v.Field(i), false)
			if f != "" {
				fields += fmt.Sprintf("%s: %v, ", typeOfS.Field(i).Name, f)
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

func getFieldValue(fieldName string, v reflect.Value, isPointer bool) string {
	var s interface{}
	var invalid bool
	switch {
	case v.Kind() == reflect.Ptr:
		if v.IsNil() {
			return ""
		} else {
			return getFieldValue(fieldName, v.Elem(), true)
		}
	case v.Kind() == reflect.String && (v.String() != "" || isPointer):
		s = fmt.Sprintf("'%s'", v.String())
	case v.Kind() == reflect.Int && (v.Int() != 0 || isPointer):
		s = v.Int()
	case v.Kind() == reflect.Int64 && (v.Int() != 0 || isPointer):
		s = v.Int()
	case v.Kind() == reflect.Uint && (v.Uint() != 0 || isPointer):
		s = v.Uint()
	case v.Kind() == reflect.Bool && (v.Bool() != false || isPointer):
		s = v.Bool()
	case v.Kind() == reflect.Float64 && (v.Float() != 0 || isPointer):
		s = v.Float()
	case v.Kind() == reflect.Struct:
		s = printChildStruct(v)
	default:
		invalid = true
	}
	if invalid || s == "" {
		return ""
	}
	return fmt.Sprint(s)
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
		f := getFieldValue(typeOfS.Field(i).Name, v.Field(i), false)
		if f != "" {
			fields += fmt.Sprintf("%s: %v, ", typeOfS.Field(i).Name, f)
		}
	}
	if fields == "" {
		return ""
	}
	fields = strings.TrimSuffix(fields, ", ")
	return fmt.Sprintf("%s{%s}", v.Type(), fields)
}
