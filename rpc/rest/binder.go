package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/windrivder/gopkg/util/valid"
)

type binder struct {
	defaultBinder *echo.DefaultBinder
	validate      *valid.Validate
}

func NewBinder(v *valid.Validate) (Binder, error) {
	return &binder{
		defaultBinder: &echo.DefaultBinder{},
		validate:      v,
	}, nil
}

func (b *binder) Bind(i interface{}, c echo.Context) (err error) {
	if err := b.defaultBinder.Bind(i, c); err != nil {
		return err
	}

	return b.validate.ValidateStruct(i)
}
