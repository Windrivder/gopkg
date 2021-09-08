package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/logx"
)

type Options struct {
	Name            string        `json:"Name"`
	Mode            string        `json:"Mode"`
	Host            string        `json:"Host"`
	Port            int           `json:"Port"`
	CertFile        string        `json:"CertFile"`
	KeyFile         string        `json:"KeyFile"`
	ShutdownTimeout time.Duration `json:"ShutdownTimeout"`
	ClientTimeout   time.Duration `json:"ClientTimeout"`
	Secret          string        `json:"Secret"`
	Expired         time.Duration `json:"Expired"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("Rest", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal rest option error")
	}

	return o, err
}

type Server struct {
	e *echo.Echo
	o Options
}

func NewServer(o Options, rs Routers) (*Server, func(), error) {
	e := echo.New()

	e.Use(middleware.Recover())

	for _, r := range rs {
		e.Add(r.Method, r.Path, r.Handler, r.Middlewares...)
	}

	s := &Server{o: o, e: e}

	return s, func() { s.Stop() }, nil
}

func (s *Server) Engine() *echo.Echo {
	return s.e
}

func (s *Server) Start() (err error) {
	addr := fmt.Sprintf("%s:%d", s.o.Host, s.o.Port)
	logx.Info().Str("addr", addr).Msg("http server starting...")

	go func() {
		s.e.Server.Addr = addr

		if s.o.CertFile == "" && s.o.KeyFile == "" {
			err = s.e.Server.ListenAndServe()
		} else {
			err = s.e.Server.ListenAndServeTLS(s.o.CertFile, s.o.KeyFile)
		}

		if err != nil && err != http.ErrServerClosed {
			logx.Fatal().Msgf("start http server err: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	logx.Info().Msg("http server stopping...")

	timeout := time.Second * s.o.ShutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.e.Shutdown(ctx)
}
