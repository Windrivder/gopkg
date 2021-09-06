package rest

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/windrivder/gopkg/encoding/jsonx"
	"github.com/windrivder/gopkg/util/valid"
)

type binder struct {
	*echo.DefaultBinder
}

func NewBinder() echo.Binder {
	return &binder{DefaultBinder: &echo.DefaultBinder{}}
}

func (b *binder) Bind(i interface{}, c echo.Context) (err error) {
	if err := b.BindPathParams(c, i); err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet || c.Request().Method == http.MethodDelete {
		if err = b.BindQueryParams(c, i); err != nil {
			return err
		}
	}

	req := c.Request()
	ctype := req.Header.Get(echo.HeaderContentType)
	if strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
		if err = jsonx.DecodeReader(req.Body, i); err != nil {
			return err
		}
	}

	if err = b.BindBody(c, i); err != nil {
		return err
	}

	return valid.ValidateStruct(i)
}
