package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/windrivder/gopkg/util/valid"
)

type binder struct {
	*echo.DefaultBinder
}

func NewBinder() (echo.Binder, error) {
	return &binder{DefaultBinder: &echo.DefaultBinder{}}, nil
}

func (b *binder) Bind(i interface{}, c echo.Context) (err error) {
	if err := b.DefaultBinder.Bind(i, c); err != nil {
		return err
	}

	return valid.ValidateStruct(i)
}
