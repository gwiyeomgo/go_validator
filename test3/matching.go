package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
)

type ChildStruct struct {
	Value int    `json:"value"`
	T     string `json:"t"`
}

type ParentStruct struct {
	Type  string      `json:"type" validate:"required"`
	Child ChildStruct `json:"child" validate:"emptyCheckChild"`
}

func (p ParentStruct) Validate(ctx echo.Context) error {
	return ctx.Validate(p)
}
func emptyCheckChild(fl validator.FieldLevel) bool {
	fmt.Println("INSIDE MY VALIDATOR")
	if value, ok := fl.Parent().Interface().(ParentStruct); ok {
		if value.Type == "B" {
			if reflect.ValueOf(value.Child).IsZero() {
				return false
			}
		}
	}
	return true
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	echoApp := echo.New()
	validator := validator.New()
	validator.RegisterValidation("emptyCheckChild", emptyCheckChild)
	echoApp.Validator = &CustomValidator{validator: validator}

	requestBody := `{ "type": "B" }`
	req := httptest.NewRequest(http.MethodPost, "/api/test", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	var content ParentStruct
	if err := ctx.Bind(&content); err != nil {
		fmt.Println("error")
	}
	if err := content.Validate(ctx); err != nil {
		fmt.Println("error")
	}
	/*data := &ParentStruct{
		Type:  "Joined",
		Child: ChildStruct{
			//Value: 10,
		},
	}*/
}
