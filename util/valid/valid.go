package valid

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type (
	Validator interface {
		Name() string
		Trans(locale string) string
		Validate(i interface{}) error
	}

	ValidateFunc func(i interface{}) error
)

var (
	v *validator.Validate
)

func NewValidor() *validator.Validate {
	v = validator.New()
	return v
}

func ValidateStruct(i interface{}) error {
	if v == nil {
		v = NewValidor()
	}

	return v.Struct(i)
}

func ValidateVar(field interface{}, tag string) error {
	if v == nil {
		v = NewValidor()
	}

	return v.Var(field, tag)
}

func RegisterValidation(locale string, trans ut.Translator, validators ...Validator) (err error) {
	for _, validator := range validators {
		tag := validator.Name()
		msg := validator.Trans(locale)

		if err = v.RegisterValidation(tag,
			wrapValidateFunc(validator.Validate), true); err != nil {
			return err
		}

		if err = RegisterTranslation(tag, msg, trans); err != nil {
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
