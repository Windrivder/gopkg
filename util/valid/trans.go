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

var (
	defaultLocale = i18n.LocateZH.String()
	translator, _ = NewTranslator(defaultLocale)
)

func NewTranslator(locale string) (trans ut.Translator, err error) {
	en := en.New()
	zh := zh.New()

	uni := ut.New(en, en, zh)
	translator, ok := uni.GetTranslator(locale)
	if !ok {
		return nil, errorx.New("validator get translator fail")
	}

	switch locale {

	case i18n.LocateEN.String():
		defaultLocale = locale
		err = en_translations.RegisterDefaultTranslations(v, translator)

	case i18n.LocateZH.String():
		defaultLocale = locale
		err = zh_translations.RegisterDefaultTranslations(v, translator)

	default:
		err = zh_translations.RegisterDefaultTranslations(v, translator)
	}

	if err != nil {
		return nil, err
	}

	return translator, nil
}

func GetTranslator() ut.Translator {
	return translator
}

func RegisterTranslation(tag string, msg string) error {
	return v.RegisterTranslation(tag, translator,

		func(trans ut.Translator) error {
			if err := trans.Add(tag, msg, true); err != nil {
				return err
			}
			return nil
		},

		func(ut ut.Translator, fe validator.FieldError) string {
			message, _ := translator.T(tag, fe.Field())
			return message
		},
	)
}
