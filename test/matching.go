package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type ChildStruct struct {
	Value int
	T     string
}

type ParentStruct struct {
	Type  string      `validate:"required"`
	Child ChildStruct `validate:"emptyCheckChild"`
}

func emptyCheckChild(fl validator.FieldLevel) bool {
	fmt.Println("INSIDE MY VALIDATOR")
	if value, ok := fl.Parent().Interface().(ParentStruct); ok {
		if value.Type == "B" {
			//만약 부모 struct 타입이 B 이면 필드 check
			//	if value.Child.Value == 0 {
			//struct empty check
			//https://freshman.tech/snippets/go/check-empty-struct/
			//참고 https://gist.github.com/NatGra/1972d860e1f4b7d2216813514431b1d9
			if reflect.ValueOf(value.Child).IsZero() {
				return false
			}

		}
	}
	return true
}

func childStructCustomTypeFunc(field reflect.Value) interface{} {
	if value, ok := field.Interface().(ChildStruct); ok {
		return value.Value
	}
	return nil
}

func main() {
	validator := validator.New()
	validator.RegisterValidation("emptyCheckChild", emptyCheckChild)
	validator.RegisterCustomTypeFunc(childStructCustomTypeFunc, ChildStruct{})
	data := &ParentStruct{
		Type: "B",
		Child: ChildStruct{
			Value: 10,
		},
	}
	validateErr := validator.Struct(data)
	if validateErr != nil { // <- This is always nil since MyValidate is never called
		fmt.Println("GOT ERROR")
		fmt.Println(validateErr)
	}
	fmt.Println("DONE")
}
