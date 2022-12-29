package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type ChildStruct struct {
	Value int `validate:"required"`
}

type ParentStruct struct {
	Child ChildStruct
}

func myValidate(fl validator.FieldLevel) bool {
	fmt.Println("INSIDE MY VALIDATOR") // <- This is never printed
	return false
}
func childStructCustomTypeFunc(field reflect.Value) interface{} {
	if value, ok := field.Interface().(ChildStruct); ok {
		return value.Value
	}
	return nil
}
func main() {
	validator := validator.New()
	validator.RegisterValidation("myValidate", myValidate)
	validator.RegisterCustomTypeFunc(childStructCustomTypeFunc, ChildStruct{})
	data := &ParentStruct{
		Child: ChildStruct{
			//	Value: 10,
		},
	}
	validateErr := validator.Struct(data)
	if validateErr != nil { // <- This is always nil since MyValidate is never called
		fmt.Println("GOT ERROR")
		fmt.Println(validateErr)
	}
	fmt.Println("DONE")
}
