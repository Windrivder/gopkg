package middleware

import (
	mw "github.com/labstack/echo/v4/middleware"
)

type (
	CORSConfig = mw.CORSConfig
)

var (
	DefaultCORSConfig = mw.DefaultCORSConfig
	CORS              = mw.CORS
	CORSWithConfig    = mw.CORSWithConfig
)
