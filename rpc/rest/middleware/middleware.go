package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// The SkipperFunc signature, used to serve the main request without logs.
// See `Configuration` too.
type SkipperFunc = middleware.Skipper

// SkipHandler 统一处理跳过函数
func SkipHandler(ctx echo.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(ctx) {
			return true
		}
	}
	return false
}

// HigherSkipperFunc 进一步包装处理函数
type HigherSkipperFunc func(...string) SkipperFunc

// AllowPathPrefixSkipper 检查请求路径是否包含指定的前缀，如果包含则跳过
func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c echo.Context) bool {
		for _, p := range prefixes {
			if strings.HasPrefix(c.Path(), p) {
				return true
			}
		}
		return false
	}
}

// AllowPathPrefixNoSkipper 检查请求路径是否包含指定的前缀，如果包含则不跳过
func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c echo.Context) bool {
		for _, p := range prefixes {
			if strings.HasPrefix(c.Path(), p) {
				return false
			}
		}
		return true
	}
}
