package pretty

import (
	"fmt"
	"reflect"
	"strings"
	"time"
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
		if v.Elem().IsValid() {
			return "*" + Print(v.Elem().Interface())
		}
		return v.Type().String() + "{nil}"
	case reflect.Struct:
		if getType(input) == "Time" {
			return fmt.Sprint(input)
		}
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
	case reflect.Slice:
		elements := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			elements[i] = Print(v.Index(i).Interface())
		}
		return fmt.Sprintf("%s: [%s]", reflect.TypeOf(input), strings.Join(elements, ", "))
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
	if typeOfS.String() == "time.Time" {
		if v.IsZero() {
			return ""
		}
		if v.CanInterface() {
			return v.Interface().(time.Time).Format(time.RFC3339)
		}
		// TODO: research how to print unexported date
		return "unexported date"
		// here we re-create a date field which might be private and unaddressable
		//return makeTime(v).String()
	}
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

// HACK: special function for translate time. used for avoid unexported dates in structs
// This func is not work correctly, magic numbers does not help
func makeTime(v reflect.Value) Time {
	var (
		magicShift int64 // I don't know what this a number, but it allows make correct unix timestamp
		timezone   *time.Location
	)

	timeFieldLoc := v.Field(2).Elem()
	if timeFieldLoc.Kind() == reflect.Invalid {
		magicShift = 1636324347739234200
		timezone = time.UTC
	} else {
		timeFieldLocFieldCacheZone := timeFieldLoc.Field(6).Elem()
		if timeFieldLocFieldCacheZone.Kind() == reflect.Invalid {
			magicShift = 1628420593
			timezone = time.UTC
		} else {
			timeShift := timeFieldLocFieldCacheZone.Field(1).Int()
			timezone = time.FixedZone(" ", int(timeShift))
			magicShift = -62135596800
			magicShift = 1628420593
		}
	}

	unix := v.Field(1).Int() + magicShift
	date := time.Unix(unix, int64(v.Field(0).Uint()))
	return Time(date.In(timezone))
}

type Time time.Time

func (t Time) String() string {
	if time.Time(t).IsZero() {
		return ""
	}
	s := time.Time(t).String()
	return strings.TrimSpace(s)
}
