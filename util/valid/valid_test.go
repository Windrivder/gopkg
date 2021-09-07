package valid

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValid(t *testing.T) {
	type User struct {
		Name   string `json:"Name" validate:"required"`
		Age    int    `json:"Age" validate:"required,min=10"`
		Mobile string `json:"Mobile" validate:"required,mobile"`
	}

	if err := RegisterValidation(validators.Mobile()); err != nil {
		t.Fatal(err)
	}

	user := User{Name: "validUser", Age: 9, Mobile: "1882222444"}
	if err := ValidateStruct(user); err != nil {
		rerr := err.(validator.ValidationErrors)

		m := map[string]string{}
		for field, errStr := range rerr.Translate(translator) {
			m[field[strings.Index(field, ".")+1:]] = errStr
		}
		t.Log(m)
	}
}
