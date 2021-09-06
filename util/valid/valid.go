package valid

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/windrivder/gopkg/container/typex"
	"github.com/windrivder/gopkg/i18n"
)

type (
	Validator interface {
		Name() string
		Trans() typex.DictStrs
		Validate(i interface{}) error
	}

	ValidateFunc func(i interface{}) error
)

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, true); err != nil {
			return err
		}
		return nil
	}
}

func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, _ := trans.T(fe.Tag(), fe.Field())
	return msg
}

func init() {
	validators := []Validator{
		Mobile(),
	}

	// translator, err := NewTranslator(i18n.LocateEN.String())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for _, valid := range validators {
		tag := valid.Name()
		msg := valid.Trans()[locate]

		RegisterValidation(tag, valid.Validate, true)
		RegisterTranslation(tag, translator, registerTranslator(tag, msg), translate)
	}
}

var (
	locate        = i18n.LocateZH.String()
	v             = validator.New()
	translator, _ = NewTranslator(locate)
)

func ValidateStruct(i interface{}) error {
	return v.Struct(i)
}

func ValidateVar(field interface{}, tag string) error {
	return v.Var(field, tag)
}

func RegisterValidation(tag string, fn ValidateFunc, callValidationEvenIfNull ...bool) error {
	return v.RegisterValidation(tag, wrapValidateFunc(fn), callValidationEvenIfNull...)
}

func RegisterTranslation(tag string, trans ut.Translator,
	registerFn validator.RegisterTranslationsFunc,
	translationFn validator.TranslationFunc) error {
	return v.RegisterTranslation(tag, trans, registerFn, translationFn)
}

func wrapValidateFunc(fn ValidateFunc) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return fn(fl.Field().Interface()) == nil
	}
}
