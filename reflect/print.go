package format

import (
	"fmt"
	"reflect"
	"strconv"
)

// Any formats any value as a string.
func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

func formatArray(v reflect.Value) (s string) {
	// v is slice or array
	s += "["
	if v.Len() > 0 {
		for i := range v.Len() - 1 {
			s += formatAtom(v.Index(i))
			s += " "
		}
		s += formatAtom(v.Index(v.Len() - 1))
	}
	s += "]"
	return s
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', 7, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 15, 64)
	case reflect.Complex64:
		return strconv.FormatComplex(v.Complex(), 'f', 7, 64)
	case reflect.Complex128:
		return strconv.FormatComplex(v.Complex(), 'f', 15, 128)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Array, reflect.Slice:
		return formatArray(v)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	fmt.Printf(display(name, reflect.ValueOf(x)))
}

func display(path string, v reflect.Value) (s string) {
	switch v.Kind() {
	case reflect.Invalid:
		s = fmt.Sprintf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			s += display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			s += display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			s += display(fmt.Sprintf("%s[%s]", path,
				formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			s = fmt.Sprintf("%s = nil\n", path)
		} else {
			s = display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			s = fmt.Sprintf("%s = nil\n", path)
		} else {
			s = fmt.Sprintf("%s.type = %s\n", path, v.Elem().Type())
			s += display(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		s = fmt.Sprintf("%s = %s\n", path, formatAtom(v))
	}
	return s
}
