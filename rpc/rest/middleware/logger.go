package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/windrivder/gopkg/logx"
)

type (
	LoggerConfig struct {
		Skipper SkipperFunc
	}
)

var (
	DefaultLoggerConfig = LoggerConfig{
		Skipper: DefaultSkipper,
	}
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			stop := time.Now()
			fields := logx.Fields{
				"id":        id,
				"path":      path,
				"remote_ip": c.RealIP(),
				"uri":       req.RequestURI,
				"method":    req.Method,
				"status":    res.Status,
				"latency":   stop.Sub(start).String(),
			}
			if err != nil {
				logx.WithFields(fields).Error(err)
			} else {
				logx.WithFields(fields).Info("ok")
			}

			return nil
		}
	}
}
