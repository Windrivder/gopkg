package valid

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/i18n"
)

type (
	Validator interface {
		Name() string
		Trans(locale string) string
		Validate(i interface{}) error
	}

	ValidateFunc func(i interface{}) error
)

type Validate struct {
	v     *validator.Validate
	o     i18n.Options
	trans ut.Translator
}

func NewValidate(o i18n.Options) (*Validate, error) {
	v := validator.New()
	trans, err := newTranslator(o, v)
	if err != nil {
		return nil, err
	}

	return &Validate{v: v, o: o, trans: trans}, nil
}

func newTranslator(o i18n.Options, v *validator.Validate) (trans ut.Translator, err error) {
	en := en.New()
	zh := zh.New()

	uni := ut.New(en, en, zh)
	trans, ok := uni.GetTranslator(o.Locale)
	if !ok {
		return nil, errorx.New("validator get translator fail")
	}

	switch o.Locale {

	case i18n.LocaleEN.String():
		err = en_translations.RegisterDefaultTranslations(v, trans)

	case i18n.LocaleZH.String():
		err = zh_translations.RegisterDefaultTranslations(v, trans)

	default:
		err = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	if err != nil {
		return nil, err
	}

	return trans, nil
}

func (v *Validate) Locale() string {
	return v.o.Locale
}

func (v *Validate) Translator() ut.Translator {
	return v.trans
}

func (v *Validate) ValidateStruct(i interface{}) error {
	return v.v.Struct(i)
}

func (v *Validate) ValidateVar(field interface{}, tag string) error {
	return v.v.Var(field, tag)
}

func (v *Validate) RegisterValidation(validators ...Validator) (err error) {
	for _, vdtor := range validators {
		tag := vdtor.Name()
		msg := vdtor.Trans(v.Locale())

		if err = v.v.RegisterValidation(tag,
			wrapValidateFunc(vdtor.Validate), true); err != nil {
			return err
		}

		if err = v.v.RegisterTranslation(tag, v.Translator(),
			func(translator ut.Translator) error {
				return translator.Add(tag, msg, true)
			},
			func(translator ut.Translator, fe validator.FieldError) string {
				message, _ := translator.T(tag, fe.Field())
				return message
			},
		); err != nil {
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
