package middleware

import (
	mw "github.com/labstack/echo/v4/middleware"
)

type (
	RequestIDConfig = mw.RequestIDConfig
)

var (
	DefaultRequestIDConfig = mw.DefaultRequestIDConfig
	RequestID              = mw.RequestID
	RequestIDWithConfig    = mw.RequestIDWithConfig
)
