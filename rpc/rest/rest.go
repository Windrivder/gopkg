package rest

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareFunc   = echo.MiddlewareFunc
	HandlerFunc      = echo.HandlerFunc
	HTTPErrorHandler = echo.HTTPErrorHandler
	Binder           = echo.Binder

	Router struct {
		Path        string
		Category    []string
		Name        string
		Key         string
		Handler     HandlerFunc
		Middlewares []MiddlewareFunc
		Method      string
		Request     interface{}
		Response    interface{}
	}
	Routers []Router
)

var ProviderSet = wire.NewSet(NewServer, NewClient, NewOptions, NewBinder)
