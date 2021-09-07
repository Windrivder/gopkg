package valid

import (
	"github.com/go-playground/validator/v10"
	"github.com/windrivder/gopkg/container/typex"
)

type (
	Validator interface {
		Name() string
		Trans() typex.DictStrs
		Validate(i interface{}) error
	}

	ValidateFunc func(i interface{}) error
)

var (
	v = validator.New()
)

func GetValidate() *validator.Validate {
	return v
}

func ValidateStruct(i interface{}) error {
	return v.Struct(i)
}

func ValidateVar(field interface{}, tag string) error {
	return v.Var(field, tag)
}

func RegisterValidation(validators ...Validator) (err error) {
	for _, validator := range validators {
		tag := validator.Name()
		msg := validator.Trans()[defaultLocale]

		if err = v.RegisterValidation(tag,
			wrapValidateFunc(validator.Validate), true); err != nil {
			return err
		}

		if err = RegisterTranslation(tag, msg); err != nil {
			return err
		}
	}

	return nil
}

func wrapValidateFunc(fn ValidateFunc) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return fn(fl.Field().Interface()) == nil
	}
}
