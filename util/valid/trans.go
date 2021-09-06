package valid

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/windrivder/gopkg/errorx"
)

func NewTranslator(locale string) (trans ut.Translator, err error) {
	en := en.New()
	zh := zh.New()

	uni := ut.New(en, en, zh)
	trans, ok := uni.GetTranslator(locale)
	if !ok {
		return nil, errorx.New("validator get translator fail")
	}

	switch locale {
	case "en":
		err = en_translations.RegisterDefaultTranslations(v, trans)
	case "zh":
		err = zh_translations.RegisterDefaultTranslations(v, trans)
	default:
		err = en_translations.RegisterDefaultTranslations(v, trans)
	}

	if err != nil {
		return nil, err
	}

	return trans, nil
}
