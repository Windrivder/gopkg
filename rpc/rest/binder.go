package rest

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/windrivder/gopkg/i18n"
	"github.com/windrivder/gopkg/util/valid"
)

type binder struct {
	*echo.DefaultBinder
	locale string
	trans  ut.Translator
}

func NewBinder(locale string) (echo.Binder, error) {
	if locale == "" {
		locale = i18n.LocateZH.String()
	}

	trans, err := valid.NewTranslator(locale)
	if err != nil {
		return nil, err
	}

	return &binder{
		DefaultBinder: &echo.DefaultBinder{},
		locale:        locale,
		trans:         trans,
	}, nil
}

func (b *binder) Bind(i interface{}, c echo.Context) (err error) {
	if err := b.DefaultBinder.Bind(i, c); err != nil {
		return err
	}

	return valid.ValidateStruct(i)
}
