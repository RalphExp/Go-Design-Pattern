package main

import (
	"fmt"
	"reflect"
)

type Secret struct {
	Username string
	Password string
}

type Record struct {
	Field1 string
	Field2 float64
	Field3 Secret
}

func (r Record) String() string {
	return fmt.Sprintf("{Field1: %s, Field2: %f , Field3: {Username: %s, Password: %s}}",
		r.Field1, r.Field2, r.Field3.Password, r.Field3.Username)
}

func main() {
	A := Record{"String value", -12.123, Secret{"Mihalis", "Tsoukalos"}}
	r := reflect.ValueOf(A)
	fmt.Printf("r's value: %v, r's type: %v\n", r.String(), r.Kind().String())

	iType := r.Type()
	fmt.Printf("i Type: %s\n", iType)
	fmt.Printf("The %d fields of %s are\n", r.NumField(), iType)
	for i := 0; i < r.NumField(); i++ {
		fmt.Printf("\t%s ", iType.Field(i).Name)
		fmt.Printf("\twith type: [%s] ", r.Field(i).Type())
		fmt.Printf("\tand value: [%v]\n", r.Field(i).Interface())

		k := reflect.TypeOf(r.Field(i).Interface()).Kind()
		if k == reflect.Struct {
			fmt.Println(r.Field(i).Type())
		}
	}

	fmt.Printf("String value: %s\n", A)
}
