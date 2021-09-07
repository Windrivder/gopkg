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

func NewTranslator(o i18n.Options) (trans ut.Translator, err error) {
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

func RegisterTranslation(tag string, msg string, trans ut.Translator) error {
	return v.RegisterTranslation(tag, trans,

		func(translator ut.Translator) error {
			if err := translator.Add(tag, msg, true); err != nil {
				return err
			}
			return nil
		},

		func(translator ut.Translator, fe validator.FieldError) string {
			message, _ := trans.T(tag, fe.Field())
			return message
		},
	)
}
