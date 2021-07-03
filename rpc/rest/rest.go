package rest

import "github.com/labstack/echo/v4"

type (
	Route              = echo.Route
	Router             = echo.Router
	Group              = echo.Group
	MiddlewareFunc     = echo.MiddlewareFunc
	HandlerFunc        = echo.HandlerFunc
	HandlerRoutersFunc func(IServer)
)

type IServer interface {
	Router() *Router
	Routers() map[string]*Router
	Pre(middleware ...MiddlewareFunc)
	Use(middleware ...MiddlewareFunc)
	CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route
	Match(methods []string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route
	Static(prefix, root string) *Route
	File(path, file string, m ...MiddlewareFunc) *Route
	Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route
	Host(name string, m ...MiddlewareFunc) (g *Group)
	Group(prefix string, m ...MiddlewareFunc) (g *Group)
	URI(handler HandlerFunc, params ...interface{}) string
	URL(h HandlerFunc, params ...interface{}) string
	Reverse(name string, params ...interface{}) string
	Routes() []*Route
	Start() error
	Stop() error
}
