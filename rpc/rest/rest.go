package rest

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareFunc   = echo.MiddlewareFunc
	HandlerFunc      = echo.HandlerFunc
	HTTPErrorHandler = echo.HTTPErrorHandler

	Router struct {
		Path        string
		Category    []string
		Name        string
		Key         string
		Handler     HandlerFunc
		Middlewares []MiddlewareFunc
		Method      string
	}
	Routers []Router
)

var ProviderSet = wire.NewSet(NewServer, NewClient, NewOptions)
